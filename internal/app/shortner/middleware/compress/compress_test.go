package compress

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/config"
	handlers "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
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
		name                   string
		url                    string
		headerAE               string // Accept-Encoding
		headerCE               string // Content-Encoding
		expectedRequestBody    string
		expectedResponseHeader string
	}{
		// TODO: Add test cases.
		{
			// name: "decoded",
			// arg:

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create POST request
			b := bytes.NewBufferString(tt.url)
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", b)
			w := httptest.NewRecorder()

			// Creat Handler
			uh.Create(w, r)

			g := GzipHandler{}
			if got := g.GzipMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipHandler.GzipMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
