package url

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/go-chi/chi/v5"
)

type URLService interface {
	// Create takes URL and returns hash.
	Create(u string) (string, error)
	// Get takes hash and returns URL.
	Get(h string) (string, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(us URLService) URLHandler {
	return URLHandler{
		urlService: us,
	}
}

func (uh URLHandler) API(w http.ResponseWriter, r *http.Request) {
	// целевой объект
	var v1 struct {
		URL string
	}

	if err := json.NewDecoder(r.Body).Decode(&v1); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("API handler - request: object", v1)
	log.Println("API handler - URL in the request is", v1.URL)

	shortURL, err := uh.urlService.Create(v1.URL)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	v2 := struct {
		Result string `json:"result"`
	}{
		Result: shortURL,
	}
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
	encoder.Encode(v2)

	log.Println("API handler - v2 is", v2.Result)
	log.Println("API handler - buf is", buf.String())

	// пишем тело ответа
	w.Write(buf.Bytes())
}

// Create passes URL to service and returns response with Hash.
func (uh URLHandler) Create(w http.ResponseWriter, r *http.Request) {

	// читаем Body
	b, err := io.ReadAll(r.Body)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert to string
	u := string(b)

	// log body from request
	log.Println("Create Handler - URL in the POST request is", u)

	shortURL, err := uh.urlService.Create(u)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	w.Write([]byte(shortURL))
}

// Get passes Hash to service and returns response with URL.
func (uh URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	// get hash
	h := chi.URLParam(r, "hash")

	// log path and hash
	log.Println("Handler Get", h, r.URL.Path)

	// Get from storage
	u, err := uh.urlService.Get(h)

	if err == errors.ErrNotFound {
		http.Error(w, "URL does not exist", http.StatusBadRequest)
		return
	}

	// set header Location
	w.Header().Set("Location", u)

	// устанавливаем статус-код 307
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))
}
