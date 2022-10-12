package routers

import (
	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/compress"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/go-chi/chi/v5"
)

func NewRouter(url url.URLHandler) chi.Router {
	// create middleware handler
	mw := compress.GzipHandler{}

	// create new router
	r := chi.NewRouter()

	r.Use(mw.GzipMiddleware)
	r.Get("/{hash}", url.Get)
	r.Post("/", url.Create)
	r.Post("/api/shorten", url.API)
	return r
}
