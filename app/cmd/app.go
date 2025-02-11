package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/api"
	"url-shortener/internal/api/gen"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage"

	"google.golang.org/grpc"
)

const (
	envPath string = "../.env"
)

func main() {
	cfg := config.MustLoad(envPath)

	logger := logger.SetupLogger(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	storageClient := storage.NewClient(ctx, &cfg.Storage, logger)
	
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		logger.Error("Failed to listen", slog.Any("Error: ", err))
		return
	}
	s := grpc.NewServer()
	srv := api.NewClient(logger, storageClient)
	gen.RegisterUrlShortenerServer(s, srv)
	go func() {
		logger.Info("AuthService started", slog.String("addr", cfg.Port))
		if err := s.Serve(lis); err != nil {
			logger.Error("Failed to serve", slog.Any("error", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	s.Stop()
	storageClient.Close()

	logger.Info("Server stoped")
}