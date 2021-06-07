package shorten

type URLShortener struct {
	s Store
}

// NewURLShortener creates an URL shortener.
func NewURLShortener(s Store) *URLShortener {
	return &URLShortener{s: s}
}

// Shorten returns the short URL key for the given long URL.
func (s *URLShortener) Shorten(long string) string {
	return ""
}

// Long returns the previously shortened long URL associated
// with the given short URL key, or false if it doesn't exist.
func (s *URLShortener) Long(short string) (string, bool) {
	return "", false
}
