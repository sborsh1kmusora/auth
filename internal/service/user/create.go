package user

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, in *model.CreateInput) (int64, error) {
	id, err := s.userRepo.Create(ctx, in)
	if err != nil {
		return 0, err
	}

	return id, nil
}
