package application

import (
	grpcapp "github.com/074yara/AuthGrpc/auth/internal/application/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//TODO: Add storage
	//TODO: Init auth service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCServer: grpcApp}
}
