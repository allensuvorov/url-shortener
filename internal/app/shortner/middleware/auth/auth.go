package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
)

var secretkey = []byte("secret key")
var idLength int = 4

// TODO inc9: task 1: AuthMiddleware()
// TODO inc9: task 2: file save/restore user history
/*
   <MW> - makes sure ID is in cookie
  	 |
  <hander> can get ID from cookie, why put a duplicate in header
*/

/* plan B
pass id via request header
*/

func checkID(r *http.Request) (string, bool) {
	cookieIdSign, err := r.Cookie("IdSign")
	if err == http.ErrNoCookie {
		return "", false
	}

	data, err := hex.DecodeString(cookieIdSign.Value)
	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:idLength])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[idLength:]) {
		id := hex.EncodeToString(data[:idLength])
		log.Println("auth/authenticate - clientExists - id:", id)

		return id, true
	} else {
		return "", false
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

func registerNewClient(w http.ResponseWriter, size int) (string, error) {
	rand, err := generateRandom(size)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write([]byte(rand))
	sign := h.Sum(nil)

	idSign := append(rand, sign...)
	stringIdSign := hex.EncodeToString(idSign)

	cookieIdSign := &http.Cookie{
		Name:  "idSign",
		Value: stringIdSign,
	}

	http.SetCookie(w, cookieIdSign)

	id := hex.EncodeToString(rand)
	log.Println("auth/registerNewClient - id:", id)

	return id, nil
}

// TODO AuthMiddleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Hello! ")
		var id string
		var ok bool
		var err error

		if id, ok = checkID(r); !ok {
			id, err = registerNewClient(w, idLength)
			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}

		}

		r.Header.Set("id", id)

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
