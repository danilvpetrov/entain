package racing_test

import (
	"slices"
	"testing"

	apiracing "github.com/danilvpetrov/entain/api/racing"
	. "github.com/danilvpetrov/entain/racing"
)

func TestListRaces(t *testing.T) {
	s := &Service{
		DB: setupDatabase(t),
	}
	client := setupServer(t, s)

	cases := []struct {
		name      string
		req       *apiracing.ListRacesRequest
		assertion func(
			t *testing.T,
			resp *apiracing.ListRacesResponse,
			err error,
		)
	}{
		{
			name: "no filter",
			req:  &apiracing.ListRacesRequest{},
			assertion: func(t *testing.T, resp *apiracing.ListRacesResponse, err error) {
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
			req: &apiracing.ListRacesRequest{
				MeetingId: []int64{1, 2, 3},
			},
			assertion: func(t *testing.T, resp *apiracing.ListRacesResponse, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if len(resp.GetRaces()) == 0 {
					t.Fatal("expected at least one race")
				}

				for _, race := range resp.GetRaces() {
					if !slices.Contains([]int64{1, 2, 3}, race.GetMeetingId()) {
						t.Errorf("unexpected meeting ID %d for race %+v", race.GetMeetingId(), race)
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
