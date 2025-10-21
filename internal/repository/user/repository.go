package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	appError "github.com/sborsh1kmusora/auth/internal/errors"
	"github.com/sborsh1kmusora/auth/internal/logger"
	"github.com/sborsh1kmusora/auth/internal/model"
	"github.com/sborsh1kmusora/auth/internal/utils"
	"github.com/sborsh1kmusora/platform_common/pkg/db"
	"github.com/sborsh1kmusora/platform_common/pkg/db/prettier"
	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, in *model.UserInfo) (int64, error)
	GetByUsername(ctx context.Context, username string) (*model.UserInfo, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
}

const (
	tableName = "users"

	idColumn        = "id"
	usernameColumn  = "username"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, in *model.UserInfo) (int64, error) {
	hashPassword, err := utils.HashPassword(in.Password)
	if err != nil {
		return 0, err
	}

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, passwordColumn, roleColumn, createdAtColumn).
		Values(in.Username, hashPassword, in.Role, time.Now()).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRow: query,
	}

	logger.Debug("Executing query",
		zap.String("sql", q.Name),
		zap.String("query", prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)),
	)

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, appError.ErrUserAlreadyExists
			}
		}
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(
		idColumn,
		usernameColumn,
		passwordColumn,
		roleColumn,
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

	logger.Debug("Executing query",
		zap.String("sql", q.Name),
		zap.String("query", prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)),
	)

	user := new(model.User)
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) GetByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	builder := sq.Select(usernameColumn, passwordColumn, roleColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{usernameColumn: username})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "users.GetByUsername",
		QueryRow: query,
	}

	logger.Debug("Executing query",
		zap.String("sql", q.Name),
		zap.String("query", prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)),
	)

	user := new(model.UserInfo)
	err = r.db.DB().ScanOneContext(ctx, user, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) Update(ctx context.Context, in *model.UpdateUser) error {
	updateBuilder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: in.ID})

	if in.Username != nil {
		updateBuilder = updateBuilder.Set(usernameColumn, in.Username)
	}

	if in.Password != nil {
		updateBuilder = updateBuilder.Set(passwordColumn, in.Password)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRow: query,
	}

	logger.Debug("Executing query",
		zap.String("sql", q.Name),
		zap.String("query", prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)),
	)
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
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

	logger.Debug("Executing query",
		zap.String("sql", q.Name),
		zap.String("query", prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)),
	)

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return err
	}

	return nil
}
