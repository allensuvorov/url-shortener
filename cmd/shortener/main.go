package main

import (
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/config"
	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/routers"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func main() {
	log.Println("Main: now will try to Config/BuildConfig")
	config.BuildConfig()
	URLStorage := storage.NewURLStorage()
	URLService := service.NewURLService(URLStorage)
	URLHandler := handler.NewURLHandler(URLService)
	r := routers.NewRouter(URLHandler)
	log.Println("Main: now will try to get sa from config")
	// sa := config.UC.SA // server address
	defaultSA := ":8080"
	sa := &defaultSA
	log.Println("Serving on port", *sa)
	log.Fatal(http.ListenAndServe(*sa, r))
}
