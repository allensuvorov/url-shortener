package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func registerNewClient(size int) error {

	// TODO generate ID, key and signature
	// TODO set ID and sgn to cookie
	// TODO save ID and key serverside

	rand, err := generateRandom(size)

	if err != nil {
		return "", err
	}

	id := hex.EncodeToString(rand)

	cookieID := &http.Cookie{
		Name:  "id",
		Value: id,
	}

	// создаём случайный ключ
	key, err := generateRandom(16)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(id))
	sg := h.Sum(nil)

	cookieSG := &http.Cookie{
		Name:  "signature",
		Value: string(sg),
	}
	return nil
}

//TODO generate signature for the client ID
//TODO read/write cookie

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieID, err := r.Cookie("id")

		// If cookieID in storage, but wrong sgn

		// if no cookieID, OR cookieID not is storage OR wrong signature
		if err == http.ErrNoCookie {

			id, err := registerNewClient(16)

			if err != nil {
				log.Printf("failed decompress data: %v", err)
			}

		}

		// if all good, then authed = true
		// then

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
