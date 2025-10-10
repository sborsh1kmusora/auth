package user

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, in *model.CreateInput) (int64, error) {
	if in.Password != in.PasswordConfirm {
		return 0, errors.New("passwords do not match")
	}

	id, err := s.userRepo.Create(ctx, in)
	if err != nil {
		return 0, err
	}

	return id, nil
}
