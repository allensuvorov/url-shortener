package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/config"
	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/routers"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

var (
	port *string
	bu   *string
	fsp  *string
)

func init() {
	// get server address from CLI flag
	port = flag.String("a", ":8080", "SERVER_ADDRESS")
	// bu = flag.String("b", ":8080", "SERVER_ADDRESS")
}

func main() {
	flag.Parse()
	log.Println("Port from Flag is", *port)
	// get server address from local env if not in cli flags
	if len(*port) == 0 {
		sa, ok := os.LookupEnv("SERVER_ADDRESS")
		if !ok {
			log.Printf("%s not set\n; passing default", "SERVER_ADDRESS")
			sa = ":8080"
		}
		*port = sa
	}

	URLStorage := storage.NewURLStorage()
	URLConfig := config.NewURLConfig()
	URLService := service.NewURLService(URLStorage, URLConfig)
	URLHandler := handler.NewURLHandler(URLService)
	r := routers.NewRouter(URLHandler)
	log.Println("Serving on port", *port)
	log.Fatal(http.ListenAndServe(*port, r))
}
