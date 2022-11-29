package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/service"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func Test_CreateForJSONClient(t *testing.T) {
	config.BuildConfig()
	usm := storage.NewURLStorage()
	us := service.NewURLService(usm)
	uh := NewURLHandler(us)

	testCases := []struct {
		name                 string
		url                  string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			name:                 "Invalid URL",
			url:                  `{"url":"htt_1_p://google.com/"}`,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: []byte("Failed to create short URL\n" + `{"result":""}` + "\n"),
		},
		{
			name:                 "Created",
			url:                  `{"url":"http://www.apple.com/store"}`,
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: []byte(`{"result":"http://localhost:8080/a7d59904"}` + "\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBufferString(tc.url)
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", b)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			uh.CreateForJSONClient(w, r)
			log.Println("Test_CreateForJSONClient, body is:", w.Body.String())
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.Bytes())
		})
	}
}

func Test_shortener(t *testing.T) {
	config.BuildConfig()
	usm := storage.NewURLStorage()
	us := service.NewURLService(usm)
	uh := NewURLHandler(us)

	testCases := []struct {
		name                 string
		url                  string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			name:                 "Invalid URL",
			url:                  "htt_1_p://google.com/",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: []byte("Failed to create short URL\n"),
		},
		{
			name:                 "Created",
			url:                  "http://www.apple.com/store",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: []byte(`http://localhost:8080/a7d59904`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBufferString(tc.url)
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080", b)
			w := httptest.NewRecorder()

			uh.Create(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.Bytes())
		})
	}
}

func Test_expander(t *testing.T) {
	config.BuildConfig()
	usm := storage.NewURLStorage()
	ue := entity.URLEntity{
		Hash: "a7d59904",
		URL:  "http://www.apple.com/store",
	}
	usm.Create(ue)
	us := service.NewURLService(usm)
	uh := NewURLHandler(us)

	testCases := []struct {
		name                     string
		hash                     string
		expectedStatusCode       int
		expectedRedirectLocation string
	}{
		{
			name:                     "Not found",
			hash:                     "abcdefg",
			expectedStatusCode:       http.StatusBadRequest, // можно и NotFound вернуть
			expectedRedirectLocation: "",
		},
		{
			name:                     "Success",
			hash:                     "a7d59904",
			expectedStatusCode:       http.StatusTemporaryRedirect,
			expectedRedirectLocation: "http://www.apple.com/store",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://localhost:8080/"+tc.hash, nil)
			// создаем тестовый контекст
			ctx := chi.NewRouteContext()
			// передаем параметры в тестовый конекст
			ctx.URLParams.Add("hash", tc.hash)
			rctx := context.WithValue(r.Context(), chi.RouteCtxKey, ctx)
			r = r.WithContext(rctx)

			w := httptest.NewRecorder()

			uh.Get(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRedirectLocation, w.Header().Get("Location"))
		})
	}
}
