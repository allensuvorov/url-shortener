package url

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
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
		},
		{
			name: "2st Test Case: invalide long URL",
			fields: fields{
				urlService: service.NewURLService(storage.NewURLStorage()),
			},
			args: args{
				longURL:    "123http://www.apple.com/store",
				requestURL: "localhost:8080/",
				StatusCode: http.StatusInternalServerError,
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

			if res.StatusCode == http.StatusCreated {
				// Check response body
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("coult not read respons: %v", err)
				}
				if string(b) != tt.args.shortURL {
					t.Fatalf("expected %s, got %s", tt.args.shortURL, string(b))
				}
			}
		})
	}
}

func TestURLHandler_Get(t *testing.T) {
	// New url entity
	// ue := &entity.URLEntity{
	// 	URL:  "123http://www.apple.com/store",
	// 	Hash: "a7d59904",
	// }

	sm := storage.NewURLStorage() // storage mock
	// log.Println(sm, ue)
	// sm.inMemory = append(sm.inMemory, ue)

	type URLStorageMock struct {
		inMemory []*entity.URLEntity
	}

	// Create new entity (pair shortURL: longURL).

	// func (usm *URLStorageMock) Create(ue *entity.URLEntity) error {
	// 	usm.inMemory = append(usm.inMemory, ue)
	// 	return nil
	// }

	// GetByHash returns entity for the matching hash, checks if hash exists.
	// func abc(b) {
	// 	return b
	// }

	// func (usm *URLStorageMock) GetURLByHash(u string) (string, error) {
	// 	return "", nil
	// }

	// GetByURL returns hash for the matching URL, checks if URL exists.
	// func (usm *URLStorageMock) GetHashByURL(u string) (string, error) {
	// 	return "", nil
	// }

	type fields struct {
		urlService URLService
	}
	type args struct {
		requestURL string
		longURL    string
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1st Test Case: positive - apple/store",
			fields: fields{
				urlService: service.NewURLService(sm),
			},
			args: args{
				requestURL: "http://localhost:8080/a7d59904",
				longURL:    "http://www.apple.com/store",
				StatusCode: http.StatusTemporaryRedirect,
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// request
			requestURL := tt.args.requestURL
			req, err := http.NewRequest("GET", requestURL, nil)

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
			if res.StatusCode != http.StatusTemporaryRedirect {
				t.Errorf("expected status TemporaryRedirect; got %v", res.Status)
			}
			if res.StatusCode == http.StatusTemporaryRedirect {
				if res.Header.Get("Location") != "http://www.apple.com/store" {
					t.Errorf("Expected Header http://www.apple.com/store, got %s", res.Header.Get("Location"))
				}
			}

		})
	}
}
