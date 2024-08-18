package grpc

import (
	"context"
	"errors"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"taskem-server/internal/pkg/jwt"
)

const authMDKey = "authorization"

func ExtractToken(ctx context.Context) (*string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Failed to get metadata")
	}

	tokens := md[authMDKey]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}
	tokenStr := tokens[0]

	// Delete "Bearer" from token
	tokenStr = tokenStr[7:]

	return &tokenStr, nil
}

// ExtractTokenPayload get token from grpc request metadata
//
// Already throws formated grpc with status.Errorf
func ExtractTokenPayload(
	ctx context.Context,
	secret string,
	tokenOpt ...*string,
) (jwt2.MapClaims, error) {
	var tokenStr string
	if len(tokenOpt) > 0 && tokenOpt[0] != nil {
		tokenStr = *tokenOpt[0]
	} else {
		token, err := ExtractToken(ctx)
		if err != nil {
			return nil, err
		}
		tokenStr = *token
	}

	payload, err := jwt.GetPayload(tokenStr, secret)

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
	return *payload, nil
}
