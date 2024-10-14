package interceptor

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/taskemapp/server/apps/server/internal/pkg/jwt"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CtxKey struct {
	key string
}

// Auth get token from grpc request metadata
//
// Already throws formated grpc with status.Errorf
func (i *Interceptor) Auth(ctx context.Context) (context.Context, error) {
	tokenMd, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	claims, err := jwt.GetPayload(tokenMd, i.c.TokenSecret)

	ctx, err = matchTokenErr(ctx, err)
	if err != nil {
		return ctx, err
	}

	payload := *claims

	_, err = i.tokenRepo.GetToken(ctx, fmt.Sprintf("%s:%s", payload["type"].(string), payload["uid"].(string)))
	if err != nil {
		switch {
		case errors.Is(err, token.ErrNotFound):
			return nil, status.Errorf(codes.InvalidArgument, "Wrong token provided")
		default:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	ctx, err = provideUserID(ctx, payload)

	return ctx, nil
}

func matchTokenErr(ctx context.Context, err error) (context.Context, error) {
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, status.Errorf(codes.Unauthenticated, "Token expired")
		case errors.Is(err, jwt.ErrTokenParse):
			return nil, status.Errorf(codes.InvalidArgument, "Token parse error")
		case errors.Is(err, jwt.ErrTokenValidation):
			return nil, status.Errorf(codes.Unauthenticated, "Token validation error")
		default:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}
	return ctx, nil
}
