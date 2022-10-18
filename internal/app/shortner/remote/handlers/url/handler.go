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

func (uh URLHandler) CreateForJSONClient(w http.ResponseWriter, r *http.Request) {
	// целевой объект
	var dv struct { // decoded value
		URL string
	}

	// TODO: Read and handle content-type header from request
	// contentType := response.Header.Get("Content-Type")
	// это может быть, например, "application/json; charset=UTF-8"

	if err := json.NewDecoder(r.Body).Decode(&dv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Handler/API - request: object", dv)
	log.Println("Handler/API - URL in the request is", dv.URL)

	shortURL, err := uh.urlService.Create(dv.URL)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	ev := struct { // encoded value
		Result string `json:"result"`
	}{
		Result: shortURL,
	}
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false) // без этой опции символ '&' будет заменён на "\u0026"
	encoder.Encode(ev)

	log.Println("API handler - v2 is", ev.Result)
	log.Println("API handler - buf is", buf.String())

	// сначала устанавливаем заголовок Content-Type
	// для передачи клиенту информации, кодированной в JSON
	w.Header().Set("content-type", "application/json")

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	w.Write(buf.Bytes())
}

// Create passes URL to service and returns response with Hash.
func (uh URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler/Create - Start")
	// читаем Body
	b, err := io.ReadAll(r.Body)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	log.Println("Handler/Create - read r.Body, no err")

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
