package task

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

const tableName = "tasks"

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

var selectTaskFields = []string{
	"id",
	"name",
	"description",
	"status",
	"team_id",
	"assigned_user_id",
	"creator",
	"created_at",
	"edited_at",
	"end_at",
}

func (p Pgx) Create(ctx context.Context, opts CreateOpts) (*Task, error) {
	fields := []string{
		"name",
		"description",
		"status",
		"team_id",
		"assigned_user_id",
		"creator",
		"created_at",
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert(tableName).
		Columns(fields...).
		Values(
			opts.Name,
			opts.Description,
			opts.Status,
			opts.TeamID,
			opts.AssignedUserID,
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

	var task Task
	err = tx.QueryRow(ctx, query, args...).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.TeamID,
		&task.AssignedUserID,
		&task.Creator,
		&task.CreatedAt,
		&task.EditedAt,
		&task.EndAt,
	)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (p Pgx) FindByID(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTaskFields...).
		From(tableName).
		Where(squirrel.Eq{"id": taskID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var task Task
	err = p.pgx.QueryRow(ctx, query, args...).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.TeamID,
		&task.AssignedUserID,
		&task.Creator,
		&task.CreatedAt,
		&task.EditedAt,
		&task.EndAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (p Pgx) FindByStatusAndUserID(ctx context.Context, userID uuid.UUID, status string) (*[]Task, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTaskFields...).
		From(tableName).
		Where(squirrel.Eq{
			"assigned_user_id": userID,
			"status":           status,
		}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)

	rows, err := p.pgx.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.TeamID,
			&task.AssignedUserID,
			&task.Creator,
			&task.CreatedAt,
			&task.EditedAt,
			&task.EndAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return &tasks, nil
}

func (p Pgx) FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error) {
	limit := opts.PerPage
	if limit == 0 {
		limit = 100
	}
	offset := opts.Page * limit

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTaskFields...).
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

	tasks := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.TeamID,
			&task.AssignedUserID,
			&task.Creator,
			&task.CreatedAt,
			&task.EditedAt,
			&task.EndAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
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
		Tasks:       tasks,
		Total:       total,
		HasNextPage: total > offset+limit,
	}, nil
}

func (p Pgx) FindManyBelongToTeam(ctx context.Context, teamID uuid.UUID, opts FindManyOpts) (*FindManyResult, error) {
	limit := opts.PerPage
	if limit == 0 {
		limit = 100
	}
	offset := opts.Page * limit

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTaskFields...).
		From(tableName).
		Where(squirrel.Eq{"team_id": teamID}).
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

	tasks := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.TeamID,
			&task.AssignedUserID,
			&task.Creator,
			&task.CreatedAt,
			&task.EditedAt,
			&task.EndAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
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
		Tasks:       tasks,
		Total:       total,
		HasNextPage: total > offset+limit,
	}, nil
}

func (p Pgx) FindManyBelongToUser(ctx context.Context, userID uuid.UUID, opts FindManyOpts) (*FindManyResult, error) {
	limit := opts.PerPage
	if limit == 0 {
		limit = 100
	}
	offset := opts.Page * limit

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectTaskFields...).
		From(tableName).
		Where(squirrel.Eq{"assigned_user_id": userID}).
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

	tasks := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.TeamID,
			&task.AssignedUserID,
			&task.Creator,
			&task.CreatedAt,
			&task.EditedAt,
			&task.EndAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
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
		Tasks:       tasks,
		Total:       total,
		HasNextPage: total > offset+limit,
	}, nil
}

func (p Pgx) Update(ctx context.Context, taskID uuid.UUID, opts UpdateOpts) (*Task, error) {
	var updateMap = map[string]interface{}{}

	if opts.Name != nil {
		updateMap["name"] = *opts.Name
	}

	if opts.Description != nil {
		updateMap["description"] = *opts.Description
	}

	if opts.Status != nil {
		updateMap["status"] = *opts.Status
	}

	if opts.EndAt != nil {
		updateMap["end_at"] = *opts.EndAt
	}

	updateMap["edited_at"] = time.Now()

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update(tableName).
		SetMap(updateMap).
		Where(squirrel.Eq{"id": taskID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = p.pgx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return p.FindByID(ctx, taskID)
}
