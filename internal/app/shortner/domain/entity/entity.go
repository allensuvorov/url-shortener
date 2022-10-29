package entity

type URLEntity struct {
	client int // foreign key
	hash   string
	url    string
}

type ClientEntity struct {
	id              int
	cookieID        string
	cookieSignature string
	key             string
}
