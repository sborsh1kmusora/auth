package auth

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/utils"
)

func (s *serv) AccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.authConfig.RefreshTokenSecretKey())
	if err != nil {
		return "", err
	}

	user, err := s.userRepository.GetByUsername(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(user, s.authConfig.AccessTokenSecretKey(), s.authConfig.AccessTokenExpiration())
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
