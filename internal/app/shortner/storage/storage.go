package storage

import (
	"log"
	"os"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

// Object with storage methods to work with DB
type URLStorage struct {
	InMemory hashmap.URLHashMap
}

// NewURLStorage creates URLStorage object
func NewURLStorage() *URLStorage {
	return &URLStorage{
		InMemory: make(hashmap.URLHashMap),
	}
}

// Create adds new URL record to storage
func (us *URLStorage) Create(h, u string) error {
	// Save to map
	us.InMemory[h] = u
	log.Println("Storage Create UE, added to map, updated map len is", len(us.InMemory))
	log.Println("Storage Create UE, added to map, updated map is", us.InMemory)
	// Save to file
	fsp, _ := os.LookupEnv("FILE_STORAGE_PATH")

	if len(fsp) > 0 { // if there is path in env var
		// TODO: save to file...
		log.Printf("Storage: saving to path - %s", fsp)
	}
	return nil
}

func (us *URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage GetHashByURL, looking for matching URL", u)
	for k, v := range us.InMemory {
		if v == u {
			log.Println("Storage GetHashByURL, found record", k)
			return k, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *URLStorage) GetURLByHash(h string) (string, error) {
	log.Println("Storage GetURLByHash, looking in map len", len(us.InMemory))
	log.Println("Storage GetURLByHash, looking for matching Hash", h)
	u, ok := us.InMemory[h]
	if !ok {
		return "", errors.ErrNotFound
	}
	return u, nil
}
