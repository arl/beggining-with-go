package shorten

import (
	"encoding/json"
	"io"
	"log"

	"github.com/jxskiss/base62"
)

type URLShortener struct {
	s Store
}

// NewURLShortener creates an URL shortener.
func NewURLShortener(s Store) *URLShortener {
	return &URLShortener{s: s}
}

// Shorten returns the short URL key for the given long URL.
func (s *URLShortener) Shorten(long string) string {
	id := s.s.Add(long)
	return string(base62.FormatUint(id))
}

// Long returns the previously shortened long URL associated
// with the given short URL key, or false if it doesn't exist.
func (s *URLShortener) Long(short string) (string, bool) {
	key, err := base62.ParseUint([]byte(short))
	if err != nil {
		log.Printf("base62 decoding failed: %v", err)
		return "", false
	}

	return s.s.Load(uint64(key))
}

type urlPair struct {
	Short string `json:"short_url"`
	Long  string `json:"long_url"`
}

type urlList struct {
	URLS []urlPair `json:"urls"`
}

func (s *URLShortener) WriteJSON(w io.Writer) error {
	urls := make([]urlPair, 0)

	s.s.ForEach(func(key uint64, value string) bool {
		urls = append(urls, urlPair{Short: string(base62.FormatUint(key)), Long: value})
		return true
	})

	enc := json.NewEncoder(w)
	return enc.Encode(urlList{URLS: urls})
}
