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

// ScanRoleIdAndRoleNameToRoleStruct takes two string arrays as input:
// role_id []string
// role_name []string
// Based on them it creates the []Role structure and returns it.
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

// GetUserWithRoles gets the user's id, email, name and a list of their roles by uuid.
func (r *repository) GetUserWithRoles(ctx context.Context, userId string) (AllUserRolesDto, error) {
	row := r.pgClient.SendsQueryToGetUserWithRoles(ctx, userId)

	return r.ScanPersonDataWitchRolesInAllUserRolesDto(row)
}

// AssignRole adds a user by uuid to a role, also by uuid.
// AssignRoleDto:
// UserId string `json: "userId"`
// RoleId string `json: "roleId"`
func (r *repository) AssignRole(ctx context.Context, dto AssignRoleDto) error {
	err := r.pgClient.SendsQueryToAssignRole(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

// TakeRole removes the role from the user by uuid, also by uuid.
// TakeRoleDto:
// UserId string `json: "userId"`
// RoleId string `json: "roleId"`
func (r *repository) TakeRole(ctx context.Context, dto TakeRoleDto) error {
	err := r.pgClient.SendsQueryToTakeRole(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

// ScanRoleNameFromPgxRowsToStringArray parses RoleName from pgx.Rows to []string
func (r *repository) ScanRoleNameFromPgxRowsToStringArray(rows pgx.Rows) ([]string, error) {
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
func (r *repository) GetUserRoleNames(ctx context.Context, userId string) ([]string, error) {
	rows, err := r.pgClient.SendsQueryToGetUserRoleNames(ctx, userId)
	if err != nil {
		return nil, err
	}

	return r.ScanRoleNameFromPgxRowsToStringArray(rows)
}
