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

	hash := fmt.Sprintf("%x", h.Sum(nil))
	var shortURL string

	// check if short URL (hash) is already in the map for a different long url, expand hash slice till unique
	for i := 8; i < len(hash); i++ {
		shortURL = hash[0:i]
		exists := storage.HashExists(shortURL)
		u := storage.GetURL(shortURL)
		if exists && u != s {
			break
		}
	}
	log.Println("created new shortURL", shortURL)
	return shortURL
}
