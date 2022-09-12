package url

import (
	"crypto/sha256"
	"fmt"
	"log"

	//"yandex/projects/urlshortner/internal/app/shortner/storage"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/storage"
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
func getUniqShortHash(h string, u string) string {
	var sh string // short hash
	for i := 8; i < len(h); i++ {
		sh = h[0:i]

		u1, ok := storage.GetURL(sh)

		// if sh is uniq (not in storage), return sh
		if !ok {
			return sh
		}
		// check it the URL is different
		if u1 == u {
			return sh
		}
	}
	return sh
}
