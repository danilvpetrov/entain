package sports

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"syreclabs.com/go/faker"
)

// testDataEvent represents a sports event from testdata.
type testDataEvent struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Competition string `json:"competition"`
}

// SeedTestData seeds the database with test data. It returns the number of
// seeded events.
//
// This function is intended to be used in tests only. Please avoid using it in
// a production setup.
func SeedTestData(ctx context.Context, db *sql.DB) (int, error) {
	raw, err := os.ReadFile("testdata/testdata.json")
	if err != nil {
		return 0, err
	}

	var events []testDataEvent
	if err := json.Unmarshal(raw, &events); err != nil {
		return 0, err
	}

	for i, ev := range events {
		if _, err := db.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO events (
				id,
				name,
				category,
				competition,
				visible,
				advertised_start_time
			)
			VALUES (?, ?, ?, ?, ?, ?)`,
			i+1,
			ev.Name,
			ev.Category,
			ev.Competition,
			faker.Number().Between(0, 1),
			faker.Time().Between(
				time.Now().AddDate(0, 0, -1),
				time.Now().AddDate(0, 0, 2),
			).Format(time.RFC3339),
		); err != nil {
			return 0, err
		}
	}

	return len(events), nil
}
