package url

import "github.com/allensuvorov/urlshortner/internal/app/shortner/domain/entity"

type URLService interface {
	// CreateURL
	CreateURL(url entity.URLEntity) error
	// GetByURL
	GetByURL(u string) (entity.URLEntity, error)
	// GetByHash
	GetByHash(h string) (entity.URLEntity, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(us URLService) URLHandler {
	return URLHandler{
		urlService: us,
	}
}
