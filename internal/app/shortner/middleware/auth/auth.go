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

var secretkey = []byte("secret key")

// TODO inc9: task 1: AuthMiddleware()
// TODO inc9: task 2: file save/restore user history

// TODO clientExists
func clientExists(r *http.Request) bool {
	// if no cookie, OR wrong signature
	cookieID, err := r.Cookie("id")
	if err == http.ErrNoCookie {
		return false
	}

	cookieSG, err := r.Cookie("signature")
	if err == http.ErrNoCookie {
		return false
	}

	// check signature
	id, err := hex.DecodeString(cookieID.Value)
	if err != nil {
		panic(err)
	}

	sg, err := hex.DecodeString(cookieSG.Value)
	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write(id)
	sign := h.Sum(nil)

	if hmac.Equal(sign, sg) {
		return true
	} else {
		return false
	}

	// return true

}

func generateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// TODO registerNewClient
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

// TODO AuthMiddleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !clientExists(r) {
			err := registerNewClient(16)

			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}
		} else {

			// TODO if all good, then authed = true, id = id
			// TODO log: ID - hash

		}

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
