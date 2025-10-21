package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	logger.Info("Deleting user with id", zap.Int64("id", req.GetId()))

	if err := i.userService.Delete(ctx, req.Id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error deleting user")
	}

	logger.Info("Successfully deleted user", zap.Int64("id", req.GetId()))

	return &empty.Empty{}, nil
}
