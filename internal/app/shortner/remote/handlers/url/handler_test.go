package url

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
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

type URLStorageMock struct {
	inMemory []*entity.URLEntity
}

// Create new entity (pair shortURL: longURL).
func (usm *URLStorageMock) Create(ue *entity.URLEntity) error {
	// usm.inMemory = append(usm.inMemory, ue)
	return nil
}

// GetByHash returns entity for the matching hash, checks if hash exists.
func (usm *URLStorageMock) GetURLByHash(u string) (string, error) {
	log.Println("Testing GetURLByHash, looking in slice len", len(usm.inMemory))
	log.Println("Testing GetURLByHash, looking for matching Hash =", u)
	log.Println("URL Storage Mock inMemory", *usm.inMemory[0])
	for _, v := range usm.inMemory {
		log.Printf("Testing GetURLByHash, comparing local %s and received %s", v.Hash, u)
		if v.Hash == u {
			log.Println("Testing GetURLByHash, found ue", v)
			return v.URL, nil
		}
	}
	return "", errors.NotFound
}

// GetByURL returns hash for the matching URL, checks if URL exists.
func (usm *URLStorageMock) GetHashByURL(u string) (string, error) {
	return "", nil
}

func TestURLHandler_Get(t *testing.T) {
	log.Println("Starting TestURLHandler_Get")

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
				urlService: service.NewURLService(&URLStorageMock{
					inMemory: []*entity.URLEntity{
						{
							URL:  "http://www.apple.com/store",
							Hash: "a7d59904",
						},
					}}),
			},
			args: args{
				requestURL: "http://localhost:8080/a7d59904",
				longURL:    "http://www.apple.com/store",
				StatusCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "2st Test Case: negative - not found",
			fields: fields{
				urlService: service.NewURLService(&URLStorageMock{
					inMemory: []*entity.URLEntity{
						{
							URL:  "http://www.apple.com/store",
							Hash: "a7d59904",
						},
					}}),
			},
			args: args{
				requestURL: "http://localhost:8080/111111",
				longURL:    "http://www.apple.com/store",
				StatusCode: http.StatusBadRequest,
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// request
			requestURL := tt.args.requestURL
			log.Println("Test Get, requestURL is", requestURL)
			req, err := http.NewRequest("GET", requestURL, nil)

			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// response
			rec := httptest.NewRecorder()

			// handler object
			// TODO: decide - should we build the object in this func?
			// New url entity
			// ue := &entity.URLEntity{
			// 	URL:  "http://www.apple.com/store",
			// 	Hash: "a7d59904",
			// }

			// usm := new(URLStorageMock)
			// usm.inMemory = append(usm.inMemory, ue)

			uh := URLHandler{
				urlService: tt.fields.urlService,
			}

			// Run the handler
			uh.Get(rec, req)

			// Check response
			res := rec.Result()
			defer res.Body.Close()

			// Check response status code
			if res.StatusCode != tt.args.StatusCode {
				t.Errorf("expected status %v; got %v", tt.args.StatusCode, res.Status)
			}
			if res.StatusCode == http.StatusTemporaryRedirect {
				if res.Header.Get("Location") != tt.args.longURL {
					t.Errorf("Expected Header %s, got %s", tt.args.longURL, res.Header.Get("Location"))
				}
			}

		})
	}
}
