package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"url-shortener/internal/api/gen"
	"url-shortener/internal/storage"
	"url-shortener/internal/utils"

	"google.golang.org/grpc"
)

type UrlService struct {
	gen.UnimplementedUrlShortenerServer
	Logger      *slog.Logger
	Repositories *storage.Repositories
}

func (r *UrlService) RegUrl(ctx context.Context, in *gen.RegUrlReq) (*gen.RegUrlResp, error) {
	r.Logger.Debug("Recived RegUrl")

	var attempts int
    var shorturl string
    var err error
    for attempts < 5 {
        shorturl, err = utils.GenerateShortURL(in.Url, 10)
        if err != nil {
            continue
        }
		if err = r.Repositories.Url.SaveUrl(ctx, in.Url, shorturl); err == nil {
			break
		}
        attempts++
        r.Logger.Error("cannot get short url", slog.Any("error", err))
    }

	r.Logger.Debug(fmt.Sprintf("generated url: %s", shorturl))

	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}
	
	resp := &gen.RegUrlResp{
		ShortUrl: shorturl,
	}

	return resp, nil
}

func (r *UrlService) GetUrl(ctx context.Context, in *gen.GetUrlReq) (*gen.GetUrlResp, error) {
	r.Logger.Debug("Recived GetUrl")

	url, err := r.Repositories.Url.GetUrl(ctx, in.ShortUrl)
	if err != nil {
		r.Logger.Error("cannot get url", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	r.Logger.Debug(fmt.Sprintf("url: %s", url))

	response := &gen.GetUrlResp{
		Url: url,
	}
	return response, nil
}

func MustNewClient(s *grpc.Server, port string, logger *slog.Logger, repositories *storage.Repositories) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Failed to listen", slog.Any("Error: ", err))
		os.Exit(1)
	}
	srv := &UrlService{Logger: logger, Repositories: repositories}
	gen.RegisterUrlShortenerServer(s, srv)
	go func() {
		logger.Info("url-shortener service started", slog.String("addr", port))
		if err := s.Serve(lis); err != nil {
			logger.Error("Failed to serve", slog.Any("error", err))
		}
	}()
}
