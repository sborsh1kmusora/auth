package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sborsh1kmusora/auth/internal/model"
	"github.com/sborsh1kmusora/auth/internal/repository"
	"github.com/sborsh1kmusora/auth/internal/repository/user/converter"
	modelRepo "github.com/sborsh1kmusora/auth/internal/repository/user/model"
)

const (
	tableName = "user_info"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	isAdminColumn   = "is_admin"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &repo{pool}
}

func (r *repo) Create(ctx context.Context, in *model.CreateInput) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, isAdminColumn, passwordColumn).
		Values(in.UserInfo.Name, in.UserInfo.Email, in.UserInfo.IsAdmin, in.Password).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	log.Printf("SQL: %s | args: %v\n", query, args)

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(
		idColumn,
		nameColumn,
		emailColumn,
		isAdminColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, in *model.UpdateInput) error {
	updateBuilder := sq.Update("user_info").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": in.ID})

	if in.Name != nil {
		updateBuilder = updateBuilder.Set("name", in.Name)
	}

	if in.Email != nil {
		updateBuilder = updateBuilder.Set("email", in.Email)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete("user_info").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
