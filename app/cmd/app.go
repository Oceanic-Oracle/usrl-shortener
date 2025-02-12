package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/api"
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
	
	s := grpc.NewServer()
	api.MustNewClient(s, cfg.Port, logger, storageClient)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	storageClient.Close()
	s.GracefulStop()

	logger.Info("Server stoped")
}