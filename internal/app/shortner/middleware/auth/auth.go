package auth

import "net/http"

//TODO generate client ID
//TODO generate signature for the client ID
//TODO read/write cookie

func AuthMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("registered")

		// no cookie
		if err == http.ErrNoCookie {
			// TODO generate ID with signature set it to cookie
			cookie = &http.Cookie{
				Name:  "registered",
				Value: "0",
			}
		}

		// cookie - no reg field
		// cookie, reg value, wrong key
		// cookie, reg value, write key

	})
}
