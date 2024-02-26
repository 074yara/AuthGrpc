package auth

import (
	"context"
	"github.com/074yara/AuthGrpc/protos/gen/authGrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type AuthService interface {
	Login(ctx context.Context, email, password string, appId uint) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (uint, error)
	IsAdmin(ctx context.Context, id uint) (bool, error)
}

type ServerAPI struct {
	authGrpc.UnimplementedAuthServer
	auth AuthService
}

func Register(gRPC *grpc.Server, auth AuthService) {
	authGrpc.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Register(ctx context.Context, req *authGrpc.RegisterRequest) (*authGrpc.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}
	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//TODO: check for error type: register of already existing user. return more informative error
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authGrpc.RegisterResponse{UserId: uint64(userId)}, nil
}

func (s *ServerAPI) Login(ctx context.Context, req *authGrpc.LoginRequest) (*authGrpc.LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}
	tokenString, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), uint(req.GetAppId()))
	if err != nil {
		//TODO: check for error type: incorrect email/password. return more informative error
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authGrpc.LoginResponse{Token: tokenString}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *authGrpc.IsAdminRequest) (*authGrpc.IsAdminResponse, error) {
	if err := validateIsAdmin(uint(req.GetUserId())); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, uint(req.GetUserId()))
	if err != nil {
		//TODO: check for error type: return more informative error
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authGrpc.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLoginRequest(request *authGrpc.LoginRequest) error {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}
	if request.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app id is empty")
	}
	return nil
}

func validateRegisterRequest(request *authGrpc.RegisterRequest) error {
	//TODO: validate email to RFC 5322 standard
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}
	return nil
}

func validateIsAdmin(userId uint) error {
	if userId == emptyValue {
		return status.Error(codes.InvalidArgument, "user id is empty")
	}
	return nil
}
