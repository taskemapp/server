package interceptor

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/taskemapp/server/apps/server/internal/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func provideUserID(ctx context.Context, payload jwt.Claims) (context.Context, error) {
	uid, err := uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Wrong token")
	}

	ctx = logging.InjectFields(
		ctx,
		logging.Fields{"user_id", uid.String()},
	)

	return context.WithValue(ctx, CtxKey{"uid"}, uid), nil
}

// GetUserID from context, only throws ErrGetUserID if uid not found in context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if uid, ok := ctx.Value(CtxKey{"uid"}).(uuid.UUID); ok {
		return uid, nil
	}

	return uuid.Nil, ErrGetUserID
}

func provideReqID(ctx context.Context) (context.Context, error) {
	rid, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate rid")
	}

	ctx = logging.InjectFields(
		ctx,
		logging.Fields{"request_id", rid.String()},
	)

	return context.WithValue(ctx, CtxKey{"rid"}, rid), nil
}

// GetRequestID from context, only throws ErrRequestID if uid not found in context
func GetRequestID(ctx context.Context) (uuid.UUID, error) {
	if uid, ok := ctx.Value(CtxKey{"uid"}).(uuid.UUID); ok {
		return uid, nil
	}

	return uuid.Nil, ErrRequestID
}
