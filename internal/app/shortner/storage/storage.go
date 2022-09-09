package storage

import (
	"errors"
	"log"
)

// map to store short urls and full urls
var Urls map[string]string = make(map[string]string)
var ErrNotFound = errors.New("Resource was not found")

// DBStorage interface (to be adopted for mock-testing):
type DBStorage interface {
	// check if longURL exists
	URLExists(u string) (bool, string)

	// check if hash exists
	HashExists(u string) bool

	// return longURL for the matching hash
	GetURL(u string) string

	// return hash for a matching longURL
	GetHash(u string) (string, error)

	// add new record - pair shortURL: longURL
	CreateURL(h, string, u string) error
}

// check if longURL exists
func GetHash(u string) (string, error) {
	for k, v := range Urls {
		if v == u {
			log.Println("URL already exists")
			return k, nil
		}
	}
	return "", ErrNotFound
}

// check if hash exists
func HashExists(h string) bool {
	if _, ok := Urls[h]; ok {
		return true
	}
	return false
}

// return longURL for the matching hash
func GetURL(h string) string {
	return Urls[h]
}

// add new record - pair shortURL: longURL
func CreateURL(h string, u string) error {
	Urls[h] = u
	return nil
}
