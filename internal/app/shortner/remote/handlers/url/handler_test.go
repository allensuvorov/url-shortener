package url

import (
	"bytes"
	"io/ioutil"
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
		longURL    string
		requestURL string
		shortURL   string
		StatusCode int
	}

	// test data
	tests := []struct {
		name string
		fields
		args args
	}{
		{
			name: "1st Test Case: apple/store",
			fields: fields{
				urlService: service.NewURLService(storage.NewURLStorage()),
			},
			args: args{
				longURL:    "http://www.apple.com/store",
				requestURL: "localhost:8080/",
				shortURL:   "http://localhost:8080/a7d59904",
				StatusCode: http.StatusCreated,
			},
		}, // TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// request
			jsonBody := []byte(tt.args.longURL)
			bodyReader := bytes.NewReader(jsonBody)
			requestURL := tt.args.requestURL
			req, err := http.NewRequest("POST", requestURL, bodyReader)

			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			// response
			rec := httptest.NewRecorder()

			// handler object
			uh := URLHandler{
				urlService: tt.fields.urlService,
			}
			// Run the handler
			uh.Create(rec, req)

			// Check response
			res := rec.Result()
			defer res.Body.Close()
			// Check response status code
			if res.StatusCode != tt.args.StatusCode {
				t.Errorf("expected status Created; got %v", res.Status)
			}
			// Check response body
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("coult not read respons: %v", err)
			}
			if string(b) != tt.args.shortURL {
				t.Fatalf("expected %s, got %s", tt.args.shortURL, string(b))
			}
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
