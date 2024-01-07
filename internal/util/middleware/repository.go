package utils

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/stonik02/proxy_service/pkg/logging"
)

type Repository struct {
	logger   *logging.Logger
	pgClient PgSQLInterface
}

func NewRepository(logger *logging.Logger, pgClient PgSQLInterface) *Repository {
	return &Repository{
		logger:   logger,
		pgClient: pgClient,
	}
}

// ScanRoleNameFromPgxRowsToStringArray parses RoleName from pgx.Rows to []string
func (r *Repository) ScanRoleNameFromPgxRowsToStringArray(rows pgx.Rows) ([]string, error) {
	var roles []string
	for rows.Next() {
		var row string
		err := rows.Scan(&row)
		if err != nil {
			r.logger.Errorf("Scan error: %s", err)
			return nil, err
		}
		roles = append(roles, row)
	}
	return roles, nil
}

// GetUserRoleNames retrieves all role_names of a person by uuid.
func (r *Repository) GetUserRoleNames(ctx context.Context, userId string) ([]string, error) {
	rows, err := r.pgClient.SendsQueryToGetUserRoleNames(ctx, userId)
	if err != nil {
		return nil, err
	}

	return r.ScanRoleNameFromPgxRowsToStringArray(rows)
}
