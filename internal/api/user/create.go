package user

import (
	"context"
	"errors"

	"github.com/sborsh1kmusora/auth/internal/converter"
	appError "github.com/sborsh1kmusora/auth/internal/errors"
	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Creating new user")

	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req.GetUserInfo()))
	if err != nil {
		if errors.Is(err, appError.ErrUserAlreadyExists) {
			logger.Warn("User already exists", zap.String("username", req.GetUserInfo().GetUsername()))
			return nil, status.Errorf(codes.AlreadyExists, "user with username %q already exists", req.GetUserInfo().GetUsername())
		}

		logger.Error("Failed to create user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	logger.Info("Successfully created user", zap.Any("id", id))

	return &desc.CreateResponse{Id: id}, nil
}
