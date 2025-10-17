package auth

import (
	"context"
	"log"

	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.RefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		log.Printf("failed to get refresh token: %v", err)
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
