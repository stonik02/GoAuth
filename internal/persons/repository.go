package person

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"

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

// HashPassword hashes the user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with an unhashed password
func CheckPasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// ScanPersonFromPgxRowsInStruct scans the persons retrieved
// from the database in type pgr.Rows into struct Person.
func (r *repository) ScanPersonFromPgxRowsInStruct(rows pgx.Rows) ([]ResponseUserDto, error) {
	var persons []ResponseUserDto

	for rows.Next() {
		var prs ResponseUserDto
		err := rows.Scan(&prs.Id, &prs.Name, &prs.Email)
		if err != nil {
			return nil, r.pgClient.LoggingSQLPgqError(err)
		}
		persons = append(persons, prs)
	}
	return persons, nil
}

// CheckUserExist checks that the user id is not a null string
func CheckUserExist(person ResponseUserDto) error {
	if person.Id == "" {
		newErr := fmt.Errorf("user does not exist")
		return newErr
	}
	return nil
}

// CheckingFieldsPersonHaveBeenChanged compares the fields of the user
// in the database with those of the user sent by the user.
// If a field in personDto has been changed,
// it is changed in personInDatabase as well.
func CheckingFieldsPersonHaveBeenChanged(personDto *Person, personInDatabase ResponseUserDto) ResponseUserDto {
	if personDto.Name != "" {
		personInDatabase.Name = personDto.Name
	}
	if personDto.Email != "" {
		personInDatabase.Email = personDto.Email
	}
	return personInDatabase
}

// Create adds a new user to the database.
// Params:
// ctx - context.Context,
// Person :
// id - string(uuid)
// Name - string
// Email - string (unique)
// Password - string (will be hash)
func (r *repository) Create(ctx context.Context, person *Person) error {
	// Hash password
	HashPassword, err := HashPassword(person.Password)
	if err != nil {
		r.logger.Errorf("CreatePerson error: %s", err)
		return err
	}
	person.Password = HashPassword

	// Create Person in DB
	err = r.pgClient.CreatePersonInDB(ctx, person)
	if err != nil {
		r.logger.Errorf("SQL error ((( error = %s", err)
		return err
	}
	return nil
}

// Delete deletes a user from the database by id.
func (r *repository) Delete(ctx context.Context, id string) {
	// Send a query to the database
	r.pgClient.DeletePersonFromDB(ctx, id)
}

// FindAll retrieves all users.
func (r *repository) FindAll(ctx context.Context) (p []ResponseUserDto, err error) {
	// Send a query to the database to retrieve all users
	rows, err := r.pgClient.FindAllPersonFromDB(ctx)
	if err != nil {
		return nil, err
	}

	// Scan the users into the structure and return
	return r.ScanPersonFromPgxRowsInStruct(rows)
}

// FindOne returns a specific person by id.
func (r *repository) FindOne(ctx context.Context, id string) (p ResponseUserDto, err error) {
	person := r.pgClient.FindOne(ctx, id)
	// Check if such a person exists
	if err = CheckUserExist(person); err != nil {
		return ResponseUserDto{}, err
	}

	return person, nil
}

// Update updates the user's data in the database.
// The method works simultaneously as Patch and Put.
// You can change name and email.
func (r *repository) Update(ctx context.Context, person *Person) error {

	personInDb, err := r.FindOne(ctx, person.Id)
	if err != nil {
		r.logger.Errorf("Update person error: %s", err)
		return err
	}

	personInDb = CheckingFieldsPersonHaveBeenChanged(person, personInDb)

	return r.pgClient.UpdatePerson(ctx, personInDb)
}

// FindByEmail returns a specific person by email.
func (r *repository) FindByEmail(ctx context.Context, email string) (ResponseUserDto, error) {
	person := r.pgClient.FindPersonByEmail(ctx, email)
	if err := CheckUserExist(person); err != nil {
		return ResponseUserDto{}, err
	}
	return person, nil
}

// AuthPerson сomparing the data entered by the user
// with the data in the database to authorize the user.
// If the data is valid, returns true, otherwise returns false.
func (r *repository) AuthPerson(ctx context.Context, dto AuthDto) (ResponseUserAuthDto, error) {
	userData := r.pgClient.GetPersonDataForAuth(ctx, dto)
	if CheckPasswordHash(dto.Password, userData.Hash_Password) {
		return userData, nil
	}
	return ResponseUserAuthDto{}, fmt.Errorf("Wrong data")
}
