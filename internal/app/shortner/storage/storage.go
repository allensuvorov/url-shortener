package storage

import "log"

// map to store short urls and full urls
var Urls map[string]string = make(map[string]string)

// DBStorage interface has methods we need here:
type DBStorage interface {
	// check if longURL exists
	URLExists(u string) (bool, string)

	// check if shortURL exists
	ShortURLExists(u string) bool

	// return longURL for the matching shortURL
	GetURL(u string) string

	// add new record - pair shortURL: longURL
	NewURL(u string) error
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

func ShortURLExists(u string) bool {
	if _, ok := Urls[u]; ok {
		return true
	}
	return false
}

// return longURL for the matching shortURL
func GetURL(u string) string {
	return Urls[u]
}

// add new record - pair shortURL: longURL
func NewURL(u string) error {
	return nil
}
