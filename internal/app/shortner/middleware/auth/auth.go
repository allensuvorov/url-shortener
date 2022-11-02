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
func registerNewClient(size int) (uint32, error) {

	// TODO set ID and sgn to cookie
	//

	rand, err := generateRandom(size)
	if err != nil {
		return 0, err
	}

	id := binary.BigEndian.Uint32(rand)

	h := hmac.New(sha256.New, secretkey)
	h.Write([]byte(rand))
	sign := h.Sum(nil)

	//TODO generate signature for the client ID

	idSign := append(rand, sign...)
	stringIdSign := hex.EncodeToString(idSign)

	cookieIdSign := &http.Cookie{
		Name:  "idSign",
		Value: stringIdSign,
	}

	http.SetCookie

	//TODO read/write cookie
	return id, nil
}

// TODO AuthMiddleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if id, ok := authenticate(r); !ok {
			id, err := registerNewClient(4)

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
