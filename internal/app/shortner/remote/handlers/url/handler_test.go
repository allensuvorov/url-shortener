package url

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func TestURLHandler_Create(t *testing.T) {
	type fields struct {
		urlService URLService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	// Mock data
	// Object
	URLService := service.NewURLService(storage.NewURLStorage())

	// Mock Request
	jsonBody := []byte("http://www.apple.com/store")
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := "localhost:8080/"

	req, err := http.NewRequest("POST", requestURL, bodyReader)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	// test data
	tests := []struct {
		name string
		fields
		args args
	}{
		{
			name:   "apple store",
			fields: fields{urlService: URLService},
			args:   args{rec, req},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := URLHandler{
				urlService: tt.fields.urlService,
			}
			uh.Create(tt.args.w, tt.args.r)
		})
	}
}

func TestURLHandler_Get(t *testing.T) {
	type fields struct {
		urlService URLService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := URLHandler{
				urlService: tt.fields.urlService,
			}
			uh.Get(tt.args.w, tt.args.r)
		})
	}
}
