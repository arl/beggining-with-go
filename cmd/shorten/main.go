package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// ShortenArgs holds the arguments to the /shorten endpoint.
type ShortenArgs struct {
	LongURL string `json:"long_url"`
}

// validate validates the content of the ShortenArgs struct.
func (s *ShortenArgs) validate() error {
	if s.LongURL == "" {
		return fmt.Errorf("long_url is empty")
	}

	_, err := url.Parse(s.LongURL)
	if err != nil || s.LongURL == "" {
		return fmt.Errorf("invalid long_url %v: %v", s.LongURL, err)
	}
	return nil
}

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	if err := http.ListenAndServe(addrFlag, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
