package storage

import (
	"context"
	"database/sql"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
	"log"
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
	db.Exec(`CREATE TABLE IF NOT EXISTS urls(
    ID SERIAL PRIMARY KEY, 
    URL TEXT, 
    hash TEXT, 
    client TEXT
                              );`)

	log.Println("created new URL Database")
	return &urlDB{
		DB: db,
	}
}

func (db urlDB) Create(ue entity.DTO) error {
	db.DB.Exec(
		`INSERT INTO urls
		(url, hash, client)
		VALUES
		($T, $T, $T);`,
		ue.URL, ue.Hash, ue.ClientID,
	)
	return nil
}

func (db urlDB) GetURLByHash(u string) (string, error) {
	row, err := db.DB.Query(`SELECT url FROM urls WHERE hash = $T;`, u)

	if err != nil {
		log.Println("urlBD/GetURLByHash, record not found")
		return "", errors.ErrNotFound
	}

	var url string
	err = row.Scan(&url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Storage GetHashByURL, found record", url)
	return url, nil
}

func (db urlDB) GetHashByURL(u string) (string, error) {
	row, err := db.DB.Query(`SELECT hash FROM urls WHERE url = $T;`, u)

	if err != nil {
		log.Println("urlBD/GetHashByURL, record not found")
		return "", errors.ErrNotFound
	}

	var hash string
	err = row.Scan(&hash)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Storage GetHashByURL, found record", hash)
	return hash, nil
}

func (db urlDB) GetClientUrls(id string) ([]entity.DTO, error) {
	log.Println("urlDB/GetClientUrls client id is:", id)
	urlEntities := make([]entity.DTO, 0)

	rows, err := db.DB.Query(`SELECT url, hash, client FROM urls WHERE client = id;`)

	defer rows.Close()

	for rows.Next() {
		var urlEntity entity.DTO

		err = rows.Scan(&urlEntity.URL, &urlEntity.Hash, &urlEntity.ClientID)
		if err != nil {
			return nil, err
		}

		urlEntities = append(urlEntities, urlEntity)
	}

	log.Println("storage/GetClientUrls client urlEntities is:", urlEntities)

	return urlEntities, nil
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
