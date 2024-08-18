package team

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"taskem-server/internal/config"
	"taskem-server/internal/grpc"
	"taskem-server/internal/mapper"
	"taskem-server/internal/repositories/token"
	"taskem-server/internal/service"
	"taskem-server/internal/service/team"
	v1 "taskem-server/tools/gen/grpc/v1"
)

type Opts struct {
	fx.In
	Team      team.Service
	Config    config.Config
	Logger    *zap.Logger
	RedisRepo token.Repository
}

type Server struct {
	v1.UnimplementedTeamServer
	team      team.Service
	config    config.Config
	logger    *zap.Logger
	redisRepo token.Repository
}

func New(opts Opts) *Server {
	return &Server{
		team:      opts.Team,
		config:    opts.Config,
		logger:    opts.Logger,
		redisRepo: opts.RedisRepo,
	}
}

func (t *Server) Get(ctx context.Context, request *v1.GetTeamRequest) (*v1.TeamResponse, error) {
	payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret, t.redisRepo)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, err
	}

	fTeam, err := t.team.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	return mapper.ToTeamResponse(fTeam), nil
}

func (t *Server) GetUserTeams(ctx context.Context, empty *emptypb.Empty) (*v1.GetAllTeamsResponse, error) {
	payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret, t.redisRepo)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, err
	}

	page := payload["page"].(int)

	res, err := t.team.GetUserTeams(ctx, uid, service.PaginationOpts{
		Page:    page,
		PerPage: 30,
	})
	if err != nil {
		return nil, err
	}
	return mapper.ToGetAllTeamsResponse(&res.Teams), nil
}

func (t *Server) GetAllCanJoin(ctx context.Context, empty *emptypb.Empty) (*v1.GetAllTeamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Server) Create(ctx context.Context, request *v1.CreateTeamRequest) (*v1.CreateTeamResponse, error) {
	payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret, t.redisRepo)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, err
	}

	c, err := t.team.Create(ctx, team.CreateOpts{
		CreatorID:   uid,
		Name:        "",
		Description: "",
	})
	if err != nil {
		t.logger.Sugar().Error(err)
		switch {
		case errors.As(err, pgconn.PgError{}):
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return mapper.ToCreateTeamResponse(c), nil
}

func (t *Server) Join(ctx context.Context, request *v1.JoinTeamRequest) (*v1.JoinTeamResponse, error) {
	panic("")
	//payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var uid uuid.UUID
	//uid, err = uuid.Parse(payload["uid"].(string))
	//if err != nil {
	//	return nil, err
	//}
	//
	//teamID, err := uuid.Parse(request.TeamId)
	//
	//err = t.team.Join(ctx, team.JoinOpts{
	//	UserID: uid,
	//	TeamID: teamID,
	//})
	//if err != nil {
	//	return nil,
	//}
	//
	//return
}

func (t *Server) GetRoles(ctx context.Context, request *v1.GetTeamRolesRequest) (*v1.GetTeamRolesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Server) ChangeRole(ctx context.Context, role *v1.ChangeTeamRole) (*v1.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Server) Leave(ctx context.Context, request *v1.LeaveTeamRequest) (*v1.LeaveTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}
