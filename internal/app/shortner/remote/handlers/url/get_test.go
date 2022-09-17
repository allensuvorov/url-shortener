package url

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/stretchr/testify/require"
)

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

func newMockHandler(u, h string) URLHandler {
	ue := &entity.URLEntity{
		URL:  u,
		Hash: h,
	}
	usm := storage.NewURLStorage()
	usm.InMemory = append(usm.InMemory, ue)
	us := service.NewURLService(usm)
	uh := NewURLHandler(us)
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
