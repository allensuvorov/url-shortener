package url

import (
	"log"
	"net/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// URLStorage interface (to be adopted for mock-testing):
type URLStorage interface {

	// Create new entity (pair shortURL: longURL).
	Create(ue entity.DTO) error

	// GetByHash returns entity for the matching hash, checks if hash exists.
	GetURLByHash(u string) (string, error)

	// GetByURL returns hash for the matching URL, checks if URL exists.
	GetHashByURL(u string) (string, error)

	GetClientActivity(id string) ([]entity.DTO, error)
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
func (us URLService) Create(ue entity.DTO) (string, error) {

	// check if URL is valid
	_, err := url.ParseRequestURI(ue.URL)
	if err != nil {
		return "", err
	}

	// get Hash from DB if the longURL already exists in storage
	h, err := us.urlStorage.GetHashByURL(ue.URL)

	// Generate new Hash if URL does not exist in storage
	if err == errors.ErrNotFound {

		// generate Hash for the shortened URL
		h = buildHash(ue.URL)

		// cut it to a short hash
		h = getUniqShortHash(h, ue.URL, us)

		log.Println("Service/Create(): created new shortURL", h)

		// add url to the storage
		ue.Hash = h
		err = us.urlStorage.Create(ue)
		if err != nil {
			return "", err
		}
		log.Println("Service/Create(): saved new shortURL in map", h)
	}

	// Get Base URL
	log.Println("Service/Create(): about go get BU from config")
	bu := config.UC.BU
	log.Println("Service: BASE_URL from local env is:", bu)
	shortURL := bu + "/" + h

	return shortURL, nil
}

func (us URLService) Get(h string) (string, error) {
	// check if hash exists, if not - return 400
	u, err := us.urlStorage.GetURLByHash(h)

	if err == errors.ErrNotFound {
		return "", err
	}
	return u, nil
}

func (us URLService) GetClientActivity(id string) ([]entity.DTO, error) {
	dtoList, err := us.urlStorage.GetClientActivity(id)
	log.Println("service/GetClientActivity client ID is:", id)
	log.Println("service/GetClientActivity dtoList is:", dtoList)
	if err != nil {
		return nil, err
	}
	return dtoList, nil
}
