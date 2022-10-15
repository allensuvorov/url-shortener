package compress

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/config"
	handlers "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/stretchr/testify/assert"
)

func TestGzipHandler_GzipMiddleware(t *testing.T) {
	/* Basic Logic
	Make request with gzip-ed body
	want - decoded URL passed to
	need: handler, request, test data
	*/

	// Create Handler
	config.BuildConfig()
	usm := storage.NewURLStorage()
	us := service.NewURLService(usm)
	uh := handlers.NewURLHandler(us)

	tests := []struct {
		name                     string
		url                      string
		headerAcceptEncoding     string // Accept-Encoding
		headerContentEncoding    string // Content-Encoding
		expectedResponseBody     []byte
		expectedResponseCEHeader string
	}{
		// TODO: Add test cases.
		{
			name:                     "decoded",
			url:                      "http://www.booking.com/",
			headerAcceptEncoding:     "gzip",
			headerContentEncoding:    "gzip",
			expectedResponseBody:     []byte(`http://localhost:8080/4cd89a20`),
			expectedResponseCEHeader: "gzip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Need to gzip-encode the URL before passing to request
			var buf bytes.Buffer       // Will write to buf
			gz := gzip.NewWriter(&buf) // создаём gzip.Writer

			_, err := gz.Write([]byte(tt.url))
			if err != nil {
				log.Fatalf("failed write data to compress temporary buffer: %v", err)
			}

			err = gz.Close()
			if err != nil {
				log.Fatalf("failed compress data: %v", err)
			}

			// Create POST request
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", &buf)
			w := httptest.NewRecorder()
			g := GzipHandler{}

			// Add headers
			r.Header.Set("Content-Encoding", tt.headerAcceptEncoding)
			r.Header.Set("Accept-Encoding", tt.headerAcceptEncoding)

			// uh.Create(w, r) // Do we need to call the original handler?

			// Handler wrapped in middleware
			h := g.GzipMiddleware(http.HandlerFunc(uh.Create)) // h - is a struct

			// Call Middleware
			h.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedResponseBody, w.Body.Bytes())

		})
	}
}
