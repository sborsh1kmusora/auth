package auth

import (
	"context"
	"errors"

	appError "github.com/sborsh1kmusora/auth/internal/errors"
	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Login(
	ctx context.Context,
	req *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	logger.Info("Login new user", zap.String("username", req.GetUsername()))

	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrUserNotFound):
			logger.Warn("User not found", zap.String("username", req.GetUsername()))
			return nil, status.Errorf(codes.NotFound, "user %q not found", req.GetUsername())
		case errors.Is(err, appError.ErrInvalidCredentials):
			logger.Warn("Invalid credentials", zap.String("username", req.GetUsername()))
			return nil, status.Errorf(codes.PermissionDenied, "invalid credentials")
		default:
			logger.Error("Failed to login", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	logger.Info("User successfully logged in", zap.String("username", req.GetUsername()))

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
