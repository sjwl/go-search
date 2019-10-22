package main

import (
	"fmt"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/sjwl/go-search/sengine"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "add", Description: "add url to index"},
		{Text: "search", Description: "search for word"},
		{Text: "reset", Description: "reset search index"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

//
// An interactive CLI that does basic web crawling and build an index of searched words of pages
//
func main() {

	index := sengine.NewIndex()

	for {
		fmt.Println("Commands: add, search and reset")
		t := strings.Split(prompt.Input("> ", completer), " ")
		switch t[0] {
		case "add":
			if len(t) < 2 {
				fmt.Println("add needs url")
			} else {
				url := t[1]
				pagesIndexed, wordsIndexed, err := index.Add(url)
				if err != nil {
					fmt.Printf("error: %s\n", err.Error())
				} else {
					fmt.Printf("pages indexed: %d\n words indexed: %d\n", pagesIndexed, wordsIndexed)
				}
			}
		case "search":
			if len(t) < 2 {
				fmt.Println("search must have a word")
			}
			word := t[1]
			if entries, ok := index.GetWord(word); !ok {
				fmt.Println("No match")
			} else {
				for _, entry := range entries {
					hits, _ := entry.GetHits(word)
					fmt.Printf("Title: %s (%s) (hits %d)\n", entry.Title, entry.Url, hits)
				}
			}
		case "reset":
			index.Reset()
			fmt.Println("Index cleared")
		case "exit":
			goto exit
		default:
			fmt.Println("Invalid command: " + t[0])
		}
	}
exit:
}
