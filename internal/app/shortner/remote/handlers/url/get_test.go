package url_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/stretchr/testify/require"
)

type handlerGetTest struct {
	name  string
	input getInput
	want  getWant
}

type md struct { // mock data
	URL  string
	Hash string
}

type getInput struct {
	md         md
	requestURL string
}

type getWant struct {
	longURL    string
	StatusCode int
}

func newMockHandler(u, h string) handler.URLHandler {
	usm := storage.NewURLStorage()
	usm.Create(h, u)
	us := service.NewURLService(usm)
	uh := handler.NewURLHandler(us)
	return uh
}

func (st handlerGetTest) run(t *testing.T) { // subtest
	// request
	req, err := http.NewRequest("GET", st.input.requestURL, nil)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// response
	rec := httptest.NewRecorder()

	// new mock handler object
	uh := newMockHandler(st.input.md.URL, st.input.md.Hash)

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

var getTests = []handlerGetTest{
	{
		name: "1st Test Case - Positive: apple/store",
		input: getInput{
			md: md{
				URL:  "http://www.apple.com/store",
				Hash: "a7d59904",
			},
			requestURL: "http://localhost:8080/a7d59904",
		},
		want: getWant{
			longURL:    "http://www.apple.com/store",
			StatusCode: http.StatusTemporaryRedirect,
		},
	},
	{
		name: "2st Test Case - Negative: not found",
		input: getInput{
			md: md{
				URL:  "http://www.apple.com/store",
				Hash: "a7d59904",
			},
			requestURL: "http://localhost:8080/1111111",
		},
		want: getWant{
			longURL:    "",
			StatusCode: http.StatusBadRequest,
		},
	}, // TODO: Add test cases.
}

func TestURLHandler_Get(t *testing.T) {
	log.Println("TestURLHandler_Get - Starting TestURLHandler_Get")
	for _, tt := range getTests {
		t.Run(tt.name, tt.run)
	}
}
