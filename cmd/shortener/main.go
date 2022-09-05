package main

import (
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/handlers"
)

func main() {
	r := chi.NewRouter()
	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.Multiplexer)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
