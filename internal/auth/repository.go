package auth

import (
	"context"
	"fmt"

	person "github.com/stonik02/proxy_service/internal/persons"
	"github.com/stonik02/proxy_service/internal/token"
	"github.com/stonik02/proxy_service/pkg/db/postgresql"
	"github.com/stonik02/proxy_service/pkg/logging"
)

type repository struct {
	client           postgresql.Client
	logger           *logging.Logger
	personRepository person.Repository
	tokenRepository  token.Repository
}

func NewRepository(client postgresql.Client, logger *logging.Logger, personRepository person.Repository, tokenRepository token.Repository) Repository {
	return &repository{
		client:           client,
		logger:           logger,
		personRepository: personRepository,
		tokenRepository:  tokenRepository,
	}
}

func (r *repository) CheckUserExist(ctx context.Context, email string) error {
	_, err := r.personRepository.FindByEmail(ctx, email)
	if err == nil {
		newErr := fmt.Errorf("CheckUserExist: personRepository.FindByEmail error: %s", err)
		r.logger.Error(newErr)
		return newErr
	}
	return nil
}

// Register implements Repository.
func (r *repository) RegisterPerson(ctx context.Context, dto RegisterDto) (*person.Person, error) {
	err := r.CheckUserExist(ctx, dto.Email)
	if err != nil {
		r.logger.Errorf("RegisterPerson: CheckUserExist error: %s", err)
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
		r.logger.Errorf("RegisterPerson: personRepository.Create error: %s", err)
		return nil, err
	}
	return &newPerson, nil
}

func ResponseUserAuthDtoToPersonDataInToken(person person.ResponseUserAuthDto) token.PersonDataInToken {
	return token.PersonDataInToken{
		Id:    person.Id,
		Email: person.Email,
	}
}

func (r *repository) CreateTokens(person person.ResponseUserAuthDto) (AuthResponseDto, error) {
	accessToken, err := r.tokenRepository.CreateJWTAccessToken(ResponseUserAuthDtoToPersonDataInToken(person))
	if err != nil {
		r.logger.Errorf("CreateTokens: create access token error: %s", err)
		return AuthResponseDto{}, err
	}
	refreshToken, err := r.tokenRepository.CreateJWTRefreshToken(ResponseUserAuthDtoToPersonDataInToken(person))
	if err != nil {
		r.logger.Errorf("CreateTokens: create refresh token error: %s", err)
		return AuthResponseDto{}, err
	}
	return AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Auth implements Repository.
func (r *repository) Auth(ctx context.Context, dto person.AuthDto) (AuthResponseDto, error) {
	hasPersonInDb, err := r.personRepository.AuthPerson(ctx, dto)
	if err != nil {
		newErr := fmt.Errorf("Auth error: wrond data")
		r.logger.Errorf("Auth person error: %s", err)
		return AuthResponseDto{}, newErr
	}

	return r.CreateTokens(hasPersonInDb)

}

// Refresh implements Repository.
func (r *repository) Refresh(ctx context.Context, dto RefreshDto) (RefreshResponseDto, error) {
	panic("unimplemented")
}
