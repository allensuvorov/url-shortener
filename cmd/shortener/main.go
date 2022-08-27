package main

import (
	"io"
	"net/http"
)

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(b)
	} else {
		w.Write([]byte("<h1>Hello, World</h1>"))
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HelloWorld)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
