package handlers

import (
	"log"
	"net/http"
	"path"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
)

// GetURL - takes short URL and returns full URL.
func GetURL(w http.ResponseWriter, r *http.Request) {
	// get hash - last element of path (after slash)
	base := path.Base(r.URL.Path)

	// log path and hash
	log.Println(base, r.URL.Path)

	// check if hash exists, if not - return 400
	if !storage.HashExists(base) {
		http.Error(w, "URL does not exist", http.StatusBadRequest)
		return
	}

	// getURL from storage
	u := storage.GetURL(base)

	// set header Location
	w.Header().Set("Location", u)

	// устанавливаем статус-код 307
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))
}
