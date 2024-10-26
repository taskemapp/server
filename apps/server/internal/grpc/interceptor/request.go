package interceptor

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ProvideRID provide request id for every single request,
// request id need for logging
func (i *Interceptor) ProvideRID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := provideReqID(ctx)
		if err != nil {
			i.logger.Warn("Failed to provide request id", zap.Error(err))
		}
		return handler(ctx, req)
	}
}
