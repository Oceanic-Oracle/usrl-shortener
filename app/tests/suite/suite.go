package suite

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"
	"url-shortener/internal/api"
	"url-shortener/internal/api/gen"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg       *config.Config
	UrlClient gen.UrlShortenerClient
	Storage   *storage.Repositories
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad("../.env")
	logger := logger.SetupLogger(cfg.Env)

	ctxReg, cancelReg := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelReg()

	srv, storageClient := regServer(ctxReg, cfg, logger)
	t.Cleanup(func() {
		srv.Stop()
	})

	ctxClient, cancelClient := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(func() {
		cancelClient()
	})

	connService, err := grpc.NewClient(cfg.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Auth gRPC server: %v", err)
	}

	t.Cleanup(func() {
		if err := connService.Close(); err != nil {
			t.Errorf("client connection close error: %v", err)
		}
	})

	return ctxClient, &Suite{
		T:         t,
		Cfg:       cfg,
		UrlClient: gen.NewUrlShortenerClient(connService),
		Storage:   storageClient,
	}
}

func regServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*grpc.Server, *storage.Repositories) {
	s := grpc.NewServer()

	storageClient := storage.NewClient(ctx, &cfg.Storage, logger)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		logger.Error("Failed to listen", slog.Any("Error: ", err))
		os.Exit(1)
	}

	srv := &api.UrlService{Logger: logger, Repositories: storageClient}
	gen.RegisterUrlShortenerServer(s, srv)
	go func() {
		logger.Info("url-shortener service started", slog.String("addr", cfg.Port))
		if err := s.Serve(lis); err != nil {
			logger.Error("Failed to serve", slog.Any("error", err))
		}
	}()

	return s, storageClient
}
