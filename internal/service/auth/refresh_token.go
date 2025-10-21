package auth

import (
	"context"
	"log"

	"github.com/sborsh1kmusora/auth/internal/utils"
)

func (s *serv) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.authConfig.RefreshTokenSecretKey())
	if err != nil {
		return "", err
	}

	user, err := s.userRepository.GetByUsername(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	newRefreshToken, err := utils.GenerateToken(
		user,
		s.authConfig.RefreshTokenSecretKey(),
		s.authConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		log.Printf("failed to generate new refresh token: %v", err)
		return "", err
	}

	return newRefreshToken, nil
}
