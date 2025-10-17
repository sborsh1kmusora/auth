package user

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, in *model.UpdateUser) error {
	if err := s.userRepo.Update(ctx, in); err != nil {
		return err
	}
	return nil
}
