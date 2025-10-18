package access

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/config"
	"github.com/sborsh1kmusora/auth/internal/repository/access"
)

type Service interface {
	Check(ctx context.Context, endpointAddress string) (bool, error)
}

type serv struct {
	authConfig config.AuthConfig

	accessRepo access.Repository
}

func NewService(authConfig config.AuthConfig, accessRepo access.Repository) Service {
	return &serv{
		authConfig: authConfig,
		accessRepo: accessRepo,
	}
}
