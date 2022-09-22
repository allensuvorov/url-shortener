package url

import (
	"crypto/sha256"
	"fmt"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// sha256 to generate the hash value
func BuildHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s)) // what's that?

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash
}

// check if short URL (hash) is already in the DB for a different long url, expand hash slice till unique
func getUniqShortHash(h string, u string, us URLService) string {
	var sh string // short hash
	for i := 8; i < len(h); i++ {
		sh = h[0:i]

		u1, err := us.urlStorage.GetURLByHash(sh)

		// if sh is uniq (not in storage), return sh
		if err == errors.ErrNotFound {
			return sh
		}
		// check it the URL is different
		if u1 == u {
			return sh
		}
	}
	return sh
}
