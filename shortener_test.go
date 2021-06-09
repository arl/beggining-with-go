package shorten

import (
	"bytes"
	"strings"
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

	bb := bytes.Buffer{}
	if err := s.WriteJSON(&bb); err != nil {
		t.Errorf("WriteJSON error: %v", err)
	}

	got := strings.TrimSpace(bb.String())
	want := `{"urls":[{"short_url":"` + short + `","long_url":"llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.co.uk"}]}`
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%v", got, want)
	}
}
