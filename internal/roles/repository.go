package roles

import (
	"context"
	// "encoding/json"

	"github.com/jackc/pgx/v4"

	"github.com/stonik02/proxy_service/pkg/logging"
)

type repository struct {
	logger   *logging.Logger
	pgClient PgSQLInterface
}

func NewRepository(logger *logging.Logger, pgClient PgSQLInterface) Repository {
	return &repository{
		logger:   logger,
		pgClient: pgClient,
	}
}

// GetAllRoles implements Repository.
func (r *repository) GetAllRoles(ctx context.Context) ([]Role, error) {
	rows, err := r.pgClient.SendsQueryToGetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := r.ScanRolesFromPgxRowsInStruct(rows)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// ScanRolesFromPgxRowsInStruct scans the persons retrieved
// from the database in type pgr.Rows into struct Person.
func (r *repository) ScanRolesFromPgxRowsInStruct(rows pgx.Rows) ([]Role, error) {
	var roles []Role

	for rows.Next() {
		var role Role
		err := rows.Scan(&role.Id, &role.RoleName)
		if err != nil {
			return nil, r.pgClient.LoggingSQLPgqError(err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// ScanPersonDataWitchRolesInAllUserRolesDto input a row with the following fields:
// person_id string,
// email string,
// name string
// role_id []string
// role_name []string
// Based on this data, it creates the AllUserRolesDto structure and returns it.
func (r *repository) ScanPersonDataWitchRolesInAllUserRolesDto(row pgx.Row) (AllUserRolesDto, error) {
	var userRoles AllUserRolesDto
	var role_id []string
	var role_name []string
	err := row.Scan(&userRoles.UserId, &userRoles.Name, &userRoles.Email, &role_id, &role_name)
	if err != nil {
		r.logger.Error(err)
		return AllUserRolesDto{}, err
	}
	roles := ScanRoleIdAndRoleNameToRoleStruct(role_id, role_name)
	userRoles.Roles = roles
	return userRoles, nil
}

// ScanRoleIdAndRoleNameToRoleStruct takes two string arrays as input.
func ScanRoleIdAndRoleNameToRoleStruct(role_id []string, role_name []string) []Role {
	var roles []Role
	for i := 0; i < len(role_id); i++ {
		var role Role
		role.Id = role_id[i]
		role.RoleName = role_name[i]
		roles = append(roles, role)
	}
	return roles
}

// GetUserWithRoles gets the user's id, email, name and a list of their roles by uuid.
func (r *repository) GetUserWithRoles(ctx context.Context, userId string) (AllUserRolesDto, error) {
	row := r.pgClient.SendsQueryToGetUserWithRoles(ctx, userId)

	return r.ScanPersonDataWitchRolesInAllUserRolesDto(row)
}

// AssignRole adds a user by uuid to a role, also by uuid.
func (r *repository) AssignRole(ctx context.Context, dto AssignRoleDto) error {
	err := r.pgClient.SendsQueryToAssignRole(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

// TakeRole removes the role from the user by uuid, also by uuid.
func (r *repository) TakeRole(ctx context.Context, dto TakeRoleDto) error {
	err := r.pgClient.SendsQueryToTakeRole(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}
