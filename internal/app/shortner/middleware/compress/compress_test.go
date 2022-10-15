package compress

import (
	"bytes"
	"compress/gzip"
	"io"
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

type gzipWriterForTest struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriterForTest) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

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
		expectedRequestBody      string
		expectedResponseCEHeader string
	}{
		// TODO: Add test cases.
		{
			name:                     "decoded",
			url:                      "http://www.booking.com/",
			headerAcceptEncoding:     "gzip",
			headerContentEncoding:    "gzip",
			expectedRequestBody:      "http://www.booking.com/",
			expectedResponseCEHeader: "gzip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// TODO: gzip the URL before passing to request
			var buf bytes.Buffer

			// создаём gzip.Writer поверх текущего w
			gz := gzip.NewWriter(&buf)
			defer gz.Close()
			_, err := gz.Write([]byte(tt.url))
			if err != nil {
				log.Fatalf("failed write data to compress temporary buffer: %v", err)
			}
			// Create POST request
			// b := bytes.NewBufferString(tt.url)
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", &buf)
			w := httptest.NewRecorder()
			g := GzipHandler{}

			// Add headers
			r.Header.Set("Content-Encoding", tt.headerAcceptEncoding)
			r.Header.Set("Accept-Encoding", tt.headerAcceptEncoding)

			// Do we need to call the original handler?
			// uh.Create(w, r)

			// Handler wrapped in middleware
			h := g.GzipMiddleware(http.HandlerFunc(uh.Create)) // h - is a struct

			// Call Middleware
			// Do we need to call the updated handler?
			h.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedRequestBody, r.Body)
			log.Println(r.Body)
			assert.Equal(t, tt.expectedResponseCEHeader, w.Header())

		})
	}
}
