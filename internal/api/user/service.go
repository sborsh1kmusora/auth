package user

import (
	"github.com/sborsh1kmusora/auth/internal/service"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
