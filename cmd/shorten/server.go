package main

import (
	"encoding/json"
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

// failMsg writes sets the HTTP header, writes an error message into w and log it.
func failMsg(w http.ResponseWriter, r *http.Request, code int, format string, args ...interface{}) {
	var sb strings.Builder

	fmt.Fprintf(&sb, r.URL.Path)
	sb.WriteByte(' ')
	fmt.Fprintf(&sb, format, args...)

	w.WriteHeader(code)
	fmt.Fprintln(w, sb.String())
	log.Println(sb.String())
}

func (s *server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		failMsg(w, r, http.StatusMethodNotAllowed, "wrong method %q", r.Method)
		return
	}

	var args ShortenArgs

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&args); err != nil {
		failMsg(w, r, http.StatusBadRequest, "invalid json %v", err)
		return
	}

	if err := args.validate(); err != nil {
		failMsg(w, r, http.StatusBadRequest, "failed arg validation %v", err)
		return
	}

	shortURL := s.shortener.Shorten(args.LongURL)
	fmt.Fprintln(w, shortURL)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	longURL, ok := s.shortener.Long(shortURL)
	if !ok {
		failMsg(w, r, http.StatusNotFound, "not found")
		return
	}

	log.Println(r.URL.Path, "redirecting to", longURL)
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
