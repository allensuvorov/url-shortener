package storage

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
	NewURL(u string) string
}

func URLExists(u string) bool, string {
	for k, v := range Urls {
		if v == u {
			return true
		}
	}
	return false
}
