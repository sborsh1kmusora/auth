package access

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/logger"
	desc "github.com/sborsh1kmusora/auth/pkg/access_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*desc.CheckResponse, error) {
	logger.Info("Checking access to endpoint", zap.String("endpoint", req.GetEndpointAddress()))

	isAllowed, err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		logger.Error("Error checking access to endpoint", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	logger.Info("Access checked", zap.String("endpoint", req.GetEndpointAddress()), zap.Bool("isAllowed", isAllowed))

	return &desc.CheckResponse{
		IsAllowed: isAllowed,
	}, nil
}
