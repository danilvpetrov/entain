package sports

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

// ApplySchema applies Racing API database schema to a database.
func ApplySchema(
	ctx context.Context,
	db *sql.DB,
) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}
