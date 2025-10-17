package user

import (
	"github.com/sborsh1kmusora/auth/internal/service/user"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	userService user.Service
}

func NewImplementation(userService user.Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
