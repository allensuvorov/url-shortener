package storage

import "log"

// map to store short urls and full urls
var Urls map[string]string = make(map[string]string)

// DBStorage interface has methods we need here:
type DBStorage interface {
	// check if longURL exists
	URLExists(u string) (bool, string)

	// check if hash exists
	HashExists(u string) bool

	// return longURL for the matching hash
	GetURL(u string) string

	// add new record - pair shortURL: longURL
	NewURL(h, string, u string) error
}

func URLExists(u string) (bool, string) {
	for k, v := range Urls {
		if v == u {
			log.Println("URL already exists")
			return true, k
		}
	}
	return false, ""
}

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
func NewURL(h string, u string) error {
	Urls[h] = u
	return nil
}
