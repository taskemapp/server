package e2e_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"github.com/taskemapp/server/libs/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	qCfg, _      = queue.NewConfig()
	s3Cfg, _     = s3.NewConfig()
	cfg, _       = config.New(qCfg, s3Cfg)
	testEmail    = gofakeit.Email()
	testPassword = gofakeit.Password(true, true, true, true, false, 8)
)

func TestSignUp(t *testing.T) {
	SkipCI(t)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(fmt.Sprint("0.0.0.0:", 50051), opts...)
	defer conn.Close()

	require.NoError(t, err)

	client := v1.NewAuthClient(conn)

	_, err = client.SignUp(context.Background(), &v1.SignupRequest{
		UserName: gofakeit.Username(),
		Email:    testEmail,
		Password: testPassword,
	})
	require.NoError(t, err)
}

func TestLogin(t *testing.T) {
	SkipCI(t)
	conn, err := grpc.NewClient(
		fmt.Sprintf("0.0.0.0:%d", cfg.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close()

	require.NoError(t, err)

	client := v1.NewAuthClient(conn)

	resp, err := client.Login(context.Background(), &v1.LoginRequest{
		Email:    testEmail,
		Password: testPassword,
	})
	require.NoError(t, err)
	t.Log(resp)
}
