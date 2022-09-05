package main

import (
	"fmt"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", handlers.PostHandler)
	r.Get("/", handlers.GetHandler)

	fmt.Println("Serving on port: 8080")
	http.ListenAndServe(":8080", r)

	// маршрутизация запросов обработчику
	// http.HandleFunc("/", handlers.Multiplexer)
	// // запуск сервера с адресом localhost, порт 8080
	// log.Fatal(http.ListenAndServe(":8080", nil)) // log.Fatal will print errors if server crashes
}
