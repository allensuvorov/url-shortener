package storage

// map to store short urls and full urls
var Urls map[string]string = make(map[string]string)

// DBStorage interface has methods we need here:
type DBStorage interface {
	// check if longURL exists
	URLExists(URL string) bool

	// return longURL for the matching shortURL
	GetURL(ShortURL string) string

	// add new record - pair shortURL: longURL
	NewURL(URL string) string

	// check if shortURL exists
	ShortURLExists(ShortURL string) bool
}
