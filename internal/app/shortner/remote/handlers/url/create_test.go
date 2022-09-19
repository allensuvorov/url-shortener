package url_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
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
	responseBody string
	statusCode   int
}

// new create hander
var ch = handler.NewURLHandler(service.NewURLService(storage.NewURLStorage()))

func (st handlerCreateTest) run(t *testing.T) {
	// Get request
	bodyReader := bytes.NewReader([]byte(st.input.longURL))
	req, err := http.NewRequest("POST", st.input.requestURL, bodyReader)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// New response recorder
	rec := httptest.NewRecorder()

	// Run the handler
	ch.Create(rec, req)

	// Get response
	res := rec.Result()
	defer res.Body.Close()

	// Check response status code
	require.Equal(t, res.StatusCode, st.want.statusCode, "expected status Created, got other")

	// Check response body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read respons: %v", err)
	}
	require.Equal(t, st.want.responseBody, string(b), "short URL is not matching")
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
			responseBody: "http://localhost:8080/a7d59904",
			statusCode:   http.StatusCreated,
		},
	},
	{
		name: "2st Test Case - Negative: invalide long URL",
		input: createInput{
			longURL:    "123http://www.apple.com/store",
			requestURL: "localhost:8080/",
		},
		want: createWant{
			responseBody: "Failed to create short URL\n",
			statusCode:   http.StatusInternalServerError,
		},
	}, // TODO: Add test cases.
}

func TestURLHandler_Create(t *testing.T) {
	for _, tt := range createTests {
		t.Run(tt.name, tt.run)
	}
}
