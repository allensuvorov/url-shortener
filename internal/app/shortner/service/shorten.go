package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

func generateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

func getUniqShortHash(h string, u string, us URLService) string {
	var shortHash string
	for i := 8; i < len(h); i++ {
		shortHash = h[0:i]
		u1, err := us.urlStorage.GetURLByHash(shortHash)

		if err == errors.ErrNotFound {
			return shortHash
		}
		if u1 == u {
			return shortHash
		}
	}
	return shortHash
}
