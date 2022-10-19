package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	handler "github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/routers"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func init() {
	log.Println("Config/getConfigFromCLI, passed flag: a,b,f")
	flag.StringVar(&config.UC.SA, "a", config.DefaultSA, "SERVER_ADDRESS")
	flag.StringVar(&config.UC.BU, "b", config.DefaultBU, "BASE_URL")
	flag.StringVar(&config.UC.FSP, "f", config.DefaultFSP, "FILE_STORAGE_PATH")
}

func main() {
	config.BuildConfig()
	URLStorage := storage.NewURLStorage()
	URLService := service.NewURLService(URLStorage)
	URLHandler := handler.NewURLHandler(URLService)
	r := routers.NewRouter(URLHandler)
	sa := config.UC.SA // server address from config
	log.Println("Serving on port", sa)
	log.Fatal(http.ListenAndServe(sa, r))
}
