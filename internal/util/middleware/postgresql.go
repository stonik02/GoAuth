package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/stonik02/proxy_service/pkg/db/postgresql"
	"github.com/stonik02/proxy_service/pkg/logging"
)

type PgSQLInterface interface {
	SendsQueryToGetUserRoleNames(ctx context.Context, userId string) (pgx.Rows, error)
}

type PgSQLClient struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewPgClient(client postgresql.Client, logger *logging.Logger) PgSQLInterface {
	return &PgSQLClient{
		client: client,
		logger: logger,
	}
}

const (
	queryGetUserRoleNames = `
	SELECT r.role_name FROM roles r
	INNER JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1;
	   `
)

func (pg *PgSQLClient) LoggingSQLPgqError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		pg.logger.Error(newErr)
		return newErr
	}
	return err
}

// GetUserRoleNames  sends a query to the database to get all role names for user by uuid.
func (pg *PgSQLClient) SendsQueryToGetUserRoleNames(ctx context.Context, userId string) (pgx.Rows, error) {
	pg.logger.Tracef("Get query: %s", strings.ReplaceAll(queryGetUserRoleNames, "\n\t", ""))

	rows, err := pg.client.Query(ctx, queryGetUserRoleNames, userId)
	if err != nil {
		return nil, pg.LoggingSQLPgqError(err)
	}
	return rows, nil
}
