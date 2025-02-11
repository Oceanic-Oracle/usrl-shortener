package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"url-shortener/internal/storage/url"
)

type repository struct {
	logger  *slog.Logger
	storage map[string]string
}

func (r *repository) SaveUrl(ctx context.Context, url, shorturl string) error {
	r.storage[shorturl] = url
    r.logger.Debug("URL saved in memory", slog.String("shortURL", shorturl), slog.String("originalURL", url))
	return nil
}

func (r *repository) GetUrl(ctx context.Context, shorturl string) (string, error) {
	val, ok := r.storage[shorturl]
	if !ok {
		return "", fmt.Errorf("url not found: %s", shorturl)
	}
	return val, nil
}

func New(logger *slog.Logger) url.UrlRepository {
	temp := &repository{
		logger:  logger,
		storage: make(map[string]string),
	}
	return temp
}
