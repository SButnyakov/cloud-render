package postgres

import (
	"cloud-render/internal/lib/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	UniqueViolationErrorCode = "23505"
)

func New(cfg config.DB) (*sql.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}
