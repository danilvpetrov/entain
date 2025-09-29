package main

import (
	"database/sql"

	"github.com/danilvpetrov/entain/sports"
)

// setupService initialises and returns a new instance of the sports service.
func setupService(db *sql.DB) *sports.Service {
	return &sports.Service{DB: db}
}
