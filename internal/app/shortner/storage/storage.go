package storage

import (
	"log"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

// Object with storage methods to work with DB
type URLStorage struct {
	InMemory InMemory
}

type InMemory struct {
	UEs []entity.URLEntity
	CEs []entity.ClientEntity
}

// NewURLStorage creates URLStorage object
func NewURLStorage() *URLStorage {
	// Restore data at start up
	fsp := config.UC.FSP
	ues := []entity.URLEntity{}
	ces := []entity.ClientEntity{}
	inm := InMemory{
		UEs: ues,
		CEs: ces,
	}

	// restore if path in config not empty
	if fsp != "" {
		inm = restore(fsp) // get map
	}

	return &URLStorage{
		InMemory: inm,
	}
}

// Create adds new URL record to storage
func (us *URLStorage) Create(h, u string) error {
	// Save to map
	us.InMemory[h] = u
	log.Println("Storage/Create(): added to map, updated map len is", len(us.InMemory))
	log.Println("Storage/Create(): added to map, updated map is", us.InMemory)

	// get file storage path from config
	fsp := config.UC.FSP

	// Save to file, if there is path in config
	if fsp != "" {
		write(h, u, fsp)
	}
	log.Printf("Storage/Create(): created hash: %s, for URL: %s. File path %s:", h, u, fsp)
	return nil
}

func (us *URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage/GetHashByURL(), looking for matching URL", u)
	for k, v := range us.InMemory {
		if v == u {
			log.Println("Storage GetHashByURL, found record", k)
			return k, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *URLStorage) GetURLByHash(h string) (string, error) {
	log.Println("Storage/GetURLByHash(), looking in map len", len(us.InMemory))
	log.Println("Storage/GetURLByHash(), looking for matching Hash", h)
	u, ok := us.InMemory[h]
	if !ok {
		return "", errors.ErrNotFound
	}
	return u, nil
}
