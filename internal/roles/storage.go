package roles

import "context"

type Repository interface {
	GetUserWithRoles(ctx context.Context, userId string) (AllUserRolesDto, error)
	AssignRole(ctx context.Context, dto AssignRoleDto) error
	TakeRole(ctx context.Context, dto TakeRoleDto) error
	GetAllRoles(ctx context.Context) ([]Role, error)
}
