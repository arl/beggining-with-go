package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arl/shorten"
)

func TestShortenAndRedirect(t *testing.T) {
	s := newServer(shorten.NewMemoryStore())

	// Call /v1/shorten with an URL.
	const url = "https://blog.golang.org/vscode-go"
	reqBody := `{"long_url": "` + url + `"}`
	req := httptest.NewRequest("POST", "http://myserver/v1/shorten", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("/v1/shorten %v want HTTP 200, got %d", reqBody, w.Code)
	}

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	shortURL := strings.TrimSpace(string(respBody))

	// Call the shortened URL.
	req = httptest.NewRequest("GET", "http://myserver.com/"+shortURL, nil)
	w = httptest.NewRecorder()

	s.ServeHTTP(w, req)
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("/v1/shorten %v want HTTP 301, got %d", reqBody, w.Code)
	}
}
