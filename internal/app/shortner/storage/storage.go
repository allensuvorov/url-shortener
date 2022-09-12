package storage

import (
	"log"
	"yandex/projects/urlshortner/internal/app/shortner/domain/entity"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	//"github.com/allensuvorov/urlshortner/internal/app/domain/entity"
)

// map to store short urls and full urls
// var Urls map[string]string = make(map[string]string)

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
	return nil
}

func (us URLStorage) GetByURL(u string) (entity.URLEntity, error) {
	for _, v := range us.inMemory {
		if v.URL == u {
			return v, nil
		}
	}
	return nil, errors.NotFound
}

// Get URL record by URL or Hash
func (us URLStorage) Get(ur entity.URLEntity) (entity.URLEntity, error) {
	if len(ur.URL) != 0 {
		for _, v := range us.inMemory {
			if v.URL == ur.URL {
				return v
			}
		}
	}
	if len(ur.Hash) != 0 {
		for _, v := range us.inMemory {
			if v.Hash == ur.Hash {
				return v
			}
		}
	}
	return ur, nil
}

// HashExists checks if hash exists.
func HashExists(h string) bool {
	if _, ok := Urls[h]; ok {
		return true
	}
	return false
}

// GetHash returns hash for the original longURL.
func GetHash(u string) (string, bool) {
	for k, v := range Urls {
		if v == u {
			log.Println("URL already exists")
			return k, true
		}
	}
	return "", false
}

// GetURL returns longURL for the matching hash.
func GetURL(h string) (string, bool) {
	if u, ok := Urls[h]; ok {
		return u, true
	}
	return "", false
}
