package authgrpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IskanderSh/taqwa-auth/internal/services/auth"
	authv1 "github.com/IskanderSh/taqwa-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
	log  *slog.Logger
}

func Register(gRPC *grpc.Server, auth Auth, log *slog.Logger) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, log: log})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	s.log.Debug("successfully validate login with params: %s, %s", req.GetEmail(), req.GetPassword())

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	s.log.Debug("get token: %s, err: %w", token, err)
	if err != nil {
		s.log.Debug("handling error: %w", err)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "arguments incorrect")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	s.log.Debug("return login response")
	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	s.log.Debug("successfully validate register func with params: %s, %s", req.GetEmail(), req.GetPassword())

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	s.log.Debug("get userID: %s, err: %w", userID, err)
	if err != nil {
		s.log.Debug("handling error: %w", err)
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	s.log.Debug("return register response")
	return &authv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
