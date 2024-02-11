package repository

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/models"
	"database/sql"
	"errors"
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

func (u *UserRepository) CheckCredentials(loginOrEmail, password string) (int64, error) {
	const fn = packagePath + "CheckCredentials"

	stmt, err := u.db.Prepare("SELECT id FROM users WHERE (login=$1 OR email=$1) AND password=$2")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var uid int64

	err = stmt.QueryRow(loginOrEmail, password).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}

		return 0, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return uid, nil
}

func (u *UserRepository) UpdateRefreshToken(uid int64, refreshToken string) error {
	const fn = packagePath + "UpdateRefreshToken"

	stmt, err := u.db.Prepare("UPDATE users SET refresh_token = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(refreshToken, uid)
	if err != nil {
		return fmt.Errorf("%s: exec statement: %w", fn, err)
	}

	return nil
}
