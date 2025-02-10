package storage

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/url"
	"url-shortener/internal/storage/url/inmemory"
	"url-shortener/internal/storage/url/postgre"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postgreType  string = "postgre"
	inMemoryType string = "inmemory"
)

type Repositories struct {
	url url.UrlRepository

	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, cfg *config.StorageCfg, logger *slog.Logger) (client *Repositories) {
	switch cfg.StorageType {
	case postgreType:
		pool, err := postgreHandler(ctx, cfg, "disable")
		if err != nil {
			log.Fatalf("failed to connect to PostgreSQL: %v", err)
		} else {
			logger.Info("Successfull connect to PostgreSQL")
		}
		client = &Repositories{
			url:  postgre.New(pool, logger),
			pool: pool,
		}
	case inMemoryType:
		client = &Repositories{
			url:  inmemory.New(logger),
			pool: nil,
		}
	default:
		log.Fatalf("unsupported storage type: %q. Available options: %q, %q", 
            cfg.StorageType, 
            postgreType, 
            inMemoryType,
        )
	}
	return client
}

func (r *Repositories) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func postgreHandler(ctx context.Context, cfg *config.StorageCfg, sslmode string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, fmt.Sprintf(
		"postgresql://%s:%s@postgres:5432/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.DataBase, sslmode,
	))
}
