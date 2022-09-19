package storage

import (
	"log"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

// Object with storage methods to work with DB
type URLStorage struct {
	InMemory *hashmap.URLHashMap
}

// NewURLStorage creates URLStorage object
func NewURLStorage() *URLStorage {
	return &URLStorage{
		InMemory: new(hashmap.URLHashMap),
	}
}

// Create adds new URL record to storage
func (us *URLStorage) Create(ue *hashmap.URLHashMap) error {
	us.InMemory[]
	log.Println("Storage Create UE, appended, updated slice len is", len(us.InMemory))
	log.Println("Storage Create UE", ue)
	return nil
}

func (us *URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage GetHashByURL, looking for matching URL", u)
	for _, v := range us.InMemory {
		if v.URL == u {
			log.Println("Storage GetHashByURL, found ue", v)
			return v.Hash, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *URLStorage) GetURLByHash(u string) (string, error) {
	log.Println("Storage GetURLByHash, looking in slice len", len(us.InMemory))
	log.Println("Storage GetURLByHash, looking for matching Hash", u)
	for _, v := range us.InMemory {
		log.Println("Storage GetURLByHash, comparing", v.Hash, u)
		if v.Hash == u {
			log.Println("Storage GetURLByHash, found ue", v)
			return v.URL, nil
		}
	}
	return "", errors.ErrNotFound
}
