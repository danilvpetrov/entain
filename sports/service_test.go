package sports_test

import (
	"slices"
	"testing"
	"time"

	sportsapi "github.com/danilvpetrov/entain/api/sports"
	. "github.com/danilvpetrov/entain/sports"
)

const (
	// maxAlphaValue is a string that is greater than any other string when
	// compared lexicographically.
	// It is used in tests to verify descending order.
	// The length of the string is chosen to be sufficiently long to cover
	// typical use cases.
	maxAlphaValue = "zzzzzzzzzzzzzzzzzzzzzz"
)

func TestListRaces(t *testing.T) { //nolint:gocognit // Explicit test cases.
	db, numberOfSeedRecords := setupDatabase(t)
	s := &Service{
		DB: db,
	}
	client := setupServer(t, s)

	cases := []struct {
		assertion func(
			t *testing.T,
			resp *sportsapi.ListEventsResponse,
			err error,
		)
		req  *sportsapi.ListEventsRequest
		name string
	}{
		{
			name: "no filter",
			req:  &sportsapi.ListEventsRequest{},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(resp.GetEvents()) != numberOfSeedRecords {
					t.Fatalf(
						"expected %d events, got %d",
						numberOfSeedRecords,
						len(resp.GetEvents()),
					)
				}
			},
		},
		{
			name: "filtered by categories",
			req: &sportsapi.ListEventsRequest{
				EventCategory: []sportsapi.Event_Category{
					sportsapi.Event_BASKETBALL,
					sportsapi.Event_SOCCER,
					sportsapi.Event_TENNIS,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(resp.GetEvents()) == 0 {
					t.Fatal("expected at least one event")
				}

				for _, event := range resp.GetEvents() {
					if !slices.Contains(
						[]sportsapi.Event_Category{
							sportsapi.Event_BASKETBALL,
							sportsapi.Event_SOCCER,
							sportsapi.Event_TENNIS,
						},
						event.GetCategory(),
					) {
						t.Errorf(
							"unexpected event category %v for event %+v",
							event.GetCategory(),
							event,
						)
					}
				}
			},
		},
		{
			name: "filtered by visible only",
			req: &sportsapi.ListEventsRequest{
				VisibleOnly: true,
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				for _, event := range resp.GetEvents() {
					if !event.GetVisible() {
						t.Errorf("expected event %+v to be visible", event)
					}
				}
			},
		},
		{
			name: "ordered by an advertised start time ascending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_ADVERTISED_START_TIME_ASC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastAdvertisedTime time.Time
				for _, event := range resp.GetEvents() {
					actual := event.GetAdvertisedStartTime().AsTime()
					if actual.Before(lastAdvertisedTime) {
						t.Fatalf(
							"expected advertised start time to be in ascending order, "+
								"got %v before %v, event %+v",
							actual,
							lastAdvertisedTime,
							event,
						)
					}
					lastAdvertisedTime = actual
				}
			},
		},
		{
			name: "ordered by an advertised start time descending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_ADVERTISED_START_TIME_DESC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				// Initialise to a time in the far future.
				lastAdvertisedTime := time.Now().AddDate(1000, 0, 0)
				for _, event := range resp.GetEvents() {
					actual := event.GetAdvertisedStartTime().AsTime()
					if actual.After(lastAdvertisedTime) {
						t.Fatalf(
							"expected advertised start time to be in descending order, "+
								"got %v after %v, event %+v",
							actual,
							lastAdvertisedTime,
							event,
						)
					}
					lastAdvertisedTime = actual
				}
			},
		},
		{
			name: "oder by name ascending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_NAME_ASC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var name string
				for _, event := range resp.GetEvents() {
					actual := event.GetName()
					if actual < name {
						t.Errorf(
							"expected name to be in ascending order, got %q before %q, event %+v",
							actual,
							name,
							event,
						)
					}
					name = actual
				}
			},
		},
		{
			name: "oder by name descending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_NAME_DESC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				name := maxAlphaValue
				for _, event := range resp.GetEvents() {
					actual := event.GetName()
					if actual > name {
						t.Errorf(
							"expected name to be in descending order, got %q after %q, event %+v",
							actual,
							name,
							event,
						)
					}
					name = actual
				}
			},
		},
		{
			name: "oder by competition ascending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_COMPETITION_ASC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var competition string
				for _, event := range resp.GetEvents() {
					actual := event.GetCompetition()
					if actual < competition {
						t.Errorf(
							"expected competition to be in ascending order, got %q before %q, event %+v",
							actual,
							competition,
							event,
						)
					}
					competition = actual
				}
			},
		},
		{
			name: "oder by competition descending",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_COMPETITION_DESC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				competition := maxAlphaValue
				for _, event := range resp.GetEvents() {
					actual := event.GetCompetition()
					if actual > competition {
						t.Errorf(
							"expected competition to be in descending order, got %q after %q, event %+v",
							actual,
							competition,
							event,
						)
					}
					competition = actual
				}
			},
		},
		{
			name: "ordered by multiple fields",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_NAME_ASC,
					sportsapi.ListEventsRequest_COMPETITION_DESC,
				},
			},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				lastName := ""
				lastCompetition := maxAlphaValue
				for _, event := range resp.GetEvents() {
					actualName := event.GetName()
					actualCompetition := event.GetCompetition()
					if actualName < lastName {
						t.Errorf(
							"expected name to be in ascending order, got %q before %q, event %+v",
							actualName,
							lastName,
							event,
						)
					}
					if actualName != lastName {
						// Reset the last competition if the name changed.
						lastCompetition = maxAlphaValue
					}
					if actualCompetition > lastCompetition {
						t.Errorf(
							"expected competition to be in descending order, got %q after %q, event %+v",
							actualCompetition,
							lastCompetition,
							event,
						)
					}
					lastName = actualName
					lastCompetition = actualCompetition
				}
			},
		},
		{
			name: "conflicted orderings",
			req: &sportsapi.ListEventsRequest{
				OrderBy: []sportsapi.ListEventsRequest_OrderBy{
					sportsapi.ListEventsRequest_COMPETITION_ASC,
					sportsapi.ListEventsRequest_COMPETITION_DESC,
				},
			},
			assertion: func(t *testing.T, _ *sportsapi.ListEventsResponse, err error) {
				if err == nil {
					t.Fatalf("expected error, got %v", err)
				}
			},
		},
		{
			name: "status field is computed correctly",
			req:  &sportsapi.ListEventsRequest{},
			assertion: func(t *testing.T, resp *sportsapi.ListEventsResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				for _, event := range resp.GetEvents() {
					expected := sportsapi.Event_OPEN
					now := time.Now()
					if event.GetAdvertisedStartTime().AsTime().Before(now) {
						expected = sportsapi.Event_CLOSED
					}

					if event.GetStatus() != expected {
						t.Errorf(
							"expected status %v, got %v",
							expected,
							event.GetStatus(),
						)
					}
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := client.ListEvents(t.Context(), c.req)
			c.assertion(t, resp, err)
		})
	}
}

func TestGetRace(t *testing.T) {
	db, _ := setupDatabase(t)
	s := &Service{
		DB: db,
	}
	client := setupServer(t, s)

	cases := []struct {
		assertion func(
			t *testing.T,
			race *sportsapi.Event,
			err error,
		)
		req  *sportsapi.GetEventRequest
		name string
	}{
		{
			name: "gets event by ID",
			req: &sportsapi.GetEventRequest{
				EventId: 1,
			},
			assertion: func(t *testing.T, event *sportsapi.Event, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if event.GetId() != 1 {
					t.Fatalf("expected event ID to be 1, got %d", event.GetId())
				}
			},
		},
		{
			name: "non-existing event ID",
			req: &sportsapi.GetEventRequest{
				EventId: 999,
			},
			assertion: func(t *testing.T, _ *sportsapi.Event, err error) {
				if err == nil {
					t.Fatalf("expected error, got %v", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := client.GetEvent(t.Context(), c.req)
			c.assertion(t, resp, err)
		})
	}
}
