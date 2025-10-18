package auth

import (
	"github.com/sborsh1kmusora/auth/internal/service/auth"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService auth.Service
}

func NewImplementation(authService auth.Service) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
