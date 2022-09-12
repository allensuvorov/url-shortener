package url

import (
	"net/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
)

// URLStorage interface (to be adopted for mock-testing):
type URLStorage interface {

	// add new record - pair shortURL: longURL
	Create(url entity.URLEntity) error

	// return entity for the matching hash, check if hash exists
	GetByHash(u string) (entity.URLEntity, error)

	// return entity for the matching URL, check if URL exists
	GetByURL(u string) (entity.URLEntity, error)
}

type URLService struct {
	urlStorage URLStorage
}

func NewURLService(us URLStorage) URLService {
	return URLService{
		urlStorage: us,
	}
}

// func CreateURL
func (us URLService) CreateURL(ue entity.URLEntity) (string, error) {

	// check if URL is valid
	_, err := url.ParseRequestURI(ue.URL)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return "", err
	}

	// get Hash if the longURL already exists in storage
	h, ok := us.urlStorage.GetByURL(ue.URL)

	// if longURL does not exist in storage
	if !ok {

		// generate shortened URL
		h = Shorten(u)

		// add url to the storage
		storage.CreateHash(h, u)
	}

	shortURL := "http://localhost:8080/" + h
}

// TODO GetByURL
// TODO GetByHash
