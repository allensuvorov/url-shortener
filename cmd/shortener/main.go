package main

import (
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", handlers.PostURL)
	r.Get("/{hash}", handlers.GetURL)

	log.Println("Serving on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
