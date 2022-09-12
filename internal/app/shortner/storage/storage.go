package storage

import (
	"yandex/projects/urlshortner/internal/app/shortner/domain/entity"
	"yandex/projects/urlshortner/internal/app/shortner/domain/errors"

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
	return nil
}

func (us URLStorage) GetHashByURL(u string) (string, error) {
	for _, v := range us.inMemory {
		if v.URL == u {
			return v.Hash, nil
		}
	}
	return "", errors.NotFound
}

func (us URLStorage) GetURLByHash(u string) (string, error) {
	for _, v := range us.inMemory {
		if v.Hash == u {
			return v.URL, nil
		}
	}
	return "", errors.NotFound
}
