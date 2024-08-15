package team_member

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const tableName = "team_members"

type Opts struct {
	fx.In
	Pgx    *pgxpool.Pool
	Logger *zap.Logger
}

type Pgx struct {
	pgx    *pgxpool.Pool
	logger *zap.Logger
}

func NewPgx(opts Opts) (*Pgx, error) {
	return &Pgx{
		pgx:    opts.Pgx,
		logger: opts.Logger,
	}, nil
}

var selectTeamFields = []string{
	"id",
	"user_id",
	"team_id",
	"joined_at",
	"leaved_at",
	"is_leaved",
}

func (p *Pgx) FindByID(ctx context.Context, tmID uuid.UUID) (*TeamMember, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTeamFields...).
		From(tableName).
		Where(squirrel.Eq{"id": tmID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var tm TeamMember
	err = p.pgx.QueryRow(ctx, query, args...).Scan(
		&tm.ID,
		&tm.UserID,
		&tm.TeamID,
		&tm.JoinedAt,
		&tm.LeavedAt,
		&tm.IsLeaved,
	)
	if err != nil {
		p.logger.Sugar().Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &tm, nil
}

func (p *Pgx) FindByUserAndTeam(ctx context.Context, userID uuid.UUID, teamID uuid.UUID) (*TeamMember, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTeamFields...).
		From(tableName).
		Where(squirrel.Eq{
			"user_id": userID,
			"team_id": teamID,
		}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var tm TeamMember
	err = p.pgx.QueryRow(ctx, query, args...).Scan(
		&tm.ID,
		&tm.UserID,
		&tm.TeamID,
		&tm.JoinedAt,
		&tm.LeavedAt,
		&tm.IsLeaved,
	)
	if err != nil {
		p.logger.Sugar().Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &tm, nil
}

func (p *Pgx) Create(ctx context.Context, opts CreateOpts) (*TeamMember, error) {
	fields := []string{
		"user_id",
		"team_id",
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert(tableName).
		Columns(fields...).
		Values(
			opts.UserID,
			opts.TeamID,
		).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := p.pgx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			p.logger.Sugar().Error(err)
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var tm TeamMember
	err = tx.QueryRow(ctx, query, args...).Scan(
		&tm.ID,
		&tm.UserID,
		&tm.TeamID,
		&tm.JoinedAt,
		&tm.LeavedAt,
		&tm.IsLeaved,
	)
	if err != nil {
		p.logger.Sugar().Error(err)
		return nil, err
	}

	return &tm, nil
}

func (p *Pgx) Update(ctx context.Context, tmID uuid.UUID, opts UpdateOpts) (*TeamMember, error) {
	var updateMap = map[string]interface{}{}

	if opts.IsLeaved != nil {
		updateMap["is_leaved"] = *opts.IsLeaved
	}

	updateMap["leaved_at"] = time.Now()

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update(tableName).
		SetMap(updateMap).
		Where(squirrel.Eq{"id": tmID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = p.pgx.Exec(ctx, query, args...)
	if err != nil {
		p.logger.Sugar().Error(err)
		return nil, err
	}

	return p.FindByID(ctx, tmID)
}
