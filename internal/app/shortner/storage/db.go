package storage

import (
	"context"
	"database/sql"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"time"
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
	db.Exec("CREATE TABLE IF NOT EXISTS url(ID INT PRIMARY KEY, URL TEXT, hash TEXT, client TEXT);")

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

	//defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.DB.PingContext(ctx); err != nil {
		return false
	}
	return true
}
