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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls(
    ID SERIAL PRIMARY KEY, 
    URL TEXT, 
    hash TEXT, 
    client TEXT,
    deleted BOOL DEFAULT 'false'
                              );`)
	if err != nil {
		log.Fatal(err)
	}

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
	row := db.DB.QueryRow(`SELECT url, deleted FROM urls WHERE hash = $1;`, u)

	var url string
	var deleted bool
	err := row.Scan(&url, &deleted)

	if err == sql.ErrNoRows {
		log.Println("urlBD/GetHashByURL, record not found")
		return "", errors.ErrNotFound
	}

	if deleted {
		return "", errors.ErrRecordDeleted
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.DB.PingContext(ctx); err != nil {
		return false
	}
	return true
}

//func (db urlDB) BatchDelete(hashList []string, clientID string) error {
//	log.Println("urlDB/BatchDelete - Hello")
//
//	// TODO: review batch update,
//	// TODO: and fan-in
//
//	for _, h := range hashList {
//
//		_, err := db.DB.Exec(`UPDATE urls SET deleted = TRUE WHERE hash = $1 AND client = $2;`, h, clientID)
//		if err != nil {
//			return err
//		}
//	}
//
//	log.Println("urlDB/BatchDelete - Bye")
//	return nil
//}

func (db urlDB) BatchDelete(hashList []string, clientID string) error {
	// шаг 1 — объявляем транзакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	//шаг 2 — готовим инструкцию
	stmt, err := tx.Prepare(
		`UPDATE urls SET deleted = TRUE WHERE hash = $1 AND client = $2;`,
	)
	if err != nil {
		return err
	}

	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	for _, h := range hashList {
		// шаг 3 — указываем, что каждая короткая ссылка будет добавлена в транзакцию
		if _, err = stmt.Exec(h, clientID); err != nil {
			return err
		}
	}
	// шаг 4 — сохраняем изменения
	return tx.Commit()
}
