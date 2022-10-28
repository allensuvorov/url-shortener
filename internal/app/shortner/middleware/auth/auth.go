package auth

import (
	"crypto/rand"
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

//TODO generate signature for the client ID
//TODO read/write cookie

const registered = "registered"

func AuthMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(registered)

		// no cookie
		if err == http.ErrNoCookie {
			// TODO generate ID with signature set it to cookie
			cookie = &http.Cookie{
				Name:  registered,
				Value: "0",
			}
		}

		// cookie - no reg field
		// cookie, reg value, wrong key
		// cookie, reg value, write key

	})
}
