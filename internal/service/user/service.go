package user

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/model"
	"github.com/sborsh1kmusora/auth/internal/repository/user"
	"github.com/sborsh1kmusora/platform_common/pkg/db"
)

type Service interface {
	Create(ctx context.Context, in *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, in *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
}

type serv struct {
	userRepo  user.Repository
	txManager db.TxManager
}

func NewService(
	authRepo user.Repository,
	txManager db.TxManager,
) Service {
	return &serv{
		userRepo:  authRepo,
		txManager: txManager,
	}
}
