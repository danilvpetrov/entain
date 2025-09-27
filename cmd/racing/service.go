package main

import (
	"database/sql"

	"github.com/danilvpetrov/entain/racing"
)

// setupService initialises and returns a new instance of the racing service.
func setupService(db *sql.DB) *racing.Service {
	return &racing.Service{DB: db}
}
