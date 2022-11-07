package routers

import (
	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/auth"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/compress"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/go-chi/chi/v5"
)

func NewRouter(url url.URLHandler) chi.Router {

	// create new router
	r := chi.NewRouter()

	r.Use(auth.AuthMiddleware, compress.GzipMiddleware)
	r.Get("/{hash}", url.Get)
	r.Post("/", url.Create)
	r.Post("/api/shorten", url.CreateForJSONClient)
	r.Get("/api/user/urls", url.GetClientActivity)
	r.Get("/ping", url.PingDB)
	return r
}
