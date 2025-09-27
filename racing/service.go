package racing

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	racingapi "github.com/danilvpetrov/entain/api/racing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service handles all requests related to racing. It implements
// apiracing.RacingServer interface.
type Service struct {
	DB *sql.DB
}

// Make sure Service implements the apiracing.RacingServer interface.
var _ racingapi.RacingServer = (*Service)(nil)

// ListRaces returns a list of all races. It can be filtered by meeting IDs.
func (s *Service) ListRaces(
	ctx context.Context,
	req *racingapi.ListRacesRequest,
) (*racingapi.ListRacesResponse, error) {
	filterQuery, args := parseFilter(req)

	rows, err := s.DB.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT
			id,
			meeting_id,
			name,
			number,
			visible,
			advertised_start_time
		 FROM races
		 WHERE id <> 0 %s`,
			filterQuery),
		args...,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Error("failed closing rows", slog.Any("error", err))
		}
	}()

	var races []*racingapi.Race
	for rows.Next() {
		var (
			race                racingapi.Race
			advertisedStartTime time.Time
		)
		if err := rows.Scan(
			&race.Id,
			&race.MeetingId,
			&race.Name,
			&race.Number,
			&race.Visible,
			&advertisedStartTime,
		); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		race.AdvertisedStartTime = timestamppb.New(advertisedStartTime)
		races = append(races, &race)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &racingapi.ListRacesResponse{
		Races: races,
	}, nil
}

// parseFilter builds SQL filter query and its arguments from the provided
// filter object.
func parseFilter(
	req *racingapi.ListRacesRequest,
) (filter string, args []any) {
	var w strings.Builder

	if len(req.GetMeetingId()) > 0 {
		_, _ = w.WriteString("AND meeting_id IN (")

		for i, id := range req.GetMeetingId() {
			if i == len(req.GetMeetingId())-1 {
				_, _ = w.WriteString("?)")
				args = append(args, id)
				break
			}

			_, _ = w.WriteString("?,")
			args = append(args, id)
		}
	}

	return w.String(), args
}
