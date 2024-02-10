package repository

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/models"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

const (
	packagePath = "repository.user."
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user models.User) error {
	const fn = packagePath + "Create"

	stmt, err := u.db.Prepare("INSERT INTO users(login, email, password) values($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(user.Login, user.Email, user.Password)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == postgres.UniqueViolationErrorCode {
			return ErrUserExists
		}

		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}
