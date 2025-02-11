package api

import (
	"context"
	"fmt"
	"log/slog"
	"url-shortener/internal/api/gen"
	"url-shortener/internal/storage"
	"url-shortener/internal/utils"
)

type UrlService struct {
	gen.UnimplementedUrlShortenerServer
	logger      *slog.Logger
	repositories *storage.Repositories
}

func (r *UrlService) RegUrl(ctx context.Context, in *gen.RegUrlReq) (*gen.RegUrlResp, error) {
	r.logger.Debug("Recived RegUrl")

	var attempts int
    var shorturl string
    var err error
    for attempts < 5 {
        shorturl, err = utils.GenerateShortURL(in.Url, 10)
        if err == nil {
            break
        }
        attempts++
        r.logger.Error("cannot get short url", slog.Any("error", err))
    }
	if err != nil {
		return nil, fmt.Errorf("failed to generate short URL after %d attempts: %w", attempts + 1, err)
	}

	r.logger.Debug(fmt.Sprintf("generated url: %s", shorturl))

	if err := r.repositories.Url.SaveUrl(ctx, in.Url, shorturl); err != nil {
		return nil, fmt.Errorf("failed to save short URL: %w", err)
	}
	resp := &gen.RegUrlResp{
		ShortUrl: shorturl,
	}

	return resp, nil
}

func (r *UrlService) GetUrl(ctx context.Context, in *gen.GetUrlReq) (*gen.GetUrlResp, error) {
	r.logger.Debug("Recived GetUrl")

	url, err := r.repositories.Url.GetUrl(ctx, in.ShortUrl)
	if err != nil {
		r.logger.Error("cannot get url", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("url: %s", url))

	response := &gen.GetUrlResp{
		Url: url,
	}
	return response, nil
}

func NewClient(logger *slog.Logger, repositories *storage.Repositories) gen.UrlShortenerServer {
	return &UrlService{logger: logger, repositories: repositories}
}
