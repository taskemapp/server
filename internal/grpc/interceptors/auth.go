package interceptors

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"taskem-server/internal/pkg/jwt"
	"taskem-server/internal/repositories/token"
)

type TokenKey struct{}
type TokenPayload struct{}

// Auth get token from grpc request metadata
//
// Already throws formated grpc with status.Errorf
func (i *Interceptor) Auth(ctx context.Context) (context.Context, error) {
	tokenMd, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	claims, err := jwt.GetPayload(tokenMd, i.c.TokenSecret)

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

	payload := *claims

	_, err = i.tokenRepo.GetToken(ctx, fmt.Sprintf("%s:%s", payload["type"].(string), payload["uid"].(string)))
	if err != nil {
		switch {
		case errors.Is(err, token.ErrNotFound):
			return nil, status.Errorf(codes.Unauthenticated, "Wrong token")
		default:
			return nil, status.Errorf(codes.Internal, "Internal error")
		}
	}

	ctx = logging.InjectFields(ctx, logging.Fields{"auth.sub", payload})
	ctx = context.WithValue(ctx, TokenKey{}, tokenMd)
	ctx = context.WithValue(ctx, TokenPayload{}, payload)

	return ctx, nil
}
