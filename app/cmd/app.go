package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage"
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
	

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	//s.Stop()
	storageClient.Close()

	logger.Info("Server stoped")
}