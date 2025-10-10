package user

import (
	"github.com/sborsh1kmusora/auth/internal/repository"
	"github.com/sborsh1kmusora/auth/internal/service"
)

type serv struct {
	userRepo repository.UserRepository
}

func NewService(authRepo repository.UserRepository) service.UserService {
	return &serv{
		userRepo: authRepo,
	}
}
