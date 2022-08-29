package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

// map to store short urls and full urls
var urls map[string]string = make(map[string]string)

// sha256 to generate the hash value
func Shorten(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	// log.Println(fmt.Sprintf("%x", h.Sum(nil)))
	shortURL := fmt.Sprintf("%x", h.Sum(nil))[0:8]
	return shortURL
}

// CreateShortURL — обработчик запроса.
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if r.Method == "GET" {

			//
			base := path.Base(r.URL.Path)
			log.Println(base)

			// set header Location
			w.Header().Set("Location", "http:/localhost:8080/"+urls[base])

			// устанавливаем статус-код 307
			w.WriteHeader(http.StatusTemporaryRedirect)

			w.Write([]byte(urls[base]))

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

		shortURL := Shorten(string(b)) // shortened URL

		// add url to the map
		urls[shortURL] = string(b)
		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte("http:/localhost:8080/" + shortURL))
	}

}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", CreateShortURL)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
