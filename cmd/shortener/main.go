package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/routers"
	service "github.com/allensuvorov/urlshortner/internal/app/shortner/service/url"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

func init() {
	log.Println("Config/getConfigFromCLI, passed flag: a,b,f")
	flag.StringVar(&config.UC.SA, "a", "", "SERVER_ADDRESS")
	flag.StringVar(&config.UC.BU, "b", "", "BASE_URL")
	flag.StringVar(&config.UC.FSP, "f", "", "FILE_STORAGE_PATH")
	flag.StringVar(&config.UC.DSN, "d", "", "DATABASE_DSN")
}

func main() {
	// for testing:
	// os.Setenv("FILE_STORAGE_PATH", "/Users/allen/go/src/yandex/projects/urlshortner/internal/app/shortner/storage/.urls.log")
	//os.Setenv("DATABASE_DSN", "postgres://postgres:sql@localhost:5432/url_db")
	config.BuildConfig()
	var URLStorage service.URLStorage
	if config.UC.DSN != "" {
		URLStorage = storage.NewURLDB()
	} else {
		URLStorage = storage.NewURLStorage()
	}
	URLService := service.NewURLService(URLStorage)
	URLHandler := handlers.NewURLHandler(URLService)
	r := routers.NewRouter(URLHandler)
	sa := config.UC.SA // server address from config
	log.Println("Serving on port", sa)
	log.Fatal(http.ListenAndServe(sa, r))
}
