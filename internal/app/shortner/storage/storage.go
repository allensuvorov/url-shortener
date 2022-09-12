package storage

import (
	// "yandex/projects/urlshortner/internal/app/shortner/domain/entity"
	// "yandex/projects/urlshortner/internal/app/shortner/domain/errors"

	"log"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// Object with storage methods to work with DB
type URLStorage struct {
	inMemory []entity.URLEntity
}

// NewURLStorage creates URLStorage object
func NewURLStorage() URLStorage {
	return URLStorage{
		inMemory: make([]entity.URLEntity, 0),
	}
}

// Create adds new URL record to storage
func (us URLStorage) Create(url entity.URLEntity) error {
	us.inMemory = append(us.inMemory, url)
	log.Println("Storage Create UE, appended, updated slice len is", len(us.inMemory))
	log.Println("Storage Create UE", url)
	return nil
}

func (us URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage GetHashByURL, looking for matching URL", u)
	for _, v := range us.inMemory {
		if v.URL == u {
			log.Println("Storage GetHashByURL, found ue", v)
			return v.Hash, nil
		}
	}
	return "", errors.NotFound
}

func (us URLStorage) GetURLByHash(u string) (string, error) {
	log.Println("Storage GetURLByHash, looking in slice len", len(us.inMemory))
	log.Println("Storage GetURLByHash, looking for matching Hash", u)
	for _, v := range us.inMemory {
		log.Println("Storage GetURLByHash, comparing", v.Hash, u)
		if v.Hash == u {
			log.Println("Storage GetURLByHash, found ue", v)
			return v.URL, nil
		}
	}
	return "", errors.NotFound
}
