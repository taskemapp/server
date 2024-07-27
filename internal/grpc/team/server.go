package team

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/emptypb"
	"taskem-server/internal/config"
	"taskem-server/internal/grpc"
	"taskem-server/internal/mapper"
	"taskem-server/internal/repositories/team"
	v1 "taskem-server/tools/gen/grpc/v1"
)

type Opts struct {
	fx.In
	TeamRepo team.Repository
	Config   config.Config
}

type Team struct {
	v1.UnimplementedTeamServer
	teamRepo team.Repository
	config   config.Config
}

func (t *Team) Get(ctx context.Context, request *v1.GetTeamRequest) (*v1.TeamResponse, error) {
	payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, err
	}

	fTeam, err := t.teamRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return mapper.ToTeamResponse(fTeam), nil
}

func (t *Team) GetUserTeams(ctx context.Context, empty *emptypb.Empty) (*v1.GetAllTeamsResponse, error) {
	payload, err := grpc.ExtractTokenPayload(ctx, t.config.TokenSecret)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.Parse(payload["uid"].(string))
	if err != nil {
		return nil, err
	}

	page := payload["page"].(int)

	res, err := t.teamRepo.FindManyBelongToUser(ctx, uid, team.FindManyOpts{
		Page:    page,
		PerPage: 30,
	})
	if err != nil {
		return nil, err
	}
	return mapper.ToGetAllTeamsResponse(&res.Teams), nil
}

func (t *Team) GetAllCanJoin(ctx context.Context, empty *emptypb.Empty) (*v1.GetAllTeamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Team) Create(ctx context.Context, request *v1.CreateTeamRequest) (*v1.CreateTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Team) Join(ctx context.Context, request *v1.JoinTeamRequest) (*v1.JoinTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Team) GetRoles(ctx context.Context, request *v1.GetTeamRolesRequest) (*v1.GetTeamRolesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Team) ChangeRole(ctx context.Context, role *v1.ChangeTeamRole) (*v1.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Team) Leave(ctx context.Context, request *v1.LeaveTeamRequest) (*v1.LeaveTeamResponse, error) {
	//TODO implement me
	panic("implement me")
}
