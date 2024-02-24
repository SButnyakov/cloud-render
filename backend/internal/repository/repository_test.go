package repository_test

import (
	"cloud-render/internal/db/postgres"
	"cloud-render/internal/lib/config"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	auth DB
	api  DB
}

type DB struct {
	db  *sql.DB
	cfg *config.Config
}

func (r *RepositoryTestSuite) SetupSuite() {
	authConfigPath := os.Getenv("AUTH_TEST_CONFIG_PATH")
	apiConfigPath := os.Getenv("API_TEST_CONFIG_PATH")

	var err error

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)

	r.auth.cfg = config.MustLoad(authConfigPath)
	r.api.cfg = config.MustLoad(apiConfigPath)

	r.auth.db, err = postgres.New(r.auth.cfg.DB)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), r.auth.db)

	r.api.db, err = postgres.New(r.api.cfg.DB)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), r.api.db)

	require.NoError(r.T(), postgres.DropMigrations(r.auth.db, r.auth.cfg.DB.MigrationsPath))
	require.NoError(r.T(), postgres.DropMigrations(r.api.db, r.api.cfg.DB.MigrationsPath))

	require.NoError(r.T(), postgres.MigrateNSteps(r.auth.db, r.auth.cfg.DB.MigrationsPath, 1))
	require.NoError(r.T(), postgres.MigrateNSteps(r.api.db, r.api.cfg.DB.MigrationsPath, 1))
}

func (r *RepositoryTestSuite) TearDownTest() {
	require.NoError(r.T(), postgres.MigrateNSteps(r.auth.db, r.auth.cfg.DB.MigrationsPath, -1))
	require.NoError(r.T(), postgres.MigrateNSteps(r.api.db, r.api.cfg.DB.MigrationsPath, -1))
	require.NoError(r.T(), postgres.MigrateNSteps(r.auth.db, r.auth.cfg.DB.MigrationsPath, 1))
	require.NoError(r.T(), postgres.MigrateNSteps(r.api.db, r.api.cfg.DB.MigrationsPath, 1))
}

func (r *RepositoryTestSuite) TearDownSuite() {
	defer r.auth.db.Close()
	defer r.api.db.Close()
	require.NoError(r.T(), postgres.DropMigrations(r.auth.db, r.auth.cfg.DB.MigrationsPath))
	require.NoError(r.T(), postgres.DropMigrations(r.api.db, r.api.cfg.DB.MigrationsPath))
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
