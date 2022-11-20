package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

type urlDB struct {
	DB *sql.DB
}

func NewURLDB() *urlDB {
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

func (db urlDB) Create(ue entity.URLEntity) error {
	_, err := db.DB.Exec(
		`INSERT INTO urls
		(url, hash, client)
		VALUES
		($1, $2, $3);`,
		ue.URL, ue.Hash, ue.ClientID,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (db urlDB) GetURLByHash(u string) (string, error) {
	row := db.DB.QueryRow(`SELECT url FROM urls WHERE hash = $1;`, u)

	var url string
	err := row.Scan(&url)

	if err == sql.ErrNoRows {
		log.Println("urlBD/GetHashByURL, record not found")
		return "", errors.ErrNotFound
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Storage GetHashByURL, found record", url)
	return url, nil
}

func (db urlDB) GetHashByURL(u string) (string, error) {
	row := db.DB.QueryRow(`SELECT hash FROM urls WHERE url = $1;`, u)

	var hash string
	err := row.Scan(&hash)
	if err == sql.ErrNoRows {
		log.Println("urlBD/GetHashByURL, record not found")
		return "", errors.ErrNotFound
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("urlBD/GetHashByURL, found record", hash)
	return hash, nil
}

func (db urlDB) GetClientUrls(id string) ([]entity.URLEntity, error) {
	log.Println("urlDB/GetClientUrls client id is:", id)
	urlEntities := make([]entity.URLEntity, 0)

	rows, err := db.DB.Query(`SELECT url, hash, client FROM urls WHERE client = $1;`, id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var urlEntity entity.URLEntity

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
