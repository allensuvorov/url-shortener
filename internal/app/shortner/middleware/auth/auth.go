package auth

import (
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/middleware/crypto"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Hello! ")
		var id string
		var ok bool
		var err error

		if id, ok = crypto.CheckID(r); !ok {
			r.Header.Set("auth", "false")
			log.Println("auth/AuthMiddleware, ID:", id)
			id, err = crypto.RegisterNewClient(w, crypto.IDLength)
			if err != nil {
				log.Printf("failed to register new client: %v", err)
			}
		} else {
			r.Header.Set("auth", "true")
		}

		r.Header.Set("id", id)

		next.ServeHTTP(w, r)
		log.Println("AuthMiddleware: Bye! ")

	})
}
