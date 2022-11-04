package storage

import (
	"log"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/hashmap"
)

// Object with storage methods to work with DB
type URLStorage struct {
	InMemory InMemory
}

type InMemory struct {
	URLHashMap     hashmap.URLHashMap
	ClientActivity hashmap.ClientActivity
}

// NewURLStorage creates URLStorage object
func NewURLStorage() *URLStorage {

	// Restore data at start up
	fsp := config.UC.FSP
	um := make(hashmap.URLHashMap) // url map
	ca := make(hashmap.ClientActivity)
	im := InMemory{um, ca}

	// TODO restore from file both um and ua

	// restore if path in config not empty
	if fsp != "" {
		im = restore(fsp) // get map
	}

	return &URLStorage{
		InMemory: im,
	}
}

// Create adds new URL record to storage
func (us *URLStorage) Create(ue entity.DTO) error {
	// Save to map
	log.Println("Storage/Create(): hello")

	us.InMemory.URLHashMap[ue.Hash] = ue.URL
	log.Println("Storage/Create(): added to map, updated map len is", len(us.InMemory.URLHashMap))
	log.Println("Storage/Create(): added to map, updated map is", us.InMemory.URLHashMap)

	_, ok := us.InMemory.ClientActivity[ue.ClientID]
	if !ok {
		us.InMemory.ClientActivity[ue.ClientID] = make(map[string]bool)
	}
	us.InMemory.ClientActivity[ue.ClientID][ue.Hash] = true

	// get file storage path from config
	fsp := config.UC.FSP

	// Save to file, if there is path in config
	if fsp != "" {
		write(ue, fsp)
	}
	log.Printf("Storage/Create(): created hash: %s, for URL: %s. File path %s:", ue.Hash, ue.URL, fsp)
	return nil
}

func (us *URLStorage) GetHashByURL(u string) (string, error) {
	log.Println("Storage/GetHashByURL(), looking for matching URL", u)
	for k, v := range us.InMemory.URLHashMap {
		if v == u {
			log.Println("Storage GetHashByURL, found record", k)
			return k, nil
		}
	}
	return "", errors.ErrNotFound
}

func (us *URLStorage) GetURLByHash(h string) (string, error) {
	log.Println("Storage/GetURLByHash(), looking in map len", len(us.InMemory.URLHashMap))
	log.Println("Storage/GetURLByHash(), looking for matching Hash", h)
	u, ok := us.InMemory.URLHashMap[h]
	if !ok {
		return "", errors.ErrNotFound
	}
	return u, nil
}

func (us *URLStorage) GetClientActivity(id string) ([]entity.DTO, error) {
	log.Println("storage/GetClientActivity client id is:", id)
	ca := us.InMemory.ClientActivity[id]
	log.Println("storage/GetClientActivity client ClientActivity is:", ca)
	dtoList := []entity.DTO{}

	for k := range ca {
		u, err := us.GetURLByHash(k)
		if err != nil {
			return nil, err
		}
		ue := entity.DTO{
			Hash: k,
			URL:  u,
		}
		dtoList = append(dtoList, ue)
	}
	log.Println("storage/GetClientActivity dtoList is:", dtoList)

	return dtoList, nil
}
