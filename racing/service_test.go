package racing_test

import (
	"slices"
	"testing"
	"time"

	racingapi "github.com/danilvpetrov/entain/api/racing"
	. "github.com/danilvpetrov/entain/racing"
)

func TestListRaces(t *testing.T) { //nolint:gocognit // Explicit test cases.
	s := &Service{
		DB: setupDatabase(t),
	}
	client := setupServer(t, s)

	cases := []struct {
		assertion func(
			t *testing.T,
			resp *racingapi.ListRacesResponse,
			err error,
		)
		req  *racingapi.ListRacesRequest
		name string
	}{
		{
			name: "no filter",
			req:  &racingapi.ListRacesRequest{},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(resp.GetRaces()) != NumberOfSeededRaces {
					t.Fatalf(
						"expected %d races, got %d",
						NumberOfSeededRaces,
						len(resp.GetRaces()),
					)
				}
			},
		},
		{
			name: "filtered by meeting IDs",
			req: &racingapi.ListRacesRequest{
				MeetingId: []int64{1, 2, 3},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(resp.GetRaces()) == 0 {
					t.Fatal("expected at least one race")
				}

				for _, race := range resp.GetRaces() {
					if !slices.Contains([]int64{1, 2, 3}, race.GetMeetingId()) {
						t.Errorf(
							"unexpected meeting ID %d for race %+v",
							race.GetMeetingId(),
							race,
						)
					}
				}
			},
		},
		{
			name: "filtered by visible only",
			req: &racingapi.ListRacesRequest{
				VisibleOnly: true,
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				for _, race := range resp.GetRaces() {
					if !race.GetVisible() {
						t.Errorf("expected race %+v to be visible", race)
					}
				}
			},
		},
		{
			name: "ordered by an advertised start time ascending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_ADVERTISED_START_TIME_ASC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastAdvertisedTime time.Time
				for _, race := range resp.GetRaces() {
					actual := race.GetAdvertisedStartTime().AsTime()
					if actual.Before(lastAdvertisedTime) {
						t.Fatalf(
							"expected advertised start time to be in ascending order, "+
								"got %v before %v, race %+v",
							actual,
							lastAdvertisedTime,
							race,
						)
					}
					lastAdvertisedTime = actual
				}
			},
		},
		{
			name: "ordered by an advertised start time descending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_ADVERTISED_START_TIME_DESC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				// Initialise to a time in the far future.
				lastAdvertisedTime := time.Now().AddDate(1000, 0, 0)
				for _, race := range resp.GetRaces() {
					actual := race.GetAdvertisedStartTime().AsTime()
					if actual.After(lastAdvertisedTime) {
						t.Fatalf(
							"expected advertised start time to be in descending order, "+
								"got %v after %v, race %+v",
							actual,
							lastAdvertisedTime,
							race,
						)
					}
					lastAdvertisedTime = actual
				}
			},
		},
		{
			name: "ordered by meeting ID ascending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_MEETING_ID_ASC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastMeetingID int64
				for _, race := range resp.GetRaces() {
					if race.GetMeetingId() < lastMeetingID {
						t.Errorf(
							"expected meeting ID to be in ascending order, got %d before %d, race %+v",
							race.GetMeetingId(),
							lastMeetingID,
							race,
						)
					}
					lastMeetingID = race.GetMeetingId()
				}
			},
		},
		{
			name: "ordered by meeting ID descending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_MEETING_ID_DESC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastMeetingID int64 = 1_000_000_000
				for _, race := range resp.GetRaces() {
					if race.GetMeetingId() > lastMeetingID {
						t.Errorf("expected meeting ID to be in descending "+
							"order, got %d after %d, race %+v",
							race.GetMeetingId(), lastMeetingID, race)
					}
					lastMeetingID = race.GetMeetingId()
				}
			},
		},
		{
			name: "oder by name ascending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_NAME_ASC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var name string
				for _, race := range resp.GetRaces() {
					actual := race.GetName()
					if actual < name {
						t.Errorf(
							"expected name to be in ascending order, got %q before %q, race %+v",
							actual,
							name,
							race,
						)
					}
					name = actual
				}
			},
		},
		{
			name: "oder by name descending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_NAME_DESC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				name := "ZZZZZZZZZZZZZZZZZZZZ"
				for _, race := range resp.GetRaces() {
					actual := race.GetName()
					if actual > name {
						t.Errorf(
							"expected name to be in descending order, got %q after %q, race %+v",
							actual,
							name,
							race,
						)
					}
					name = actual
				}
			},
		},
		{
			name: "ordered by number ascending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_NUMBER_ASC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastNumber int64
				for _, race := range resp.GetRaces() {
					actual := race.GetNumber()
					if actual < lastNumber {
						t.Errorf(
							"expected number to be in ascending order, got %d before %d, race %+v",
							actual,
							lastNumber,
							race,
						)
					}
					lastNumber = actual
				}
			},
		},
		{
			name: "ordered by number descending",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_NUMBER_DESC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var lastNumber int64 = 1_000_000
				for _, race := range resp.GetRaces() {
					actual := race.GetNumber()
					if actual > lastNumber {
						t.Errorf(
							"expected number to be in descending order, got %d after %d, race %+v",
							actual,
							lastNumber,
							race,
						)
					}
					lastNumber = actual
				}
			},
		},
		{
			name: "ordered by multiple fields",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_MEETING_ID_ASC,
					racingapi.ListRacesRequest_NUMBER_DESC,
				},
			},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				var (
					lastMeetingID int64
					lastNumber    int64 = 1_000_000
				)
				for _, race := range resp.GetRaces() {
					actualMeetingID := race.GetMeetingId()
					actualNumber := race.GetNumber()
					if actualMeetingID < lastMeetingID {
						t.Errorf(
							"expected meeting ID to be in ascending order, got %d before %d, race %+v",
							actualMeetingID,
							lastMeetingID,
							race,
						)
					}

					if actualMeetingID != lastMeetingID {
						// Reset the last number if the meeting ID changed.
						lastNumber = 1_000_000
					}

					if actualNumber > lastNumber {
						t.Errorf(
							"expected number to be in descending order, got %d after %d, race %+v",
							actualNumber,
							lastNumber,
							race,
						)
					}
					lastMeetingID = actualMeetingID
					lastNumber = actualNumber
				}
			},
		},
		{
			name: "conflicted orderings",
			req: &racingapi.ListRacesRequest{
				OrderBy: []racingapi.ListRacesRequest_OrderBy{
					racingapi.ListRacesRequest_MEETING_ID_ASC,
					racingapi.ListRacesRequest_MEETING_ID_DESC,
				},
			},
			assertion: func(t *testing.T, _ *racingapi.ListRacesResponse, err error) {
				if err == nil {
					t.Fatalf("expected error, got %v", err)
				}
			},
		},
		{
			name: "status field is computed correctly",
			req:  &racingapi.ListRacesRequest{},
			assertion: func(t *testing.T, resp *racingapi.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				for _, race := range resp.GetRaces() {
					expected := racingapi.Race_OPEN
					now := time.Now()
					if race.GetAdvertisedStartTime().AsTime().Before(now) {
						expected = racingapi.Race_CLOSED
					}

					if race.GetStatus() != expected {
						t.Errorf(
							"expected status %v, got %v",
							expected,
							race.GetStatus(),
						)
					}
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := client.ListRaces(t.Context(), c.req)
			c.assertion(t, resp, err)
		})
	}
}

func TestGetRace(t *testing.T) {
	s := &Service{
		DB: setupDatabase(t),
	}
	client := setupServer(t, s)

	cases := []struct {
		assertion func(
			t *testing.T,
			race *racingapi.Race,
			err error,
		)
		req  *racingapi.GetRaceRequest
		name string
	}{
		{
			name: "gets race by ID",
			req: &racingapi.GetRaceRequest{
				RaceId: 1,
			},
			assertion: func(t *testing.T, race *racingapi.Race, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if race.GetId() != 1 {
					t.Fatalf("expected race ID to be 1, got %d", race.GetId())
				}
			},
		},
		{
			name: "non-existing race",
			req: &racingapi.GetRaceRequest{
				RaceId: 999,
			},
			assertion: func(t *testing.T, _ *racingapi.Race, err error) {
				if err == nil {
					t.Fatalf("expected error, got %v", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := client.GetRace(t.Context(), c.req)
			c.assertion(t, resp, err)
		})
	}
}
