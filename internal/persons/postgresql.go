package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
)

type PgSQLInterface interface {
	LoggingSQLPgqError(err error) error
	CreatePersonInDB(ctx context.Context, person *Person) error
	DeletePersonFromDB(ctx context.Context, id string)
	FindAllPersonFromDB(ctx context.Context) (rows pgx.Rows, err error)
	FindOne(ctx context.Context, id string) (p ResponseUserDto)
	UpdatePerson(ctx context.Context, person ResponseUserDto) error
	FindPersonByEmail(ctx context.Context, email string) ResponseUserDto
	GetPersonDataForAuth(ctx context.Context, dto AuthDto) AuthDto
}

type PgSQLClient struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewPgClient(client postgresql.Client, logger *logging.Logger) PgSQLInterface {
	return &PgSQLClient{
		client: client,
		logger: logger,
	}
}

const (
	queryCreatePerson                     = `INSERT INTO public.person (name, email, hash_password) VALUES ($1, $2, $3) RETURNING id;`
	queryDeletePerson                     = `DELETE FROM public.person WHERE id = $1;`
	queryFindAllPerson                    = `SELECT id, name, email FROM public.person;`
	queryFindPersonById                   = `SELECT id, name, email FROM public.person WHERE id = $1;`
	queryUpdatePerson                     = `UPDATE public.person SET (name, email) = ($1, $2) WHERE id = $3;`
	queryFindPersonByEmail                = `SELECT id, name, email FROM public.person WHERE email = $1;`
	queryGetEmailAndPAsswordPersonForAuth = `SELECT email, hash_password FROM public.person WHERE email = $1`
)

// LoggingSQLPgqError logs sql errors of type *pgconn.PgError
func (pg *PgSQLClient) LoggingSQLPgqError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		pg.logger.Error(newErr)
		return newErr
	}
	return err
}

// CreatePerson sends a query to the database to add a new person.
func (pg *PgSQLClient) CreatePersonInDB(ctx context.Context, person *Person) error {
	pg.logger.Tracef("Get query: %s", queryCreatePerson)

	err := pg.client.QueryRow(ctx, queryCreatePerson, person.Name, person.Email, HashPassword).Scan(&person.Id)
	if err != nil {
		return pg.LoggingSQLPgqError(err)
	}
	return nil
}

// DeletePerson sends a query to the database to delete the person.
func (pg *PgSQLClient) DeletePersonFromDB(ctx context.Context, id string) {
	pg.logger.Tracef("Delete query: %s", queryDeletePerson)
	pg.client.QueryRow(ctx, queryDeletePerson, id)
}

// FindAllPersonFromDB sends a query to the database to retrieve all users.
func (pg *PgSQLClient) FindAllPersonFromDB(ctx context.Context) (rows pgx.Rows, err error) {
	pg.logger.Tracef("Get query: %s", queryFindAllPerson)
	rows, err = pg.client.Query(ctx, queryFindAllPerson)
	if err != nil {
		return nil, pg.LoggingSQLPgqError(err)
	}
	return rows, nil
}

// FindOne sends a query to the database to retrieve a specific user by id.
func (pg *PgSQLClient) FindOne(ctx context.Context, id string) ResponseUserDto {
	pg.logger.Tracef("Get query: %s", queryFindPersonById)
	var person ResponseUserDto
	pg.client.QueryRow(ctx, queryFindPersonById, id).Scan(&person.Id, &person.Name, &person.Email)
	return person
}

// UpdatePerson sends a query to the database to update all user data.
func (pg *PgSQLClient) UpdatePerson(ctx context.Context, person ResponseUserDto) error {
	pg.logger.Tracef("Update query: %s", queryUpdatePerson)
	_, err := pg.client.Query(ctx, queryUpdatePerson, person.Name, person.Email, person.Id)
	if err != nil {
		return pg.LoggingSQLPgqError(err)
	}
	return nil
}

// FindByEmail sends a query to the database to retrieve a specific user by email.
func (pg *PgSQLClient) FindPersonByEmail(ctx context.Context, email string) ResponseUserDto {
	pg.logger.Tracef("Get query: %s", queryFindPersonByEmail)
	var person ResponseUserDto

	pg.client.QueryRow(ctx, queryFindPersonByEmail, email).Scan(&person.Id, &person.Name, &person.Email)

	return person
}

// GetPersonDataForAuth sends a query to the database to retrieve a password and email
// from database for authorization
func (pg *PgSQLClient) GetPersonDataForAuth(ctx context.Context, dto AuthDto) AuthDto {
	pg.logger.Tracef("Get query: %s", queryGetEmailAndPAsswordPersonForAuth)
	var userData AuthDto
	pg.client.QueryRow(ctx, queryGetEmailAndPAsswordPersonForAuth, dto.Email).Scan(&userData.Email, &userData.Password)
	return userData
}
