package access

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/sborsh1kmusora/auth/internal/model"
	"github.com/sborsh1kmusora/platform_common/pkg/db"
)

const (
	tableName = "accesses"

	idColumn              = "id"
	endpointAddressColumn = "endpoint_address"
	roleColumn            = "role"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type Repository interface {
	GetList(ctx context.Context) ([]*model.AccessInfo, error)
}

type repo struct {
	client db.Client
}

func NewRepository(client db.Client) Repository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetList(ctx context.Context) ([]*model.AccessInfo, error) {
	builder := sq.Select(idColumn, endpointAddressColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "accesses.GetList",
		QueryRow: query,
	}

	var accessInfo []*model.AccessInfo
	err = r.client.DB().ScanAllContext(ctx, &accessInfo, q, args...)
	if err != nil {
		return nil, err
	}

	return accessInfo, nil
}
