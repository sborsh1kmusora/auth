package service

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, in *model.CreateInput) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, in *model.UpdateInput) error
	Delete(ctx context.Context, id int64) error
}
