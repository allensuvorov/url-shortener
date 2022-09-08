package handlers

import (
	"io"
	"net/http"
	"path"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/util"
)

// GetURL - takes short URL and returns full URL.
func GetURL(w http.ResponseWriter, r *http.Request) {
	// get hash - last element of path (after slash)
	base := path.Base(r.URL.Path)

	// check if hash exists, if not - return 400
	if !storage.HashExists(base) {
		http.Error(w, "URL does not exist", http.StatusBadRequest)
		return
	}

	// getURL from storage
	u := storage.GetURL(base)

	// set header Location
	w.Header().Set("Location", u)

	// устанавливаем статус-код 307
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))
}

// PostURL - takes full URL and returns short URL.
func PostURL(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if long URL is empty string
	if len(b) == 0 {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}
	// // check it URL is valid
	// u, err := url.ParseRequestURI(string(b))
	// log.Println("parsed URL", u)

	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }

	// check if long url is already in the map
	exists, h := storage.URLExists(string(b))

	if !exists {

		// generate shortened URL
		h = util.Shorten(string(b))

		// add url to the map
		storage.NewURL(h, string(b))
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	shortURL := "http://localhost:8080/" + h

	// пишем тело ответа
	w.Write([]byte(shortURL))

}
