package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sborsh1kmusora/auth/internal/converter"
	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	logger.Info("Updating user with id", zap.Int64("id", req.GetId()))

	if err := i.userService.Update(ctx, converter.ToUpdateUserFromDesc(req)); err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error updating user")
	}

	logger.Info("Successfully updated user", zap.Int64("id", req.GetId()))

	return &empty.Empty{}, nil
}
