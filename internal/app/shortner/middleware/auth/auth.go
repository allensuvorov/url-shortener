package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

var secretkey = []byte("secret key")

// TODO inc9: task 1: AuthMiddleware()
// TODO inc9: task 2: file save/restore user history

func authenticate(r *http.Request) (uint32, bool) {
	cookieIdSign, err := r.Cookie("IdSign")
	if err == http.ErrNoCookie {
		return 0, false
	}

	data, err := hex.DecodeString(cookieIdSign.Value)
	if err != nil {
		panic(err)
	}

	id := binary.BigEndian.Uint32(data[:4])

	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		log.Println("auth/clientExists - id:", id)
		return id, true
	} else {
		return 0, false
	}
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

		if id, ok := authenticate(r); !ok {
			err := registerNewClient(16)

			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}
		} else {
			log.Println(id)
			// TODO if all good, then authed = true, id = id
			// TODO log: ID - hash

		}

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
