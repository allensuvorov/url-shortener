package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"yandex/projects/urlshortner/internal/app/shortner/storage"
	"yandex/projects/urlshortner/internal/app/shortner/util"
)

// CreateShortURL — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("path is /", r.URL.Path)
		if r.Method == "GET" {

			// get part after last slash
			base := path.Base(r.URL.Path)
			log.Println("after last slash", base)
			log.Println(r.URL, r.URL.Path)

			// check if hash exists
			if _, ok := storage.Urls[base]; !ok {
				http.Error(w, "URL does not exist", 400)
				return
			}

			// set header Location
			w.Header().Set("Location", storage.Urls[base])

			// устанавливаем статус-код 307
			w.WriteHeader(http.StatusTemporaryRedirect)

			w.Write([]byte(storage.Urls[base]))

			return
		}
	}
	if r.Method == "POST" {
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
		// check it URL is valid
		u, err := url.ParseRequestURI(string(b))
		log.Println("parsed URL", u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			//log.Printf("hi/there?: err=%+v url=%+v\n", err, u)
			return
		}

		// check if long url is already in the map
		var exists bool
		var shortURL string
		for k, v := range storage.Urls {
			if v == string(b) {
				exists = true
				shortURL = k
				break
			}
		}

		if !exists {

			// get shortened URL
			shortURL = util.Shorten(string(b))

			// add url to the map
			storage.Urls[shortURL] = string(b)
		}

		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)

		shortURL = "http://localhost:8080/" + shortURL

		// пишем тело ответа
		w.Write([]byte(shortURL))
	}

}
