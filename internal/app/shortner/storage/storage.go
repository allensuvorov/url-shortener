package storage

// map to store short urls and full urls
var Urls map[string]string = make(map[string]string)

// methods we need here:
// add new record - pair shortURL: longURL
// return longURL for the matching shortURL
// check if shortURL exists
// check if longURL exists
