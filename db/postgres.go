package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// PGClient is a wrapper for the user database connection.
type PGClient struct {
	DB *sqlx.DB
}

type assetFunc func(name string) ([]byte, error)
type assetDirFunc func(name string) ([]string, error)

// NewPGClient creates a new user database connection.
func NewPGClient(ctx context.Context, dbURL string) (*PGClient, error) {
	// Connect to database and test connection integrity.
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "opening database")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "pinging database")
	}

	// Limit maximum connections (default is unlimited).
	db.SetMaxOpenConns(10)

	// Run and log database migrations.
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}
	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal().Err(err).Msg("database migrations failed!")
	}
	log.Info().Msgf("applied new database migrations: %d", n)
	migrationRecords, err := migrate.GetMigrationRecords(db.DB, "postgres")
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't read back database migration records")
	}
	if len(migrationRecords) == 0 {
		log.Info().Msg("no database migrations currently applied")
	} else {
		for _, m := range migrationRecords {
			log.Info().
				Str("migration", m.Id).
				Time("applied_at", m.AppliedAt).
				Msg("database migration")
		}
	}

	return &PGClient{db}, nil
}
