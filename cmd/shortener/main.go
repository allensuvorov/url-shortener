package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

// map to store short urls and full urls
var urls map[string]string = make(map[string]string)

// sha256 to generate the hash value
func Shorten(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	// log.Println(fmt.Sprintf("%x", h.Sum(nil)))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	var shortURL string

	// check if short URL is already in the map for a different long url, expand hash slice till unique
	for i := 8; i < len(hash); i++ {
		shortURL = hash[0:i]
		if v, ok := urls[shortURL]; !(ok && v != s) {
			break
		}
	}
	log.Println("created new shortURL", shortURL)
	return shortURL
}

// CreateShortURL — обработчик запроса.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if r.Method == "GET" {

			// get part after last slash
			base := path.Base(r.URL.Path)
			log.Println("after last slash", base)

			if _, ok := urls[base]; !ok {
				http.Error(w, "URL does not exist", 400)
				return
			}

			// set header Location
			w.Header().Set("Location", urls[base])

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
		var urlIsNew = true
		var shortURL string
		for k, v := range urls {
			if v == string(b) {
				urlIsNew = false
				shortURL = k
			}
		}

		if urlIsNew {

			// get shortened URL
			shortURL = Shorten(string(b))

			// add url to the map
			urls[shortURL] = string(b)
		}

		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)

		// пишем тело ответа
		w.Write([]byte(shortURL))
	}

}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", Handler)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
