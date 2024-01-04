package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/logging/db/postgresql"
	"golang.org/x/crypto/bcrypt"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Create implements user.Repository.
func (r *repository) Create(ctx context.Context, person *Person) error {
	HashPassword, err := HashPassword(person.Password)
	if err != nil {
		return err
	}
	query := `insert into public.person (name, email, hash_password) values ($1, $2, $3) returning id;`
	r.logger.Tracef("Get query: %s", query)

	err = r.client.QueryRow(ctx, query, person.Name, person.Email, HashPassword).Scan(&person.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
	}
	return nil
}

// Delete implements user.Repository.
func (r *repository) Delete(ctx context.Context, id string) error {
	query := "DELETE from public.person where id = $1;"
	r.logger.Tracef("Get query: %s", query)

	err := r.client.QueryRow(ctx, query, id)
	if err != nil {

		r.logger.Error("SQL error: %s", err)
		error := fmt.Errorf("SQL error: %s", err)
		return error

	}
	return nil
}

// FindAll implements user.Repository.
func (r *repository) FindAll(ctx context.Context) (p []ResponseUserDto, err error) {
	// query := `insert into public person (name, age) values ($1, $2) returning id`
	q2 := `select id, name, email from public.person;`

	rows, err := r.client.Query(context.TODO(), q2)

	persons := make([]ResponseUserDto, 0)

	for rows.Next() {
		var prs ResponseUserDto

		err = rows.Scan(&prs.Id, &prs.Name, &prs.Email)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.logger.Error(newErr)
				return nil, newErr
			}
		}

		persons = append(persons, prs)
	}
	fmt.Println(persons)

	return persons, nil
}

// FindOne implements user.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (p ResponseUserDto, err error) {
	query := "select id, name, email from public.person where id = $1;"
	r.logger.Tracef("Get query: %s", query)

	var prs ResponseUserDto

	err = r.client.QueryRow(ctx, query, "234234234234").Scan(&prs.Id, &prs.Name, &prs.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return ResponseUserDto{}, newErr
		}
	}

	if prs.Id == "" {
		newErr := fmt.Errorf("User does not exist")
		return ResponseUserDto{}, newErr
	}

	return prs, nil
}

// Update implements user.Repository.
func (r *repository) Update(ctx context.Context, person *Person) error {
	var query_get_person, query_update_person string
	query_get_person = "select id, name, email from public.person where id = $1;"
	query_update_person = "update public.person set (name, email) = ($1, $2) where id = $3;"
	var PersonInDb Person

	r.logger.Tracef("Get query: %s", query_get_person)
	err := r.client.QueryRow(ctx, query_get_person, person.Id).Scan(&PersonInDb.Id, &PersonInDb.Name, &PersonInDb.Email)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	if person.Name != "" {
		PersonInDb.Name = person.Name
	}
	if person.Email != "" {
		PersonInDb.Email = person.Email
	}
	r.logger.Tracef("Update query: %s", query_update_person)
	_, err = r.client.Query(ctx, query_update_person, PersonInDb.Name, PersonInDb.Email, PersonInDb.Id)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (userExist int64, err error) {
	query_check_user_exist := `select email from public.person where email = $1;`
	r.logger.Tracef("Get query: %s", query_check_user_exist)
	userExist = 0

	err = r.client.QueryRow(ctx, query_check_user_exist, email).Scan(&userExist)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return 0, newErr
		}
	}

	return userExist, nil
}
