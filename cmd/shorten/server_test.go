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
	const url = "https://blog.golang.org/vscode-go"

	tests := []struct {
		name string
		body string
		want int
	}{
		{
			name: "shorten+redirect",
			body: `{"long_url": "` + url + `"}`,
			want: http.StatusOK,
		},

		// error cases
		{
			name: "empty url",
			body: `{"long_url": ""}`,
			want: http.StatusBadRequest,
		},
		{
			name: "malformed url",
			body: `{"long_url": ":/::"}`,
			want: http.StatusBadRequest,
		},
		{
			name: "invalid body",
			body: `{"long_url": "`,
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := newServer(shorten.NewMemoryStore())

			// Call /v1/shorten with an URL.
			req := httptest.NewRequest("POST", "http://myserver/v1/shorten", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			s.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Fatalf("/v1/shorten want HTTP %d, got %d", tt.want, w.Code)
			}

			if w.Code != http.StatusOK {
				return
			}

			resp := w.Result()
			respBody, _ := io.ReadAll(resp.Body)
			shortURL := strings.TrimSpace(string(respBody))

			// Call the shortened URL.
			req = httptest.NewRequest("GET", "http://myserver.com/"+shortURL, nil)
			w = httptest.NewRecorder()

			s.ServeHTTP(w, req)
			if w.Code != http.StatusMovedPermanently {
				t.Errorf("/v1/shorten want HTTP 301, got %d", w.Code)
			}
		})
	}
}
