package user

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

const tableName = "users"

type Opts struct {
	fx.In
	Pgx *pgxpool.Pool
}

type Pgx struct {
	pgx *pgxpool.Pool
}

func NewPgx(opts Opts) (*Pgx, error) {
	return &Pgx{
		pgx: opts.Pgx,
	}, nil
}

var selectUserFields = []string{
	"id",
	"name",
	"display_name",
	"email",
	"avatar_url",
	"created_at",
	"edited_at",
}

func (p *Pgx) Update(ctx context.Context, userID uuid.UUID, opts UpdateOpts) (*User, error) {
	var updateMap = map[string]interface{}{}

	if opts.DisplayName != nil {
		updateMap["display_name"] = *opts.DisplayName
	}

	if opts.Email != nil {
		updateMap["email"] = *opts.DisplayName
	}

	if opts.AvatarUrl != nil {
		updateMap["avatar_url"] = *opts.DisplayName
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update(tableName).
		SetMap(updateMap).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = p.pgx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return p.FindByID(ctx, userID)
}

func (p *Pgx) FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error) {
	limit := opts.PerPage
	if limit == 0 {
		limit = 100
	}
	offset := opts.Page * limit

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From(tableName).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pgx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.DisplayName,
			&user.Email,
			&user.AvatarUrl,
			&user.CreatedAt,
			&user.EditedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
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
		Users:       users,
		Total:       total,
		HasNextPage: total > offset+limit,
	}, nil
}

func (p *Pgx) FindByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From(tableName).
		Where(squirrel.Eq{"id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user User
	err = p.pgx.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.EditedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) Create(ctx context.Context, opts CreateOpts) (*User, error) {
	fields := []string{
		"name",
		"display_name",
		"email",
		"avatar_url",
		"created_at",
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert(tableName).
		Columns(fields...).
		Values(
			opts.Name,
			opts.DisplayName,
			opts.Email,
			opts.AvatarUrl,
			time.Now(),
		).ToSql()
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

	var user User
	err = tx.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.EditedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) FindByName(ctx context.Context, name string) (*User, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(selectUserFields...).
		From(tableName).
		Where(squirrel.Eq{"name": name}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user User
	err = p.pgx.QueryRow(ctx, query, args).Scan(
		&user.ID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.EditedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Pgx) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete(tableName).
		Where(squirrel.Eq{"id": userID}).ToSql()
	if err != nil {
		return err
	}

	tag, err := p.pgx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
