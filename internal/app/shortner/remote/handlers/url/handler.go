package url

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/go-chi/chi/v5"
)

type URLService interface {
	// Create takes URL and returns hash.
	Create(ue entity.DTO) (string, error)
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
	var decVal struct { // decoded value
		URL string
	}

	// TODO: Read and handle content-type header from request
	// contentType := response.Header.Get("Content-Type")
	// это может быть, например, "application/json; charset=UTF-8"

	if err := json.NewDecoder(r.Body).Decode(&decVal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Handler/CreateForJSONClient - request: object", decVal)
	log.Println("Handler/CreateForJSONClient - URL in the request is", decVal.URL)

	ue := entity.DTO{
		URL: decVal.URL,
	}

	shortURL, err := uh.urlService.Create(ue)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	encVal := struct { // encoded value
		Result string `json:"result"`
	}{
		Result: shortURL,
	}

	log.Println("Handler/CreateForJSONClient: ev is", encVal.Result)

	// сначала устанавливаем заголовок Content-Type
	// для передачи клиенту информации, кодированной в JSON
	w.Header().Set("content-type", "application/json")

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	json.NewEncoder(w).Encode(encVal)
}

// Create passes URL to service and returns response with Hash.
func (uh URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler/Create - Start")

	log.Println("Handler/Create - Set-Cookie:", w.Header().Get("Set-Cookie"))
	log.Println("Handler/Create - ID header:", r.Header.Get("id"))

	ue := entity.DTO{
		ClientID: r.Header.Get("id"),
	}

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
	ue.URL = string(b)

	// log body from request
	log.Println("Create Handler - URL in the POST request is", ue.URL)

	ue.Hash, err = uh.urlService.Create(ue)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	w.Write([]byte(ue.Hash))
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
