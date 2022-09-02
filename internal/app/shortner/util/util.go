package util

import (
	"crypto/sha256"
	"fmt"
	"log"

	"yandex/projects/urlshortner/internal/app/shortner/storage"
)

// sha256 to generate the hash value
func Shorten(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	// log.Println(fmt.Sprintf("%x", h.Sum(nil)))

	hash := fmt.Sprintf("%x", h.Sum(nil))
	var shortURL string

	// check if short URL is already in the map for a different long url, expand hash slice till unique
	for i := 8; i < len(hash); i++ {
		shortURL = hash[0:i]
		if v, ok := storage.Urls[shortURL]; !(ok && v != s) {
			break
		}
	}
	log.Println("created new shortURL", shortURL)
	return shortURL
}
