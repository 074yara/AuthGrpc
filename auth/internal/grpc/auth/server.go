package auth

import (
	"context"
	"github.com/074yara/AuthGrpc/protos/gen/authGrpc"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	authGrpc.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authGrpc.RegisterAuthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) Register(ctx context.Context, req *authGrpc.RegisterRequest) (*authGrpc.RegisterResponse, error) {
	return nil, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *authGrpc.LoginRequest) (*authGrpc.LoginResponse, error) {
	return &authGrpc.LoginResponse{Token: "hallo"}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *authGrpc.IsAdminRequest) (*authGrpc.IsAdminResponse, error) {
	return nil, nil
}
