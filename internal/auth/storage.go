package auth

import (
	"context"

	"github.com/stonik02/proxy_service/internal/persons"
)

type Repository interface {
	RegisterPerson(ctx context.Context, dto RegisterDto) (*person.Person, error)
	Auth(ctx context.Context, dto person.AuthDto) (AuthResponseDto, error)
	Refresh(ctx context.Context, dto RefreshDto) (RefreshResponseDto, error)
}
