package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/routers"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func main() {
	// Add vars to env
	// Commands to add vars to env
	// export SERVER_ADDRESS=:8080
	// export BASE_URL=http://localhost:8080/
	sa, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		log.Printf("%s not set\n; passing default", "SERVER_ADDRESS")
		sa = ":8080"
	}

	URLStorage := storage.NewURLStorage()
	URLService := service.NewURLService(URLStorage)
	URLHandler := handler.NewURLHandler(URLService)
	r := routers.NewRouter(URLHandler)
	log.Println("Serving on port", sa)
	log.Fatal(http.ListenAndServe(sa, r))
}
