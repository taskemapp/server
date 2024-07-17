package team

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const tableName = "teams"

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
	"name",
	"description",
	"header_image_url",
	"creator",
	"created_at",
	"edited_at",
}

func (p Pgx) Create(ctx context.Context, opts CreateOpts) (*Team, error) {
	fields := []string{
		"name",
		"description",
		"creator",
		"created_at",
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert(tableName).
		Columns(fields...).
		Values(
			opts.Name,
			opts.Description,
			opts.Creator,
			time.Now().Unix(),
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

	var team Team
	err = tx.QueryRow(ctx, query, args...).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.HeaderImageUrl,
		&team.Creator,
		&team.CreatedAt,
		&team.EditedAt,
	)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (p Pgx) FindByID(ctx context.Context, teamID uuid.UUID) (*Team, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTeamFields...).
		From(tableName).
		Where(squirrel.Eq{"id": teamID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var team Team
	err = p.pgx.QueryRow(ctx, query, args...).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.HeaderImageUrl,
		&team.Creator,
		&team.CreatedAt,
		&team.EditedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &team, nil
}

func (p Pgx) FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error) {
	limit := opts.PerPage
	if limit == 0 {
		limit = 100
	}
	offset := opts.Page * limit

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTeamFields...).
		From(tableName).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pgx.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		team := Team{}
		err = rows.Scan(
			&team.ID,
			&team.Name,
			&team.Description,
			&team.HeaderImageUrl,
			&team.Creator,
			&team.CreatedAt,
			&team.EditedAt,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	totalQuery, totalArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("COUNT(*)").
		From(tableName).
		ToSql()
	if err != nil {
		return nil, err
	}

	var total int
	err = p.pgx.QueryRow(ctx, totalQuery, totalArgs...).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &FindManyResult{
		Teams:       teams,
		Total:       total,
		HasNextPage: total > offset+limit,
	}, nil
}

func (p Pgx) FindManyBelongToUser(ctx context.Context, userID uuid.UUID, opts FindManyOpts) (*FindManyResult, error) {
	//limit := opts.PerPage
	//if limit == 0 {
	//	limit = 100
	//}
	//offset := opts.Page * limit

	//query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
	//	Select(selectTeamFields...).
	//	From(tableName).
	//	Where(squirrel.Eq{"team_id": userID}).
	//	Offset(uint64(offset)).
	//	Limit(uint64(limit)).
	//	ToSql()
	//if err != nil {
	//	return nil, err
	//}
	//
	//rows, err := p.pgx.Query(ctx, query, args...)
	//if err != nil {
	//	if errors.Is(err, pgx.ErrNoRows) {
	//		return nil, ErrNotFound
	//	}
	//	return nil, err
	//}
	//
	//defer rows.Close()
	//
	//teams := make([]Team, 0)
	//for rows.Next() {
	//	team := Team{}
	//	err = rows.Scan(
	//		&team.ID,
	//		&team.Name,
	//		&team.Description,
	//		&team.HeaderImageUrl,
	//		&team.Creator,
	//		&team.CreatedAt,
	//		&team.EditedAt,
	//	)
	//	if err != nil {
	//		return nil, err
	//	}
	//	teams = append(teams, team)
	//}
	//
	//totalQuery, totalArgs, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
	//	Select("COUNT(*)").
	//	From(tableName).
	//	ToSql()
	//if err != nil {
	//	return nil, err
	//}
	//
	//var total int
	//err = p.pgx.QueryRow(ctx, totalQuery, totalArgs...).Scan(&total)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &FindManyResult{
	//	Teams:       teams,
	//	Total:       total,
	//	HasNextPage: total > offset+limit,
	//}, nil
	panic("")
}

func (p Pgx) Update(ctx context.Context, teamID uuid.UUID, opts UpdateOpts) (*Team, error) {
	var updateMap = map[string]interface{}{}

	if opts.Name != nil {
		updateMap["name"] = *opts.Name
	}

	if opts.Description != nil {
		updateMap["description"] = *opts.Description
	}

	if opts.HeaderImage != nil {
		updateMap["header_image_url"] = *opts.HeaderImage
	}

	if opts.Creator != nil {
		updateMap["creator"] = *opts.Creator
	}

	updateMap["edited_at"] = time.Now()

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update(tableName).
		SetMap(updateMap).
		Where(squirrel.Eq{"id": teamID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = p.pgx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return p.FindByID(ctx, teamID)
}
