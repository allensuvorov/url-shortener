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
	h.Write([]byte(s)) // what's that?

	hash := fmt.Sprintf("%x", h.Sum(nil))
	sh := getUniqShortHash(hash, s)
	log.Println("created new shortURL", sh)
	return sh
}

// check if short URL (hash) is already in the map for a different long url, expand hash slice till unique
func getUniqShortHash(h string, s string) string {
	var sh string // short hash
	for i := 8; i < len(h); i++ {
		sh = h[0:i]
		if !storage.HashExists(sh) {
			return sh
		}
		if storage.GetURL(sh) == s {
			return sh
		}
	}
	return sh
}
