package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/danilvpetrov/entain/sports"
	_ "github.com/mattn/go-sqlite3" // underscore import for the SQLite driver
)

var (
	sportsDBPath        = os.Getenv("SPORTS_DB_PATH")
	defaultSportsDBPath = "artefacts/sports.db"
)

// setupDB initialises the database connection and applies the necessary schema.
func setupDB(ctx context.Context) (*sql.DB, error) {
	if sportsDBPath == "" {
		sportsDBPath = defaultSportsDBPath
		// Make sure the artefacts directory exists if we are using the default
		// path.
		if err := os.MkdirAll("artefacts", os.ModePerm); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", sportsDBPath)
	if err != nil {
		return nil, err
	}

	if err := sports.ApplySchema(ctx, db); err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Warning: this step is done only for demo purposes. In a real-world
	// application, you would not want to seed test data to the production
	// database mixing test and real data. Remove this step in such cases.
	if _, err := sports.SeedTestData(
		ctx,
		db,
		"sports/testdata/testdata.json",
	); err != nil {
		return nil, err
	}

	return db, nil
}
