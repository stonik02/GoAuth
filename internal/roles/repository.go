package roles

import (
	"context"
	// "encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

const (
	query_get_all_roles = `
	SELECT id, role_name FROM roles
	`
	query_get_user_with_roles = `
	SELECT p.name,p.email, r.id,  r.role_name AS roles
	FROM roles r
	LEFT JOIN user_roles ur ON ur.role_id = r.id
	LEFT JOIN person p ON p.id = ur.user_id
	WHERE ur.user_id = $1;
   `
	query_assign_role = `
	INSERT INTO user_roles(user_id, role_id) VALUES($1, $2);
	`
	query_take_role = `
	DELETE FROM user_roles WHERE user_id = $1 and role_id = $2;
	`
	query_get_user_role_names = `
	SELECT r.role_name FROM roles r
	INNER JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1;
	   `
)

// GetAllRoles implements Repository.
func (r *repository) GetAllRoles(ctx context.Context) ([]Role, error) {
	r.logger.Tracef("Get query: %s", query_get_all_roles)
	var roles []Role

	rows, err := r.client.Query(ctx, query_get_all_roles)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return roles, newErr
		}
	}
	for rows.Next() {
		var rol Role

		err = rows.Scan(&rol.Id, &rol.RoleName)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return nil, newErr
			}
		}

		roles = append(roles, rol)
	}
	return roles, nil
}

func (r *repository) GetUserWithRoles(ctx context.Context, userId string) (AllUserRolesDto, error) {
	r.logger.Tracef("Get query: %s", query_get_user_with_roles)
	var userRoles AllUserRolesDto

	rows, err := r.client.Query(ctx, query_get_user_with_roles, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return AllUserRolesDto{}, newErr
		}
	}
	flag := 0
	for rows.Next() {
		var name, email string
		var rol Role

		err = rows.Scan(&name, &email, &rol.Id, &rol.RoleName)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return AllUserRolesDto{}, newErr
			}
		}

		if flag == 0 {
			userRoles.Email = email
			userRoles.Name = name
		}

		userRoles.Roles = append(userRoles.Roles, rol)
		flag++
	}
	userRoles.UserId = userId
	return userRoles, nil

}

// AssignRole implements Repository.
func (r *repository) AssignRole(ctx context.Context, dto AssignRoleDto) error {
	r.logger.Tracef("Insert query: %s", query_assign_role)

	_, err := r.client.Query(ctx, query_assign_role, dto.UserId, dto.RoleId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
	}

	return nil
}

// TakeRole implements Repository.
func (r *repository) TakeRole(ctx context.Context, dto TakeRoleDto) error {
	r.logger.Tracef("Delete query: %s", query_take_role)

	_, err := r.client.Query(ctx, query_take_role, dto.UserId, dto.RoleId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
	}

	return nil
}

func (r *repository) GetUserRoleNames(ctx context.Context, userId string) (roles []string, err error) {
	r.logger.Tracef("Get query: %s", query_get_user_role_names)

	rows, err := r.client.Query(ctx, query_get_user_role_names, userId)
	if err != nil {
		r.logger.Errorf("SQL error: %s", err)
		return nil, err
	}
	for rows.Next() {
		var row string
		err = rows.Scan(&row)
		if err != nil {
			r.logger.Errorf("Scan error: %s", err)
			return nil, err
		}
		roles = append(roles, row)
	}

	return roles, nil
}
