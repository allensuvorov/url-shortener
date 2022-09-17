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
			name: "1st Test Case - Positive: apple/store",
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
			name: "2st Test Case - Negative: invalide long URL",
			fields: fields{
				urlService: service.NewURLService(storage.NewURLStorage()),
			},
			args: args{
				longURL:    "123http://www.apple.com/store",
				requestURL: "localhost:8080/",
				shortURL:   "Failed to create short URL\n",
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
			require.Equal(t, res.StatusCode, tt.args.StatusCode, "expected status Created, got other")

			// Check response body
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read respons: %v", err)
			}
			require.Equal(t, tt.args.shortURL, string(b), "short URL is not matching")

		})
	}
}

// data for handlerGetTest
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
	longURL    string
	StatusCode int
}

func (st handlerGetTest) run(t *testing.T) { // subtest
	// request
	req, err := http.NewRequest("GET", st.input.requestURL, nil)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// response
	rec := httptest.NewRecorder()

	// handler object
	// New url entity
	ue := &entity.URLEntity{
		URL:  st.input.md.URL,
		Hash: st.input.md.Hash,
	}
	usm := storage.NewURLStorage()
	usm.InMemory = append(usm.InMemory, ue)
	us := service.NewURLService(usm)
	uh := NewURLHandler(us)

	// Run the handler
	uh.Get(rec, req)

	// Check response
	res := rec.Result()
	defer res.Body.Close()

	// Check response status code
	require.Equal(t, st.want.StatusCode, res.StatusCode, "response status does not match expected")

	// Check response header
	require.Equal(t, st.want.longURL, res.Header.Get("Location"), "response header location does not match expected")
}

var tests = []handlerGetTest{
	{
		name: "1st Test Case - Positive: apple/store",
		input: input{
			md: md{
				URL:  "http://www.apple.com/store",
				Hash: "a7d59904",
			},
			requestURL: "http://localhost:8080/a7d59904",
		},
		want: want{
			longURL:    "http://www.apple.com/store",
			StatusCode: http.StatusTemporaryRedirect,
		},
	},
	{
		name: "2st Test Case - Negative: not found",
		input: input{
			md: md{
				URL:  "http://www.apple.com/store",
				Hash: "a7d59904",
			},
			requestURL: "http://localhost:8080/1111111",
		},
		want: want{
			longURL:    "",
			StatusCode: http.StatusBadRequest,
		},
	}, // TODO: Add test cases.
}

func TestURLHandler_Get(t *testing.T) {
	log.Println("TestURLHandler_Get - Starting TestURLHandler_Get")
	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
