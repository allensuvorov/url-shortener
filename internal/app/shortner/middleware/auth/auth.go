package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net/http"
)

var secretkey = []byte("secret key")

// TODO inc9: task 1: AuthMiddleware()
// TODO inc9: task 2: file save/restore user history
/*
   <MW> - makes sure ID is in cookie
  	 |
  <hander> can get ID from cookie, why put a duplicate in header
*/
func authenticate(r *http.Request) bool {
	cookieIdSign, err := r.Cookie("IdSign")
	if err == http.ErrNoCookie {
		return false
	}

	data, err := hex.DecodeString(cookieIdSign.Value)
	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		log.Println("auth/clientExists - id:")
		return true
	} else {
		return false
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

func registerNewClient(w http.ResponseWriter, size int) (uint32, error) {
	rand, err := generateRandom(size)
	if err != nil {
		return 0, err
	}

	id := binary.BigEndian.Uint32(rand)

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
	return id, nil
}

// TODO AuthMiddleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !authenticate(r) {
			id, err := registerNewClient(w, 4)

			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}
			// how to pass id to handler?
			r.Header.Set("id", string(id))

		} else {
			log.Println(id)
			// TODO if all good, then authed = true, id = id
			// TODO log: ID - hash

		}

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
