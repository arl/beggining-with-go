package shorten

import (
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
