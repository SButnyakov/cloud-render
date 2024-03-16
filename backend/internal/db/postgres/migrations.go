package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func prepareDB(pg *sql.DB, migrationsPath string) (*migrate.Migrate, error) {
	driver, err := pgmigrate.WithInstance(pg, &pgmigrate.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create driver. Err: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration: %w", err)
	}
	return m, nil
}

func MigrateNSteps(pg *sql.DB, migrationsPath string, n int) error {
	m, err := prepareDB(pg, migrationsPath)
	if err != nil {
		return err
	}
	return m.Steps(n)
}

func MigrateTop(pg *sql.DB, migrationsPath string) error {
	m, err := prepareDB(pg, migrationsPath)
	if err != nil {
		return err
	}

	for {
		err = m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func DropMigrations(pg *sql.DB, migrationsPath string) error {
	m, err := prepareDB(pg, migrationsPath)
	if err != nil {
		return err
	}
	return m.Drop()
}
