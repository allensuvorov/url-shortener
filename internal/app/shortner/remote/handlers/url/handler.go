package url

import "github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"

type UrlService interface {
	// CreateUrl
	CreateUrl(url entity.URLEntity) error
	// GetByUrl
	GetByURL(u string) (entity.URLEntity, error)
	// GetByHash
	GetByHash(h string) (entity.URLEntity, error)
}
