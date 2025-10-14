package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sborsh1kmusora/auth/internal/model"
	"github.com/sborsh1kmusora/auth/internal/repository"
	"github.com/sborsh1kmusora/auth/internal/repository/user/converter"
	modelRepo "github.com/sborsh1kmusora/auth/internal/repository/user/model"
	"github.com/sborsh1kmusora/platform_common/pkg/db"
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
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{
		db: db,
	}
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

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRow: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
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

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRow: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
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

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
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

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRow: query,
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}
