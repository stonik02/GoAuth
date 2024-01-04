package roles

import "context"

type Repository interface {
	GetUserWithRoles(ctx context.Context, userId string) (AllUserRolesDto, error)
	// TODO: Сделать возможность сразу присваивать несколько ролей
	AssignRole(ctx context.Context, dto AssignRoleDto) error
	TakeRole(ctx context.Context, dto TakeRoleDto) error
	GetAllRoles(ctx context.Context) ([]Role, error)
	GetUserRoleNames(ctx context.Context, userId string) (roles []string, err error)
}
