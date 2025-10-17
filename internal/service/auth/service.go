package auth

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/config"
	userRepo "github.com/sborsh1kmusora/auth/internal/repository/user"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	AccessToken(ctx context.Context, refreshToken string) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type serv struct {
	authConfig     config.AuthConfig
	userRepository userRepo.Repository
}

func NewService(authConfig config.AuthConfig, userRepository userRepo.Repository) Service {
	return &serv{
		authConfig:     authConfig,
		userRepository: userRepository,
	}
}
