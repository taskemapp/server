package team

import (
	"context"
	"github.com/google/uuid"
	"github.com/taskemapp/server/apps/server/internal/repositories/team"
	"github.com/taskemapp/server/apps/server/internal/repositories/team_member"
	"github.com/taskemapp/server/apps/server/internal/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In
	TeamRepo       team.Repository
	TeamMemberRepo team_member.Repository
	Logger         *zap.Logger
}

type Team struct {
	teamRepo       team.Repository
	teamMemberRepo team_member.Repository
	logger         *zap.Logger
}

func New(opts Opts) *Team {
	return &Team{
		teamRepo:       opts.TeamRepo,
		teamMemberRepo: opts.TeamMemberRepo,
		logger:         opts.Logger,
	}
}

func (t *Team) Get(ctx context.Context, id uuid.UUID) (*team.Team, error) {
	f, err := t.teamRepo.FindByID(ctx, id)
	if err != nil {
		t.logger.Sugar().Error(err)
		return nil, err
	}
	return f, nil
}

func (t *Team) GetUserTeams(ctx context.Context, userID uuid.UUID, pgOpts service.PaginationOpts) (*team.FindManyResult, error) {
	fMany, err := t.teamRepo.FindManyBelongToUser(ctx, userID, team.FindManyOpts{
		Page:    pgOpts.Page,
		PerPage: pgOpts.PerPage,
	})
	if err != nil {
		t.logger.Sugar().Error(err)
		return nil, err
	}
	return fMany, nil
}

func (t *Team) Create(ctx context.Context, opts CreateOpts) (*team.Team, error) {
	c, err := t.teamRepo.Create(ctx, team.CreateOpts{
		Name:        opts.Name,
		Description: opts.Description,
		Creator:     opts.CreatorID,
	})
	if err != nil {
		t.logger.Sugar().Error(err)
		return nil, err
	}
	return c, nil
}

func (t *Team) Join(ctx context.Context, opts JoinOpts) error {
	_, err := t.teamMemberRepo.Create(ctx, team_member.CreateOpts{
		UserID: opts.UserID,
		TeamID: opts.TeamID,
	})
	if err != nil {
		t.logger.Sugar().Error(err)
		return err
	}
	return err
}

func (t *Team) Leave(ctx context.Context, opts LeaveOpts) error {
	tm, err := t.teamMemberRepo.FindByUserAndTeam(ctx, opts.UserID, opts.TeamID)
	if err != nil {
		t.logger.Sugar().Error(err)
		return err
	}

	leaved := true
	_, err = t.teamMemberRepo.Update(
		ctx,
		tm.ID,
		team_member.UpdateOpts{
			IsLeaved: &leaved,
		},
	)
	if err != nil {
		t.logger.Sugar().Error(err)
		return err
	}

	return err
}

func (t *Team) GetRoles(ctx context.Context, teamID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
