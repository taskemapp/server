package grpc

import (
	"context"
	"errors"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"server/internal/pkg/jwt"
)

const authMDKey = "authorization"

// ExtractTokenPayload get token from grpc request metadata
//
// Already throws formated grpc with status.Errorf
func ExtractTokenPayload(ctx context.Context, secret string) (jwt2.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Failed to get metadata")
	}

	tokens := md[authMDKey]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}
	tokenStr := tokens[0]

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
