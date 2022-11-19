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
	Create(ue entity.URLEntity) error

	// GetURLByHash returns entity for the matching hash, checks if hash exists.
	GetURLByHash(u string) (string, error)

	// GetHashByURL returns hash for the matching URL, checks if URL exists.
	GetHashByURL(u string) (string, error)

	GetClientUrls(id string) ([]entity.URLEntity, error)

	PingDB() bool
}

type URLService struct {
	urlStorage URLStorage
}

func NewURLService(us URLStorage) URLService {
	return URLService{
		urlStorage: us,
	}
}

func (us URLService) Create(ue entity.URLEntity) (string, error) {

	// check if URL is valid
	_, err := url.ParseRequestURI(ue.URL)
	if err != nil {
		return "", err
	}

	// get Hash from DB if the longURL already exists in storage
	h, err := us.urlStorage.GetHashByURL(ue.URL)
	if err == nil {
		err = errors.ErrRecordExists
	}
	// Get Base URL
	log.Println("Service/Create(): about go get BU from config")
	bu := config.UC.BU
	log.Println("Service: BASE_URL from local env is:", bu)

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
		err = nil
	}
	shortURL := bu + "/" + h
	return shortURL, err
}

func (us URLService) Get(h string) (string, error) {
	// check if hash exists, if not - return 400
	u, err := us.urlStorage.GetURLByHash(h)

	if err == errors.ErrNotFound {
		return "", err
	}
	return u, nil
}

func (us URLService) GetClientActivity(id string) ([]entity.URLEntity, error) {
	dtoList, err := us.urlStorage.GetClientUrls(id)
	log.Println("service/GetClientUrls client ID is:", id)
	log.Println("service/GetClientUrls dtoList is:", dtoList)
	if err != nil {
		return nil, err
	}
	return dtoList, nil
}

func (us URLService) PingDB() bool {
	return us.urlStorage.PingDB()
}
