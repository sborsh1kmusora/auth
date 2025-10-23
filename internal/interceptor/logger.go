package interceptor

import (
	"context"
	"time"

	"github.com/sborsh1kmusora/auth/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err == nil {
		logger.Debug("Success", zap.String("method", info.FullMethod), zap.Any("req", req), zap.Any("res", res), zap.Duration("duration", time.Since(now)))
	}

	return res, err
}