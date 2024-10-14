package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
)

func (i *Interceptor) AuthMatcher(ctx context.Context, c interceptors.CallMeta) bool {
	return loginSkip(ctx, c) &&
		refreshSkip(ctx, c) &&
		signUpSkip(ctx, c) &&
		reflectionSkip(ctx, c)
}

func loginSkip(_ context.Context, c interceptors.CallMeta) bool {
	return c.FullMethod() != "/v1.auth.Auth/Login"
}

func signUpSkip(_ context.Context, c interceptors.CallMeta) bool {
	return c.FullMethod() != "/v1.auth.Auth/SignUp"
}

func refreshSkip(_ context.Context, c interceptors.CallMeta) bool {
	return c.FullMethod() != "/v1.auth.Auth/RefreshToken"
}

func reflectionSkip(_ context.Context, c interceptors.CallMeta) bool {
	return c.FullMethod() != "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo"
}
