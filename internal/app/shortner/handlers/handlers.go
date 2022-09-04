package handlers

import (
	"io"
	"net/http"
	"path"

	// "yandex/projects/urlshortner/internal/app/shortner/storage"
	// "yandex/projects/urlshortner/internal/app/shortner/util"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/util"
)

// Multiplexer - is a request router.
func Multiplexer(w http.ResponseWriter, r *http.Request) {
	//log.Println("path is /", r.URL.Path)

	// проверяем, каким методом получили запрос
	switch r.Method {
	// если методом POST
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	default:
		return
	}
}

// getHandler - handles GET requests.
func getHandler(w http.ResponseWriter, r *http.Request) {
	// get hash - the part after last slash
	base := path.Base(r.URL.Path)
	// log.Println("Path after last slash", base)
	// log.Println(r.URL, r.URL.Path)

	// check if hash exists
	if !storage.HashExists(base) {
		http.Error(w, "URL does not exist", 400)
		return
	}

	u := storage.GetURL(base)
	// set header Location
	w.Header().Set("Location", u)

	// устанавливаем статус-код 307
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))

	return
}

// postHandler - handles POST requests.
func postHandler(w http.ResponseWriter, r *http.Request) {
	// читаем Body
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// check if long URL is empty string
	if len(b) == 0 {
		http.Error(w, "empty URL", 400)
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
