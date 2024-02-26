package grpcapp

import (
	"fmt"
	"github.com/074yara/AuthGrpc/auth/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService auth.AuthService, port int) *App {
	gRPCServer := grpc.NewServer()
	auth.Register(gRPCServer, authService)
	return &App{
		log:        log,
		GRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))

	if err = a.GRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(
		slog.String("op", op)).
		Info(
			"stopping gRPC server", slog.Int("port", a.port),
		)
	a.GRPCServer.GracefulStop()
}
