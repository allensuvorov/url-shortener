package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
)

// TODO generate client ID
func generateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateID(size int) (string, error) {
	rand, err := generateRandom(size)

	if err != nil {
		return "", err
	}

	encoded := hex.EncodeToString(rand)
	return encoded, nil
}

//TODO generate signature for the client ID
//TODO read/write cookie

func AuthMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("id")

		// no cookie
		if err == http.ErrNoCookie {
			// TODO generate ID with signature set it to cookie
			id, err := generateID(16)

			if err != nil {
				log.Printf("failed decompress data: %v", err)
			}

			cookie = &http.Cookie{
				Name:  "id",
				Value: id,
			}
		}

		// cookie - no reg field
		// cookie, reg value, wrong key
		// cookie, reg value, write key

	})
}
