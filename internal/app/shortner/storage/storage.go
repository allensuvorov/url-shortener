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

func (us URLStorage) GetByURL(u string) (entity.URLEntity, error) {
	for _, v := range us.inMemory {
		if v.URL == u {
			return v, nil
		}
	}
	return nil, errors.NotFound
}

func (us URLStorage) GetByHash(u string) (entity.URLEntity, error) {
	for _, v := range us.inMemory {
		if v.Hash == u {
			return v, nil
		}
	}
	return nil, errors.NotFound
}
