package url

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/stretchr/testify/require"
)

type handlerCreateTest struct {
	name  string
	input createInput
	want  createWant
}

type createInput struct {
	longURL    string
	requestURL string
}

type createWant struct {
	shortURL   string
	StatusCode int
}

// new create hander
var ch = NewURLHandler(service.NewURLService(storage.NewURLStorage()))

func (st handlerCreateTest) run(t *testing.T) {
	jsonBody := []byte(st.input.longURL)
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := st.input.requestURL
	req, err := http.NewRequest("POST", requestURL, bodyReader)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	// response
	rec := httptest.NewRecorder()

	// Run the handler
	ch.Create(rec, req)

	// Check response
	res := rec.Result()
	defer res.Body.Close()
	// Check response status code
	require.Equal(t, res.StatusCode, st.want.StatusCode, "expected status Created, got other")

	// Check response body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read respons: %v", err)
	}
	require.Equal(t, st.want.shortURL, string(b), "short URL is not matching")
}

// test data

var createTests = []handlerCreateTest{
	{
		name: "1st Test Case - Positive: apple/store",
		input: createInput{
			longURL:    "http://www.apple.com/store",
			requestURL: "localhost:8080/",
		},
		want: createWant{
			shortURL:   "http://localhost:8080/a7d59904",
			StatusCode: http.StatusCreated,
		},
	},
	{
		name: "2st Test Case - Negative: invalide long URL",
		input: createInput{
			longURL:    "123http://www.apple.com/store",
			requestURL: "localhost:8080/",
		},
		want: createWant{
			shortURL:   "Failed to create short URL\n",
			StatusCode: http.StatusInternalServerError,
		},
	}, // TODO: Add test cases.
}

func TestURLHandler_Create(t *testing.T) {
	for _, tt := range getTests {
		t.Run(tt.name, tt.run)
	}
}
