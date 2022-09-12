package url

import (
	"io"
	"log"
	"net/http"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
)

type URLService interface {
	// CreateURL
	CreateURL(ue entity.URLEntity) (string, error)
	// GetByURL
	GetByURL(u string) (entity.URLEntity, error)
	// GetByHash
	GetByHash(h string) (entity.URLEntity, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(us URLService) URLHandler {
	return URLHandler{
		urlService: us,
	}
}

func (uh URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var ue entity.URLEntity // url entity

	// читаем Body
	b, err := io.ReadAll(r.Body)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// convert to string
	u := string(b)
	ue.URL = u

	// check if long URL is empty string
	if len(b) == 0 {
		http.Error(w, "empty URL", http.StatusBadRequest)
		return
	}

	// log body from request
	log.Println("URL in the POST request is", u)

	shortURL, err := uh.urlService.CreateURL(ue)

	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
	}

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	// пишем тело ответа
	w.Write([]byte(shortURL))
}

func (uh URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {

}
