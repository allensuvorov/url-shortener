package url

import (
	"log"
	"net/url"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// URLStorage interface (to be adopted for mock-testing):
type URLStorage interface {

	// Create new entity (pair shortURL: longURL).
	Create(h, u string) error

	// GetByHash returns entity for the matching hash, checks if hash exists.
	GetURLByHash(u string) (string, error)

	// GetByURL returns hash for the matching URL, checks if URL exists.
	GetHashByURL(u string) (string, error)
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
		return "", err
	}

	// get Hash from DB if the longURL already exists in storage
	h, err := us.urlStorage.GetHashByURL(u)

	// Generate new Hash if URL does not exist in storage
	if err == errors.ErrNotFound {

		// generate Hash for the shortened URL
		h = buildHash(u)

		// cut it to a short hash
		h = getUniqShortHash(h, u, us)

		log.Println("created new shortURL", h)

		// add url to the storage
		err = us.urlStorage.Create(h, u)
		if err != nil {
			return "", err
		}
	}
	// // Set local env var
	// os.Setenv("BASE_URL", "http://localhost:8080/")
	// // Get local env var
	// bu := os.Getenv("BASE_URL")

	bu, ok := os.LookupEnv("BASE_URL")
	if !ok {
		log.Printf("%s not set\n; passing default", "BASE_URL")
		bu = "http://localhost:8080"
	}

	log.Println("Service: BASE_URL from local env is:", bu)
	shortURL := bu + "/" + h
	return shortURL, nil
}

// TODO Get
func (us URLService) Get(h string) (string, error) {
	// check if hash exists, if not - return 400
	u, err := us.urlStorage.GetURLByHash(h)

	if err == errors.ErrNotFound {
		return "", err
	}
	return u, nil
}
