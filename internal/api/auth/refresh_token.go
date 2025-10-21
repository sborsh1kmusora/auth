package auth

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"go.uber.org/zap"
)

func (i *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	logger.Info("Getting new refresh token")

	refreshToken, err := i.authService.RefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		logger.Error("Failed to get refresh token", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully got new refresh token")

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
