package user

import (
	"github.com/sborsh1kmusora/auth/internal/client/db"
	"github.com/sborsh1kmusora/auth/internal/repository"
	"github.com/sborsh1kmusora/auth/internal/service"
)

type serv struct {
	userRepo  repository.UserRepository
	txManager db.TxManager
}

func NewService(
	authRepo repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepo:  authRepo,
		txManager: txManager,
	}
}
