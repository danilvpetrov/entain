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

	orderBy, err := parseOrderBy(req)
	if err != nil {
		return nil, err
	}

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
		 	WHERE id <> 0 %s %s`,
			filterQuery,
			orderBy,
		),
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
		_, _ = w.WriteString(" AND meeting_id IN (")

		for i, id := range req.GetMeetingId() {
			_, _ = w.WriteString("?")
			args = append(args, id)

			// If it's the last element, close the parenthesis and break.
			// Otherwise, add a comma.
			if i == len(req.GetMeetingId())-1 {
				_, _ = w.WriteString(")")
				break
			}

			_, _ = w.WriteString(", ")
		}
	}

	if req.GetVisibleOnly() {
		_, _ = w.WriteString(" AND visible = true")
	}

	return w.String(), args
}

// conflictingOrdering maps each ordering option to its conflicting counterpart.
var conflictingOrdering = map[racingapi.ListRacesRequest_OrderBy]racingapi.ListRacesRequest_OrderBy{
	racingapi.ListRacesRequest_ADVERTISED_START_TIME_ASC:  racingapi.ListRacesRequest_ADVERTISED_START_TIME_DESC,
	racingapi.ListRacesRequest_ADVERTISED_START_TIME_DESC: racingapi.ListRacesRequest_ADVERTISED_START_TIME_ASC,
	racingapi.ListRacesRequest_MEETING_ID_ASC:             racingapi.ListRacesRequest_MEETING_ID_DESC,
	racingapi.ListRacesRequest_MEETING_ID_DESC:            racingapi.ListRacesRequest_MEETING_ID_ASC,
	racingapi.ListRacesRequest_NAME_ASC:                   racingapi.ListRacesRequest_NAME_DESC,
	racingapi.ListRacesRequest_NAME_DESC:                  racingapi.ListRacesRequest_NAME_ASC,
	racingapi.ListRacesRequest_NUMBER_ASC:                 racingapi.ListRacesRequest_NUMBER_DESC,
	racingapi.ListRacesRequest_NUMBER_DESC:                racingapi.ListRacesRequest_NUMBER_ASC,
}

func parseOrderBy(req *racingapi.ListRacesRequest) (string, error) {
	var w strings.Builder

	if len(req.GetOrderBy()) == 0 {
		return "", nil
	}

	w.WriteString(" ORDER BY ")
	visited := map[racingapi.ListRacesRequest_OrderBy]bool{}

	for i, order := range req.GetOrderBy() {
		if visited[conflictingOrdering[order]] {
			return "", status.Error(
				codes.InvalidArgument,
				"conflicting order by fields",
			)
		}

		switch order {
		case racingapi.ListRacesRequest_ADVERTISED_START_TIME_ASC:
			_, _ = w.WriteString("advertised_start_time ASC")
		case racingapi.ListRacesRequest_ADVERTISED_START_TIME_DESC:
			_, _ = w.WriteString("advertised_start_time DESC")
		case racingapi.ListRacesRequest_MEETING_ID_ASC:
			_, _ = w.WriteString("meeting_id ASC")
		case racingapi.ListRacesRequest_MEETING_ID_DESC:
			_, _ = w.WriteString("meeting_id DESC")
		case racingapi.ListRacesRequest_NAME_ASC:
			_, _ = w.WriteString("name ASC")
		case racingapi.ListRacesRequest_NAME_DESC:
			_, _ = w.WriteString("name DESC")
		case racingapi.ListRacesRequest_NUMBER_ASC:
			_, _ = w.WriteString("number ASC")
		case racingapi.ListRacesRequest_NUMBER_DESC:
			_, _ = w.WriteString("number DESC")
		}

		visited[order] = true

		// Add a comma if it's not the last element.
		if i != len(req.GetOrderBy())-1 {
			_, _ = w.WriteString(", ")
		}
	}

	return w.String(), nil
}
