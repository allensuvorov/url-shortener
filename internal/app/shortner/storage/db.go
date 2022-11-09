package storage

import (
	"database/sql"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
)

type URLEntity struct {
	URL      string
	Hash     string
	ClientID string
}

type urlDB struct {
	DB *sql.DB
}

func NewUrlDB() *urlDB {
	db, err := sql.Open("pgx",
		config.UC.DSN)
	if err != nil {
		panic(err)
	}
	return &urlDB{
		DB: db,
	}
}

func (db urlDB) Create(ue entity.DTO) error {
	return nil
}

func (db urlDB) GetURLByHash(u string) (string, error) {
	return "", nil
}

func (db urlDB) GetHashByURL(u string) (string, error) {
	return "", nil
}

func (db urlDB) GetClientActivity(id string) ([]entity.DTO, error) {
	return nil, nil
}

func (db urlDB) PingDB() bool {
	return true
}
