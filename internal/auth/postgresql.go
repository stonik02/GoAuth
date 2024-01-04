package auth

import (
	"context"
	"fmt"

	"github.com/stonik02/proxy_service/internal/person"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
)

type repository struct {
	client           postgresql.Client
	personRepository person.Repository
	logger           *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger, personRepository person.Repository) Repository {
	return &repository{
		client:           client,
		personRepository: personRepository,
		logger:           logger,
	}
}

// Register implements Repository.
func (r *repository) Register(ctx context.Context, dto RegisterDto) (*person.Person, error) {
	userExist, err := r.personRepository.FindByEmail(ctx, dto.Email)
	if err != nil {
		r.logger.Errorf("Sql error: %s", err)
		return nil, err
	}
	if userExist != 0 {
		newErr := fmt.Errorf("Registration error: user with email = %s is exist", dto.Email)
		r.logger.Error(newErr)
		return nil, newErr
	}

	// TODO: На каком этапе хешировать пароль???
	newPerson := person.Person{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}

	err = r.personRepository.Create(ctx, &newPerson)

	if err != nil {
		r.logger.Errorf("Sql error: %s", err)
		return nil, err
	}
	return &newPerson, nil
}

// Auth implements Repository.
func (r *repository) Auth(ctx context.Context, dto AuthDto) (AuthResponseDto, error) {
	panic("unimplemented")
}

// Refresh implements Repository.
func (r *repository) Refresh(ctx context.Context, dto RefreshDto) (RefreshResponseDto, error) {
	panic("unimplemented")
}
