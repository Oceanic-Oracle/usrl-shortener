package postgre

import (
	"log/slog"
	"url-shortener/internal/storage/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (r *repository) SaveUrl(string) error {
	return nil
}

func (r *repository) GetUrl(string) (string, error) {
	return "", nil
}

func New(pool *pgxpool.Pool, logger *slog.Logger) url.UrlRepository {
	temp := &repository{
		pool:   pool,
		logger: logger,
	}
	return temp
}
