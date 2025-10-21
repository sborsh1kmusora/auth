package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sborsh1kmusora/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, username, password string) (string, error) {
	userInfo, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if err := utils.VerifyPassword(userInfo.Password, password); err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(userInfo,
		s.authConfig.RefreshTokenSecretKey(),
		s.authConfig.RefreshTokenExpiration())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
