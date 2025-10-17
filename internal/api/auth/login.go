package auth

import (
	"context"
	"log"

	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

func (i *Implementation) Login(
	ctx context.Context,
	req *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		log.Printf("failed to login: %v\n", err)
		return nil, err
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
