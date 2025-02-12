package postgre

import (
	"context"
	"fmt"
	"log/slog"
	"url-shortener/internal/storage/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (r *repository) SaveUrl(ctx context.Context, url string, shortUrl string) error {
	q := `
	INSERT INTO urls (url, short_url) 
    VALUES ($1, $2) 
	`

	if _, err := r.pool.Exec(ctx, q, url, shortUrl); err != nil {
		return fmt.Errorf("cannot save url %s and short url %s. err: %w", url, shortUrl, err)
	}

	return nil
}

func (r *repository) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	q := `
	SELECT url
	FROM urls
	WHERE short_url = $1
	`

	var url string
	if err := r.pool.QueryRow(ctx, q, shortUrl).Scan(&url); err != nil {
		return "", fmt.Errorf("cannot get url. err: %w", err)
	}

	return url, nil
}

func New(pool *pgxpool.Pool, logger *slog.Logger) url.UrlRepository {
	temp := &repository{
		pool:   pool,
		logger: logger,
	}
	return temp
}
