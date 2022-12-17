package handlers

import (
	"encoding/json"
	errors2 "errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

type URLService interface {
	Create(ue entity.URLEntity) (string, error)

	Get(h string) (string, error)

	GetClientActivity(id string) ([]entity.URLEntity, error)

	PingDB() bool

	BatchDelete(hashList *[]string, clientID string) error
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
	var decVal struct { // decoded value
		URL string
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, errors.ErrWrongContentType.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&decVal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Handler/CreateForJSONClient - request: object", decVal)
	log.Println("Handler/CreateForJSONClient - URL in the request is", decVal.URL)

	ue := entity.URLEntity{
		URL: decVal.URL,
	}

	ue.ClientID = r.Header.Get("id")

	shortURL, err := uh.urlService.Create(ue)

	if err != nil && !errors2.Is(err, errors.ErrRecordExists) {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	encVal := struct { // encoded value
		Result string `json:"result"`
	}{
		Result: shortURL,
	}

	log.Println("Handler/CreateForJSONClient: ev is", encVal.Result)

	w.Header().Set("content-type", "application/json")

	if errors2.Is(err, errors.ErrRecordExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(encVal)
}

func (uh URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler/Create - Start")
	log.Println("Handler/Create - Set-Cookie:", w.Header().Get("Set-Cookie"))
	log.Println("Handler/Create - ID header:", r.Header.Get("id"))

	ue := entity.URLEntity{
		ClientID: r.Header.Get("id"),
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	log.Println("Handler/Create - read r.Body, no err")

	ue.URL = string(b)
	log.Println("Handler/Create - URL in the POST request is", ue.URL)

	ue.Hash, err = uh.urlService.Create(ue)

	if err != nil && !errors2.Is(err, errors.ErrRecordExists) {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	if errors2.Is(err, errors.ErrRecordExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte(ue.Hash))
}

func (uh URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	h := chi.URLParam(r, "hash")

	log.Println("Handler Get", h, r.URL.Path)

	u, err := uh.urlService.Get(h)

	if err == errors.ErrNotFound {
		http.Error(w, "URL does not exist", http.StatusBadRequest)
		return
	}

	if err == errors.ErrRecordDeleted {
		http.Error(w, errors.ErrRecordDeleted.Error(), http.StatusGone)
		return
	}

	w.Header().Set("Location", u)

	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(u))
}

func (uh URLHandler) GetClientActivity(w http.ResponseWriter, r *http.Request) {
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

		log.Println("Handler/GetClientUrls: clientID is", clientID)
		log.Println("Handler/GetClientUrls: dtoList is", dtoList)

		if dtoList == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write(nil)
			return
		}

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

		log.Println("Handler/GetClientUrls: ev is", encVal)

		w.Header().Set("content-type", "application/json")

		w.WriteHeader(http.StatusOK)

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

func (uh URLHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler/BatchCreate - Hello")
	log.Println("Handler/BatchCreate body:", r.Body)
	var decVals []struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, errors.ErrWrongContentType.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&decVals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Handler/BatchCreate - request: object", decVals)

	ue := entity.URLEntity{}

	clientID := r.Header.Get("id")

	var encVals []struct {
		ID   string `json:"correlation_id"`
		Hash string `json:"short_url"`
	}

	for _, v := range decVals {
		ue.URL = v.URL
		ue.ClientID = clientID
		shortURL, err := uh.urlService.Create(ue)
		if err != nil {
			http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		}

		encVals = append(encVals, struct {
			ID   string `json:"correlation_id"`
			Hash string `json:"short_url"`
		}{ID: v.ID, Hash: shortURL})
	}

	log.Println("Handler/BatchCreate: encVal is", encVals)

	w.Header().Set("content-type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(encVals)
}

func (uh URLHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("Handlers/BatchDelete - Hello")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, errors.ErrWrongContentType.Error(), http.StatusBadRequest)
		return
	}

	var decVals []string

	if err := json.NewDecoder(r.Body).Decode(&decVals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Handlers/BatchDelete - decoded request: object", decVals)

	err := uh.urlService.BatchDelete(&decVals, r.Header.Get("id"))
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusAccepted)

	log.Println("Handlers/BatchDelete - Buy")
}
