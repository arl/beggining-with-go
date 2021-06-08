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

type server struct {
	shortener *shorten.URLShortener
	router    *http.ServeMux
}

func newServer(store shorten.Store) *server {
	s := &server{
		shortener: shorten.NewURLShortener(store),
		router:    http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *server) routes() {
	s.router.HandleFunc("/v1/shorten", s.handleShorten)
	s.router.HandleFunc("/", s.handleIndex)
}

func (s *server) handleShorten(w http.ResponseWriter, r *http.Request) {
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

	shortURL := s.shortener.Shorten(args.LongURL)
	fmt.Fprintln(w, shortURL)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	longURL, ok := s.shortener.Long(shortURL)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, r.URL.Path, "not found")
		log.Println(r.URL.Path, "not found")
		return
	}

	log.Println(r.URL.Path, "redirecting to", longURL)
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func main() {
	addrFlag := ""
	flag.StringVar(&addrFlag, "addr", ":8080", "server listen address")
	flag.Parse()

	s := newServer(shorten.NewMemoryStore())

	if err := http.ListenAndServe(addrFlag, s); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
