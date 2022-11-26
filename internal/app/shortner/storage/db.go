package storage

import (
	"context"
	"database/sql"
	errors2 "errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

type urlDB struct {
	DB *sql.DB
	//buffer   []string
	BufferCh chan urlDeleteRequest
	//clientID string
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
		DB:       db,
		BufferCh: make(chan urlDeleteRequest, 1000),
		//buffer: make([]string, 0, 1000),
	}
}

func (db *urlDB) Create(ue entity.URLEntity) error {
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

func (db *urlDB) GetURLByHash(u string) (string, error) {
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

func (db *urlDB) GetHashByURL(u string) (string, error) {
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

func (db *urlDB) GetClientUrls(id string) ([]entity.URLEntity, error) {
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

func (db *urlDB) PingDB() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.DB.PingContext(ctx); err != nil {
		return false
	}
	return true
}

type urlDeleteRequest struct {
	hash     string
	clientID string
}

func (db *urlDB) BatchDelete(hashLIst *[]string, clientID string) error {
	log.Println("urlDB/BatchDelete - Hello")

	go func(hashList *[]string, clientID string) {
		wg := &sync.WaitGroup{}

		for _, h := range *hashList {
			wg.Add(1)
			udr := urlDeleteRequest{h, clientID}

			go func(v urlDeleteRequest) {
				defer wg.Done()
				select {
				case db.BufferCh <- v:
				default:
					db.flushBufferToDB()
				}

			}(udr)
		}
		wg.Wait()
	}(hashLIst, clientID)
	go func() {
		db.flushBufferToDB()
	}()

	log.Println("urlDB/BatchDelete - Bye")

	return nil
}

func (db *urlDB) flushBufferToDB() error {
	log.Println("urlDB/flushBufferToDB - Hello")
	startTimer := time.Now()

	if db.DB == nil {
		return errors2.New("you haven`t opened the database connection")
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`UPDATE urls SET deleted = TRUE WHERE hash = $1 AND client = $2 AND deleted = FALSE;`,
	)
	if err != nil {
		return err
	}

	defer stmt.Close()

loop:
	for {
		select {
		case v := <-db.BufferCh:
			if _, err = stmt.Exec(v.hash, v.clientID); err != nil {
				return err
			}
		default:
			break loop
		}
	}
	duration := time.Since(startTimer)
	fmt.Printf("urlDB/flushBufferToDB - Execution Time ms %d\n", duration.Milliseconds())
	log.Println("urlDB/flushBufferToDB - Bye")

	return tx.Commit()
}
