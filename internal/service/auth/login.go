package auth

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/sborsh1kmusora/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, username, password string) (string, error) {
	userInfo, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("faield to get user by username: %v\n", err)
		return "", err
	}

	if !utils.VerifyPassword(userInfo.Password, password) {
		return "", errors.Errorf("invalid password")
	}

	token, err := utils.GenerateToken(userInfo,
		s.authConfig.RefreshTokenSecretKey(),
		s.authConfig.RefreshTokenExpiration())
	if err != nil {
		log.Printf("faield to generate token: %v\n", err)
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
