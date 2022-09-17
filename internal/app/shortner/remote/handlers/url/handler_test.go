package url

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			assert.Equal(t, res.StatusCode, tt.args.StatusCode, "expected status Created, got other")

			if res.StatusCode == http.StatusCreated {
				// Check response body
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("coult not read respons: %v", err)
				}
				require.Equal(t, tt.args.shortURL, string(b), "short URL is not matching")
			}
		})
	}
}

type handlerGetTest struct {
	name  string
	input input
	want  want
}
type md struct { // mock data
	URL  string
	Hash string
}
type input struct {
	md         md
	requestURL string
}
type want struct {
	requestURL string
	longURL    string
	StatusCode int
}

func (st handlerGetTest) run(t *testing.T) {
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
	// New url entity
	ue := &entity.URLEntity{
		URL:  tt.fields.URL,
		Hash: tt.fields.Hash,
	}

	usm := storage.NewURLStorage()
	usm.InMemory = append(usm.InMemory, ue)

	uh := URLHandler{
		urlService: service.NewURLService(usm),
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
}

var tests = []handlerGetTest{
	{
		name: "1st Test Case: positive - apple/store",
		fields: md{
			URL:  "http://www.apple.com/store",
			Hash: "a7d59904",
		},
		args: want{
			requestURL: "http://localhost:8080/a7d59904",
			longURL:    "http://www.apple.com/store",
			StatusCode: http.StatusTemporaryRedirect,
		},
	},
	{
		name: "2st Test Case: negative - not found",
		fields: md{
			URL:  "http://www.apple.com/store",
			Hash: "a7d59904",
		},
		args: want{
			requestURL: "http://localhost:8080/111111",
			longURL:    "http://www.apple.com/store",
			StatusCode: http.StatusBadRequest,
		},
	}, // TODO: Add test cases.
}

func TestURLHandler_Get(t *testing.T) {
	log.Println("TestURLHandler_Get - Starting TestURLHandler_Get")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T))
	}
}
