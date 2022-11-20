package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
)

var secretkey = []byte("secret key")
var IdLength int = 4

func CheckID(r *http.Request) (string, bool) {
	cookieIDSign, err := r.Cookie("idSign")

	if err == http.ErrNoCookie {
		log.Println("auth/CheckID, no idSign in cookie.")
		return "", false
	}
	log.Println("auth/CheckID, ID from cookie:", cookieIDSign.Value)

	data, err := hex.DecodeString(cookieIDSign.Value)
	if err != nil {
		panic(err)
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:IdLength])
	sign := h.Sum(nil)

	// log.Println("auth/CheckID, ID:", hex.EncodeToString(data[:IdLength]))

	if hmac.Equal(sign, data[IdLength:]) {
		id := hex.EncodeToString(data[:IdLength])
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

func RegisterNewClient(w http.ResponseWriter, size int) (string, error) {
	rand, err := generateRandom(size)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, secretkey)
	h.Write([]byte(rand))
	sign := h.Sum(nil)

	idSign := append(rand, sign...)
	stringIDSign := hex.EncodeToString(idSign)

	cookieIDSign := &http.Cookie{
		Name:  "idSign",
		Value: stringIDSign,
	}

	http.SetCookie(w, cookieIDSign)

	id := hex.EncodeToString(rand)
	log.Println("auth/RegisterNewClient - id:", id)

	return id, nil
}
