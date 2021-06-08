package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/arl/shorten"
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

	store := shorten.NewMemoryStore()
	shortener := shorten.NewURLShortener(store)

	http.HandleFunc("/v1/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Println(r.URL.Path)
			fmt.Fprintln(w, "/v1/shorten: wrong method", r.Method)
			log.Println("/v1/shorten: wrong method", r.Method)
			return
		}

		var args ShortenArgs

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&args); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "/v1/shorten: invalid json:", err)
			log.Println("/v1/shorten: invalid json:", err)
			return
		}

		if err := args.validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "/v1/shorten: failed arg validation:", err)
			log.Println("/v1/shorten: failed arg validation:", err)
			return
		}

		shortURL := shortener.Shorten(args.LongURL)
		fmt.Fprintln(w, shortURL)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimPrefix(r.URL.Path, "/")
		longURL, ok := shortener.Long(shortURL)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, r.URL.Path, "not found")
			log.Println(r.URL.Path, "not found")
			return
		}

		fmt.Println("redirecting to", longURL)
	})

	if err := http.ListenAndServe(addrFlag, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
