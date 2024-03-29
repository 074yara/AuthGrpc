package application

import (
	grpcapp "github.com/074yara/AuthGrpc/auth/internal/application/grpc"
	"github.com/074yara/AuthGrpc/auth/internal/domain/services/auth"
	"github.com/074yara/AuthGrpc/auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCServer: grpcApp}
}
