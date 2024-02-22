package main

import (
	"github.com/074yara/AuthGrpc/auth/internal/application"
	"github.com/074yara/AuthGrpc/auth/internal/config"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logger := setupLogger(cfg.Env)
	logger.Info("Starting auth service",
		slog.Any("config", cfg))
	app := application.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)

	go func() {
		if err := app.GRPCServer.Run(); err != nil {
			logger.Error("Failed to start gRPC server", slog.String("error", err.Error()))
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GRPCServer.Stop()
	logger.Info("gGRPC server gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
