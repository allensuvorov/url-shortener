package routers

import (
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/go-chi/chi/v5"
)

func NewRouter(url url.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{hash}", url.Middleware(url.Get))
	r.Post("/", url.Middleware(url.Create))
	r.Post("/api/shorten", url.Middleware(url.API))
	return r
}
