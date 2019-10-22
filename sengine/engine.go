package sengine

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

//
// Entry is the result of a single crawl of one page
//
type Entry struct {
	wordHits map[string]uint64 // how many times a word appears
	Title    string            // the title of the page
	Url      string            // page url
	lock     *sync.RWMutex     // makes accessing map conurrently safe
}

//
// GetHits
//  - given a word return number of search hits
func (e Entry) GetHits(word string) (uint64, bool) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	hits, found := e.wordHits[word]
	return hits, found
}

type WordToEntries map[string][]*Entry

type Index struct {
	Client      *http.Client
	UrlToEntry  map[string]*Entry // the string will be a url
	WordResults WordToEntries     // word to a list of Entries
	lock        *sync.RWMutex     // makes accessing map concurrently safe
	lock2       *sync.RWMutex
}

//
// NewIndex
//   - create new instance of Index
func NewIndex() *Index {
	wu := make(map[string][]*Entry)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &Index{&client, map[string]*Entry{}, wu, &sync.RWMutex{}, &sync.RWMutex{}}
}

//
// GetEntry
//   - Given url return entry
//   - concurrency safe
func (i Index) GetEntry(Url string) (*Entry, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	val, found := i.UrlToEntry[Url]
	return val, found
}

//
// PutEntry
//   - Given url and an entry, store it
//   - concurrency safe
func (i Index) PutEntry(Url string, entry *Entry) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.UrlToEntry[Url] = entry
}

//
// GetWord
//   - Given a word, return a list of entries
//   - concurrency safe
func (i Index) GetWord(word string) ([]*Entry, bool) {
	i.lock2.RLock()
	defer i.lock2.RUnlock()
	val, found := i.WordResults[word]
	sort.Slice(val, func(i, j int) bool {
		h1, _ := val[i].GetHits(word)
		h2, _ := val[j].GetHits(word)
		return h1 > h2
	})
	return val, found
}

//
// PutWord
//   - Given word and entry, add to the list
//   - concurrency safe
func (i Index) PutWord(word string, entry *Entry) {
	i.lock2.Lock()
	defer i.lock2.Unlock()
	i.WordResults[word] = append(i.WordResults[word], entry)
}

//
// GetHits
//   - Given a url and word, return hits
func (i Index) GetHits(Url string, word string) uint64 {
	if e, ok := i.GetEntry(Url); ok {
		if h, ok := e.GetHits(word); ok {
			return h
		}
	}

	return uint64(0)
}

//
// Reset
//   - clear index and start empty
func (i Index) Reset() {

	for k := range i.UrlToEntry {
		delete(i.UrlToEntry, k)
	}

	for k := range i.WordResults {
		delete(i.WordResults, k)
	}
}

type counter struct {
	pages, words int
}

//
// Add -
//  Add a url to index and return the total pages indexed and words found
//
func (i Index) Add(Url string) (int, int, error) {
	var wg sync.WaitGroup
	var totalPages, totalWords int
	counts := make(chan counter)

	wg.Add(1)
	go i.Crawl(Url, 3, counts, &wg)

	go func() {
		for work := range counts {
			totalPages += work.pages
			totalWords += work.words
		}
	}()
	wg.Wait()
	fmt.Println("DONE!")
	close(counts)

	return totalPages, totalWords, nil
}

//
// Crawl -
//  take a url find all the words on the page
//  will also recursively crawl any links found on the page
//  depth is how many levels deep will look for embedded links
//  returns total pages indexed and words found
func (i Index) Crawl(Url string, depth int, counts chan counter, wg *sync.WaitGroup) {

	defer wg.Done()

	wh := make(map[string]uint64)

	entry := Entry{wh, "", Url, &sync.RWMutex{}}

	resp, err := i.Client.Get(Url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	previousStartToken := z.Token()
	for {
		switch z.Next() {
		case html.ErrorToken:
			// End of document
			goto forEnd
		case html.StartTagToken:
			previousStartToken = z.Token()
			if depth > 0 {
				// look for embedded links to crawl further
				if previousStartToken.Data == "a" { //anchor
					for _, a := range previousStartToken.Attr {
						if a.Key == "href" {
							// skip anything that is not a url
							if _, err := url.Parse(a.Val); err == nil {
								// don't include in page anchors (#)
								if strings.HasPrefix(a.Val, "#") {
									continue
								}
								resolvedUrl := fixUrl(a.Val, Url)
								if strings.HasPrefix(resolvedUrl, "http") {

									wg.Add(1)
									go i.Crawl(resolvedUrl, depth-1, counts, wg)
									fmt.Printf(".")
								}
							}
						}
					}
				}
			}

		case html.TextToken:

			// skip if url already indexed
			if _, ok := i.GetEntry(Url); ok {
				continue
			}

			if previousStartToken.Data == "script" ||
				previousStartToken.Data == "style" {
				continue
			}
			// the first title we find must be <head><title>
			if previousStartToken.Data == "title" && entry.Title == "" {
				entry.Title = string(z.Text())
			}
			f := func(c rune) bool {
				return !unicode.IsLetter(c)
			}
			for _, w := range strings.FieldsFunc(string(z.Text()), f) {
				if _, ok := wh[w]; !ok {
					wh[w] = uint64(1)
					i.PutWord(w, &entry)
				} else {
					wh[w]++
				}
			}
		}

	}
forEnd:
	// only increment page and word counts if newly indexed url
	if _, ok := i.GetEntry(Url); !ok {
		i.PutEntry(Url, &entry)
		c := counter{1, len(wh)}
		counts <- c
	}
}

// convert relative paths into full paths
func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
