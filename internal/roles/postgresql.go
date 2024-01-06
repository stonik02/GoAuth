package roles

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
)

type PgSQLInterface interface {
	LoggingSQLPgqError(err error) error
	SendsQueryToGetAllRoles(ctx context.Context) (pgx.Rows, error)
	SendsQueryToGetUserWithRoles(ctx context.Context, userId string) pgx.Row
	SendsQueryToAssignRole(ctx context.Context, dto AssignRoleDto) error
	SendsQueryToTakeRole(ctx context.Context, dto TakeRoleDto) error
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
	queryGetAllRoles = `SELECT id, role_name FROM roles`
	// Тут добавил p.id, если будет ошибка - убрать
	// 	queryGetUserWithRoles = `
	// 	SELECT p.id, p.name,p.email, r.id, r.role_name AS roles FROM roles r
	// 	LEFT JOIN user_roles ur ON ur.role_id = r.id
	// 	LEFT JOIN person p ON p.id = ur.user_id
	// 	WHERE ur.user_id = $1;
	//    `
	queryGetUserWithRoles = `
	SELECT p.id as person_id, p.name, p.email,
	ARRAY_AGG(r.id) as role_id, ARRAY_AGG(r.role_name) AS roles 
	FROM roles r 
	LEFT JOIN user_roles ur ON ur.role_id = r.id 
	LEFT JOIN person p ON p.id = ur.user_id 
	WHERE ur.user_id = $1
	GROUP BY person_id;`

	queryAssignRole       = `INSERT INTO user_roles(user_id, role_id) VALUES($1, $2);`
	queryTakeRole         = `DELETE FROM user_roles WHERE user_id = $1 and role_id = $2;`
	queryGetUserRoleNames = `
	SELECT r.role_name FROM roles r
	INNER JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1;
	   `
)

// LoggingSQLPgqError logs sql errors of type *pgconn.PgError
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

// AssignRole sends a query to the database to assign new role to person.
func (pg *PgSQLClient) SendsQueryToAssignRole(ctx context.Context, dto AssignRoleDto) error {
	pg.logger.Tracef("Insert query: %s", queryAssignRole)

	_, err := pg.client.Query(ctx, queryAssignRole, dto.UserId, dto.RoleId)
	if err != nil {
		return pg.LoggingSQLPgqError(err)
	}
	return nil
}

// GetAllRoles sends a query to the database to get all roles.
func (pg *PgSQLClient) SendsQueryToGetAllRoles(ctx context.Context) (pgx.Rows, error) {
	pg.logger.Tracef("Get query: %s", queryGetAllRoles)

	rows, err := pg.client.Query(ctx, queryGetAllRoles)
	if err != nil {
		return nil, pg.LoggingSQLPgqError(err)
	}
	return rows, nil
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

// GetUserWithRoles sends a query to the database to get the person's id, email, name and a list of their roles by uuid.
func (pg *PgSQLClient) SendsQueryToGetUserWithRoles(ctx context.Context, userId string) pgx.Row {
	pg.logger.Tracef("Get query: %s", strings.ReplaceAll(queryGetUserWithRoles, "\n\t", ""))

	row := pg.client.QueryRow(ctx, queryGetUserWithRoles, userId)

	return row
}

// AssignRole sends a query to the database to remove the role from the user.
func (pg *PgSQLClient) SendsQueryToTakeRole(ctx context.Context, dto TakeRoleDto) error {
	pg.logger.Tracef("Delete query: %s", queryTakeRole)

	_, err := pg.client.Query(ctx, queryTakeRole, dto.UserId, dto.RoleId)
	if err != nil {
		return pg.LoggingSQLPgqError(err)
	}
	return nil
}
