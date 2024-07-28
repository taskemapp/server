package team_member

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

func (p *Pgx) Create(ctx context.Context, opts CreateOpts) (*TeamMember, error) {
	fields := []string{
		"id",
		"user_id",
		"team_id",
		"joined_at",
		"leaved_at",
		"is_leaved",
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
		return nil, err
	}

	return &tm, nil
}

func (p *Pgx) Update(ctx context.Context, tmID uuid.UUID, opts UpdateOpts) (*TeamMember, error) {
	//TODO implement me
	panic("implement me")
}
