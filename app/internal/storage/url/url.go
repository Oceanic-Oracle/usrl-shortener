package url

import "context"

type UrlRepository interface {
	SaveUrl(context.Context, string, string) error
	GetUrl(context.Context, string) (string, error)
}