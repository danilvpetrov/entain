package racing

import (
	"context"
	"database/sql"
	"time"

	"syreclabs.com/go/faker"
)

const (
	// NumberOfSeededRaces defines how many races are seeded in the database.
	NumberOfSeededRaces = 100
)

// SeedTestData seeds the database with test data.
// This function is intended to be used in tests only. Please avoid using it
// in a production setup.
func SeedTestData(ctx context.Context, db *sql.DB) error {
	for i := range NumberOfSeededRaces {
		if _, err := db.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO races(
					id,
					meeting_id,
					name,
					number,
					visible,
					advertised_start_time
				) VALUES (?,?,?,?,?,?)`,
			i+1,
			faker.Number().Between(1, 10),
			faker.Team().Name(),
			faker.Number().Between(1, 12),
			faker.Number().Between(0, 1),
			faker.Time().Between(
				time.Now().AddDate(0, 0, -1),
				time.Now().AddDate(0, 0, 2),
			).Format(time.RFC3339),
		); err != nil {
			return err
		}
	}

	return nil
}
