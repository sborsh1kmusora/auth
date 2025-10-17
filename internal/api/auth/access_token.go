package auth

import (
	"context"
	"log"

	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

func (i *Implementation) GetAccessToken(
	ctx context.Context,
	req *desc.GetAccessTokenRequest,
) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.AccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		log.Printf("failed to get access token: %v", err)
		return nil, err
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
