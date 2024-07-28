package grpc

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authMDKey = "authorization"

func ExtractTokenPayload(ctx context.Context, secret string) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Failed to get metadata")
	}

	tokens := md[authMDKey]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}
	tokenStr := tokens[0]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token claims")
	}
}
