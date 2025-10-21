package user

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/converter"
	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting user with id", zap.Int64("id", req.Id))

	userObj, err := i.userService.Get(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting user")
	}

	logger.Info("Successfully got user", zap.Int64("id", req.Id))

	return &desc.GetResponse{
		User: converter.ToDescFromUser(userObj),
	}, nil
}
