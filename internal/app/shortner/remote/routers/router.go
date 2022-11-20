package routers

import (
	"github.com/go-chi/chi/v5"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/auth"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/compress"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers"
)

func NewRouter(url handlers.URLHandler) chi.Router {

	r := chi.NewRouter()

	r.Use(auth.AuthMiddleware, compress.GzipMiddleware)
	r.Get("/{hash}", url.Get)
	r.Post("/", url.Create)
	r.Post("/api/shorten", url.CreateForJSONClient)
	r.Get("/api/user/urls", url.GetClientActivity)
	r.Get("/ping", url.PingDB)
	r.Post("/api/shorten/batch", url.BatchCreate)
	return r
}
