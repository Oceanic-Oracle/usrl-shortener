package url

type UrlRepository interface {
	SaveUrl(string) error
	GetUrl(string) (string, error)
}