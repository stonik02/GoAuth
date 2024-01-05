package person

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, person *Person) error
	FindAll(ctx context.Context) (persons []ResponseUserDto, err error)
	FindByEmail(ctx context.Context, email string) (ResponseUserDto, error)
	FindOne(ctx context.Context, id string) (ResponseUserDto, error)
	Update(ctx context.Context, person *Person) error
	Delete(ctx context.Context, id string)
	AuthPerson(ctx context.Context, dto AuthDto) bool
}
