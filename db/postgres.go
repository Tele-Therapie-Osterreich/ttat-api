package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"

	"github.com/Tele-Therapie-Osterreich/ttat-api/chassis"
)

// PGClient is a wrapper for the user database connection.
type PGClient struct {
	DB *sqlx.DB
}

// NewPGClient creates a new user database connection.
func NewPGClient(ctx context.Context, dbURL string) (*PGClient, error) {
	db, err := chassis.DBConnect(ctx, dbURL, Asset, AssetDir)
	if err != nil {
		return nil, err
	}
	return &PGClient{db}, nil
}
