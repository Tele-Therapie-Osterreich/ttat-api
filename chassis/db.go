package chassis

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
)

type assetFunc func(name string) ([]byte, error)
type assetDirFunc func(name string) ([]string, error)

// DBConnect creates a new database connection.
func DBConnect(ctx context.Context, dbURL string,
	asset assetFunc, assetDir assetDirFunc) (*sqlx.DB, error) {
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
		Asset:    asset,
		AssetDir: assetDir,
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

	return db, nil
}

// Paginate does the API-wide processing of pagination controls:
// maximum limit is 100, default limit is 30, default offset is zero.
func Paginate(inLimit, inOffset *uint) string {
	limit := uint(30)
	if inLimit != nil {
		limit = *inLimit
	}
	if limit > 100 {
		limit = uint(100)
	}
	offset := uint(0)
	if inOffset != nil {
		offset = *inOffset
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}
