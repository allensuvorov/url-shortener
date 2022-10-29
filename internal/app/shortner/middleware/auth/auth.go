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

	//TODO generate signature for the client ID
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
	//TODO read/write cookie
	return nil
}

func clientExists(cookieID *http.Cookie, err error) bool {
	// if no cookieID, OR cookieID not is storage OR wrong signature
	if err == http.ErrNoCookie {
		return false
	}
	cookieID.Value
	storage
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieID, err := r.Cookie("id")

		if !clientExists(cookieID, err) {
			err := registerNewClient(16)

			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}
		} else {

			// if all good, then authed = true, id = id

		}

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
