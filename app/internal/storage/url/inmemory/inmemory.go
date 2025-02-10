package inmemory

import (
	"log/slog"
	"url-shortener/internal/storage/url"
)

type repository struct {
	logger *slog.Logger
}

func (r *repository) SaveUrl(string) error {
	return nil
}

func (r *repository) GetUrl(string) (string, error) {
	return "", nil
}

func New(logger *slog.Logger) url.UrlRepository {
	temp := &repository{
		logger: logger,
	}
	return temp
}