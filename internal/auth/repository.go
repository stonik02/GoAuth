package auth

import (
	"context"
	"fmt"

	person "github.com/stonik02/proxy_service/internal/persons"
	"github.com/stonik02/proxy_service/pkg/db/postgresql"
	"github.com/stonik02/proxy_service/pkg/logging"

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

func (r *repository) CheckUserExist(ctx context.Context, email string) error {
	_, err := r.personRepository.FindByEmail(ctx, email)
	if err == nil {
		newErr := fmt.Errorf("Registration error: user with email = %s is exist", email)
		r.logger.Error(newErr)
		return newErr
	}
	return nil
}

// Register implements Repository.
func (r *repository) RegisterPerson(ctx context.Context, dto RegisterDto) (*person.Person, error) {
	err := r.CheckUserExist(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	r.logger.Tracef("dto create user = %s", dto)
	newPerson := person.Person{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	err = r.personRepository.Create(ctx, &newPerson)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return &newPerson, nil
}

// Auth implements Repository.
func (r *repository) Auth(ctx context.Context, dto person.AuthDto) (AuthResponseDto, error) {
	var response AuthResponseDto
	hasPersonInDb := r.personRepository.AuthPerson(ctx, dto)
	if !hasPersonInDb {
		newErr := fmt.Errorf("wrong data")
		return AuthResponseDto{}, newErr
	}

	//TODO: Тут получаем jwt, все дела
	//
	//
	//

	response = AuthResponseDto{
		accessToken:  "Типо токен",
		refreshToken: "Типо токен",
	}

	return response, nil

}

// Refresh implements Repository.
func (r *repository) Refresh(ctx context.Context, dto RefreshDto) (RefreshResponseDto, error) {
	panic("unimplemented")
}
