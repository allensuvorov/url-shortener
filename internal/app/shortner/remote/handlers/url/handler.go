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

	GetClientActivity(id string) ([]entity.DTO, error)

	PingDB() bool
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

	ue.ClientID = r.Header.Get("id")

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

func (uh URLHandler) GetClientActivity(w http.ResponseWriter, r *http.Request) {
	// auth false
	if r.Header.Get("auth") == "false" {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
	}

	if r.Header.Get("auth") == "true" {
		clientID := r.Header.Get("id")
		dtoList, err := uh.urlService.GetClientActivity(clientID)
		if err != nil {
			http.Error(w, "Failed to get client activity", http.StatusInternalServerError)
		}

		log.Println("Handler/GetClientActivity: clientID is", clientID)
		log.Println("Handler/GetClientActivity: dtoList is", dtoList)

		// auth true, but no records
		if dtoList == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write(nil)
			return
		}

		// auth true, records exist

		encVal := []struct { // encoded value
			Hash string `json:"short_url"`
			URL  string `json:"original_url"`
		}{}

		for _, dto := range dtoList {
			encVal = append(encVal, struct {
				Hash string "json:\"short_url\""
				URL  string "json:\"original_url\""
			}{Hash: dto.Hash, URL: dto.URL})
		}

		log.Println("Handler/GetClientActivity: ev is", encVal)

		// сначала устанавливаем заголовок Content-Type
		// для передачи клиенту информации, кодированной в JSON
		w.Header().Set("content-type", "application/json")

		// устанавливаем статус-код 200
		w.WriteHeader(http.StatusOK)

		// пишем тело ответа
		json.NewEncoder(w).Encode(encVal)

	}
}

func (uh URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	if uh.urlService.PingDB() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
