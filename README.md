go-search
=========

- [Interview Excercise](#introduction)
- [Solution](#solution)
- [Usage](#usage)
- [Building](#building)

## Introduction

Google for a Day

It’s the year 1992, and you are the developer in charge of implementing a simple search engine to rank web pages based on the number of occurrences of certain words. To do this, you’ll create an application featuring two primary functions: constructing and resetting a search index, and a simple interface to search and rank pages based on the supplied criteria.

To create the search index the application will take a URL as an argument and perform the following actions:

● It will open and parse the content of such URL.

● If it finds a link going out to another page, it will follow that link and parse that URL as well.

● It will also index every word (non-HTML tag) that finds on the page.

Since the web could get quite complicated to crawl, you should limit the depth of the crawling functionality to no more than three levels deep. For example, let's assume the following list of pages:

● Page A links to Page B

● Page B links to Page C

● Page C links to Page D

Assuming Page A is the one provided to your application, you should index pages A, B, and C (being Page A level 1, Page B level 2, and Page C level 3).

By repeatedly indexing multiple URL’s we should be able to increase the index size of the application; every time the user provides a URL, your application should add to the existing index. There will also be a button to clear the index and start fresh.

Any time the application processes a new URL, the page should display the number of indexed pages and the number of indexed words.

Keep in mind the following:

● You shouldn’t index HTML tags on the page, only the real content.

● Keep and store the title of the page (the one found at `<head><title>`). You’ll need this to display the results.

● Links going out to other pages are those specified as part of an HTML anchor tag (`<a>`). You’ll always find the link itself in the href attribute of the anchor tag.

● Links could go to a different website or be relative links to a different page on the same current domain.

● Make sure the crawler doesn’t get into infinite loops (Page A links to Page B that links back to Page A).

The search page will also be very simple and will let users search for a particular word. For the purpose of this assignment, you shouldn’t worry about phrases at all; users will only use single words to search.

After searching, the results should display the title of the page and the number of occurrences of the word on that particular page. Results should be ordered by the number of occurrences (the more, the higher the rank) and you should not display pages with zero occurrences of the word. Clicking on the title of the page should open the URL in a different tab of the browser.

Here are some additional notes to help you work through this problem:

● Your goal is to create a very simple application, not reproduce Google.com. Searching is a complex problem to solve, but we just need a very simple prototype. Feel free to make assumptions that don’t subtract from the usability of your prototype but help you deliver on time.

● There are great products for indexing information (for example, Apache Lucene). For this exercise, we’d like you not to use any of them. Also, don’t over think the indexing too much. Stringing together a bunch of services and infrastructure isn’t required.

● You’ll have exactly seven days from the time we give you this problem. We believe it should take you less time than that, but we want to give you as much room to succeed as possible.

● Your exercise should be pushed to a publically accessible place such as Github, Gitlab, Bitbucket.

● Feel free to ask as many questions for clarifications, by email, as you need.

## Solution
The solution is a golang command line _interactive_ application.

## Usage

```
> ./go-search
Commands: add, search and reset
> add https://yahoo.com
.............................................................................................................................................................................................................................DONE!
pages indexed: 108
 words indexed: 31525
Commands: add, search and reset
> search free
Title: Yahoo Finance Premium (https://finance.yahoo.com/premium-marketing) (hits 11)
Title: Pets Hub | Yahoo Lifestyle (https://www.yahoo.com/lifestyle/tagged/pets-hub) (hits 5)
Title: How Instagram Influencers Make Life Miserable For Small Tourist Towns (https://www.yahoo.com/huffpost/instagram-influencers-overtourism.html) (hits 2)
Title: Style | Yahoo Lifestyle (https://www.yahoo.com/lifestyle/style/) (hits 2)
Title: Shopping | Yahoo Lifestyle (https://www.yahoo.com/lifestyle/tagged/shopping/) (hits 1)
Commands: add, search and reset
>
```
## Building
```
go build -mod=vendor
```
the binary `go-search` (or `go-search.exe` on Windows) is generated in current directory
