package shorten

import (
	"testing"
)

func TestURLShortener(t *testing.T) {
	const longURL = "llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.co.uk"

	s := NewURLShortener(NewMemoryStore())
	short := s.Shorten(longURL)

	if got, ok := s.Long(short); got != longURL || !ok {
		t.Errorf("Long(%s) = (%v, %v) want (longURL, true)", short, got, ok)
	}

	if got, ok := s.Long("foobar"); ok {
		t.Errorf("Long(foobar) = (%v, %v) want (_, false)", got, ok)
	}
}
