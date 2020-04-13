package db

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
)

var pg *PGClient

// InitTestDB does setup for the test database.
func InitTestDB() {
	var err error
	pgdsn := os.Getenv("TTAT_TEST_DB")

	dbtmp, err := sqlx.Open("postgres", pgdsn)
	if err != nil {
		log.Fatal().Err(err).Msg("opening test database")
	}
	defer dbtmp.Close()
	MultiExec(dbtmp, dropSchema)

	pg, err = NewPGClient(context.Background(), pgdsn)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't connect to test database")
	}
}

// RunWithSchema runs a test within a context with a database
// connection to a database initialised with the application schema.
func RunWithSchema(t *testing.T, test func(pg *PGClient, t *testing.T)) {
	defer func() {
		MultiExec(pg.DB, dropSchema)
	}()

	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}
	_, err := migrate.Exec(pg.DB.DB, "postgres", migrations, migrate.Up)
	assert.Nil(t, err, "database migrations failed!")

	test(pg, t)
}

// A full pre-test schema reset, dropping all tables and creating a
// new empty schema.
var dropSchema = `
SET ROLE postgres;
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
`

// LoadDefaultFixture loads an SQL test fixture.
func LoadDefaultFixture(db *PGClient, t *testing.T) {
	f, err := os.Open("fixture.sql")
	assert.Nil(t, err)
	defer f.Close()
	fixture, err := ioutil.ReadAll(f)
	assert.Nil(t, err)
	tx := pg.DB.MustBegin()
	MultiExec(tx, string(fixture))
	tx.Commit()
}

// Execer is an interface for database command execution.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// MultiExec executes a script of SQL statements from a multiline
// string.
func MultiExec(db Execer, query string) {
	stmts := strings.Split(query, ";\n")
	if len(strings.Trim(stmts[len(stmts)-1], " \n\t\r")) == 0 {
		stmts = stmts[:len(stmts)-1]
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			log.Fatal().Msgf("executing '%s': %s", s, err.Error())
		}
	}
}
