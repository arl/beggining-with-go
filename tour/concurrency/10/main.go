package main

import (
	"fmt"
	"sync"
)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	c := crawler{
		fetched: make(map[string]struct{}),
	}

	c.crawl(url, depth, fetcher)
	c.wg.Wait()
}

type crawler struct {
	wg sync.WaitGroup

	// Protects the URL map from concurrent accesses.
	mu      sync.Mutex
	fetched map[string]struct{}
}

func (c *crawler) crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	// Dont fetch the same URL twice.
	c.mu.Lock()
	_, ok := c.fetched[url]
	if ok {
		// We've already fetched that URL.
		c.mu.Unlock()
		return
	}

	c.fetched[url] = struct{}{}
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			c.crawl(u, depth-1, fetcher)
		}
	}()
}

// --

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
