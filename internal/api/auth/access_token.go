package auth

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"go.uber.org/zap"
)

func (i *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	logger.Info("Getting access token")

	accessToken, err := i.authService.AccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		logger.Error("Failed to get access token", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully got access token")

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
