package repository

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user models.User) error {
	const fn = packagePath + "user.CreateUser"

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

func (u *UserRepository) GetOneUser(uid int64) (*models.User, error) {
	const fn = packagePath + "user.GetOneUser"

	var resUser models.User

	stmt, err := u.db.Prepare("SELECT id, login, email, password FROM users WHERE id=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	err = stmt.QueryRow(uid).Scan(&resUser.Id, &resUser.Login, &resUser.Email, &resUser.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return &resUser, nil
}

func (u *UserRepository) UpdateUser(user models.User) error {
	const fn = packagePath + "user.Update"

	stmt, err := u.db.Prepare("UPDATE users SET login=$2, email=$3, password=$4 WHERE id=$1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	res, err := stmt.Exec(user.Id, user.Login, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) GetHashedPassword(loginOrEmail, password string) ([]models.User, error) {
	const fn = packagePath + "user.GetHashedPassword"

	stmt, err := u.db.Prepare("SELECT id, password FROM users WHERE login=$1 OR email=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query(loginOrEmail)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	users := make([]models.User, 0)

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Id, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepository) UpdateRefreshToken(uid int64, refreshToken string) error {
	const fn = packagePath + "user.UpdateRefreshToken"

	stmt, err := u.db.Prepare("UPDATE users SET refresh_token = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	res, err := stmt.Exec(refreshToken, uid)
	if err != nil {
		return fmt.Errorf("%s: exec statement: %w", fn, err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (u *UserRepository) GetRefreshToken(uid int64) (string, error) {
	const fn = packagePath + "user.GetRefreshToken"

	stmt, err := u.db.Prepare("SELECT refresh_token FROM users WHERE id = $1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var token string

	err = stmt.QueryRow(uid).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return token, nil
}
