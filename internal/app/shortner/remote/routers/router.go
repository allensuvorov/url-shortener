package routers

import "github.com/go-chi/chi/v5"

func NewRouter(url url.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{hash}", url.GetUrl)
	r.Post("/", url.CreateHash)
	return r
}
