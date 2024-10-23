package profile

import (
	"context"
	"github.com/taskemapp/server/apps/server/internal/grpc/interceptor"
	"github.com/taskemapp/server/apps/server/internal/service/profile"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Opts struct {
	fx.In
	Logger  *zap.Logger
	Profile profile.Service
}

func New(opts Opts) *Server {
	return &Server{profile: opts.Profile}
}

type Server struct {
	v1.UnimplementedProfileServer
	profile profile.Service
}

func (s *Server) AddOrUpdateAvatar(ctx context.Context, request *v1.AddOrUpdateAvatarRequest) (*emptypb.Empty, error) {
	a := request.AvatarImage

	uid, err := interceptor.GetUserID(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Wrong token provided")
	}

	err = s.profile.UploadAvatar(
		ctx,
		uid,
		profile.UploadAvatarOpts{
			Avatar: a,
		},
	)
	if err != nil {
		return nil, err
	}

	return nil, err
}
