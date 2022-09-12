package url

import (
	"net/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// URLStorage interface (to be adopted for mock-testing):
type URLStorage interface {

	// Create new entity (pair shortURL: longURL).
	Create(url entity.URLEntity) error

	// GetByHash returns entity for the matching hash, checks if hash exists.
	GetByHash(u string) (entity.URLEntity, error)

	// GetByURL returns hash for the matching URL, checks if URL exists.
	GetHashByURL(u string) (h string, error)
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
func (us URLService) Create(u string) (string, error) {

	// check if URL is valid
	_, err := url.ParseRequestURI(u)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return "", err
	}

	// get Hash from DB if the longURL already exists in storage
	h, err := us.urlStorage.GetHashByURL(u)

	// Generate new Hash if URL does not exist in storage
	if err == errors.NotFound {

		// generate shortened URL
		h = Shorten(u)
		
		// New url entity
		ue := entity.URLEntity{ 
			URL: u, 
			Hash: h,
		} 

		// add url to the storage
		us.urlStorage.Create(ue)
	}

	shortURL := "http://localhost:8080/" + h
	return shortURL, nil
}

// TODO Get
func (us URLService) Get(u string) (string, error) {

}
