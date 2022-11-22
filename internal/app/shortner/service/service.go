package service

import (
	"log"
	"net/url"

	"github.com/allensuvorov/urlshortner/internal/app/shortner/config"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"
	"github.com/allensuvorov/urlshortner/internal/app/shortner/domain/errors"
)

type URLStorage interface {
	Create(ue entity.URLEntity) error

	GetURLByHash(u string) (string, error)

	GetHashByURL(u string) (string, error)

	GetClientUrls(id string) ([]entity.URLEntity, error)

	PingDB() bool

	BatchDelete(hashList []string, clientID string) error
}

type URLService struct {
	urlStorage URLStorage
}

func NewURLService(us URLStorage) URLService {
	return URLService{
		urlStorage: us,
	}
}

func (us URLService) Create(ue entity.URLEntity) (string, error) {

	_, err := url.ParseRequestURI(ue.URL)
	if err != nil {
		return "", err
	}

	h, err := us.urlStorage.GetHashByURL(ue.URL)
	if err == nil {
		err = errors.ErrRecordExists
	}
	log.Println("Service/Create(): about go get BU from config")
	bu := config.UC.BU
	log.Println("Service: BASE_URL from local env is:", bu)

	if err == errors.ErrNotFound {

		h = generateHash(ue.URL)

		h = getUniqShortHash(h, ue.URL, us)

		log.Println("Service/Create(): created new shortURL", h)

		ue.Hash = h
		err = us.urlStorage.Create(ue)
		if err != nil {
			return "", err
		}
		log.Println("Service/Create(): saved new shortURL in map", h)
		err = nil
	}
	shortURL := bu + "/" + h
	return shortURL, err
}

func (us URLService) Get(h string) (string, error) {
	u, err := us.urlStorage.GetURLByHash(h)

	//if err == errors.ErrNotFound {
	//	return "", err
	//}
	if err != nil {
		return "", err
	}
	return u, nil
}

func (us URLService) GetClientActivity(id string) ([]entity.URLEntity, error) {
	dtoList, err := us.urlStorage.GetClientUrls(id)
	log.Println("service/GetClientUrls client ID is:", id)
	log.Println("service/GetClientUrls dtoList is:", dtoList)
	if err != nil {
		return nil, err
	}
	return dtoList, nil
}

func (us URLService) PingDB() bool {
	return us.urlStorage.PingDB()
}

func (us URLService) BatchDelete(hashList []string, clientID string) error {
	log.Println("service/BatchDelete - Hello")
	log.Println("service/BatchDelete - Bye")
	return us.urlStorage.BatchDelete(hashList, clientID)
}
