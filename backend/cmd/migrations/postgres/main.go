package main

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/lib/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("operation not specified")
	}

	var envName string
	var migrationsPath string

	switch os.Args[1] {
	case "api":
		envName = "API_CONFIG_PATH"
		migrationsPath = "file://../../../migrations/api/postgres"
	case "auth":
		envName = "AUTH_CONFIG_PATH"
		migrationsPath = "file://../../../migrations/auth/postgres"
	default:
		log.Fatalf("undefined service")
	}

	envPath := os.Getenv(envName)

	cfg := config.MustLoad(envPath)

	pg, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to create database instance. Err: %s", err.Error())
	}
	defer pg.Close()

	switch os.Args[2] {
	case "top":
		err = migrateTop(pg, migrationsPath)
	case "drop":
		err = dropMigrations(pg, migrationsPath)
	default:
		err = migrateNSteps(pg, migrationsPath, os.Args[1])
	}
}

func migrateTop(pg *sql.DB, migrationsPath string) error {
	return postgres.MigrateTop(pg, migrationsPath)
}

func dropMigrations(pg *sql.DB, migrationsPath string) error {
	return postgres.DropMigrations(pg, migrationsPath)
}

func migrateNSteps(pg *sql.DB, migrationsPath string, n string) error {
	steps, err := strconv.Atoi(n)
	if err != nil {
		return fmt.Errorf("wrong type of argument: %w", err)
	}
	if steps == 0 {
		return fmt.Errorf("wrong number of steps. n > 0 to migrate up and n < 0 to migrate down")
	}

	return postgres.MigrateNSteps(pg, migrationsPath, steps)
}
