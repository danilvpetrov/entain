package sports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	sportsapi "github.com/danilvpetrov/entain/api/sports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service handles all requests related to sports. It implements
// sportsapi.SportsServer interface.
type Service struct {
	// DB is a database connection pool used to perform queries against the
	// underlying database.
	DB *sql.DB
}

// Make sure Service implements the sportsapi.SportsServer interface.
var _ sportsapi.SportsServer = (*Service)(nil)

// ListEvents returns a list of all sports events.
func (s *Service) ListEvents(
	ctx context.Context,
	req *sportsapi.ListEventsRequest,
) (*sportsapi.ListEventsResponse, error) {
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
				name,
				category,
				competition,
				visible,
				advertised_start_time
			FROM events
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

	var events []*sportsapi.Event
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sportsapi.ListEventsResponse{
		Events: events,
	}, nil
}

// GetEvent returns a specific sports event by its ID.
func (s *Service) GetEvent(
	ctx context.Context,
	req *sportsapi.GetEventRequest,
) (*sportsapi.Event, error) {
	row := s.DB.QueryRowContext(
		ctx,
		`SELECT
			id,
			name,
			category,
			competition,
			visible,
			advertised_start_time
		FROM events
		WHERE id = ?`,
		req.GetEventId(),
	)

	event, err := scanEvent(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return event, nil
}

// scanner is an interface that abstracts sql.Row and sql.Rows types.
type scanner interface {
	Scan(dest ...any) error
}

// scanEvent scans a sports event from the given scanner.
func scanEvent(s scanner) (*sportsapi.Event, error) {
	var (
		event               sportsapi.Event
		category            string
		advertisedStartTime time.Time
	)
	if err := s.Scan(
		&event.Id,
		&event.Name,
		&category,
		&event.Competition,
		&event.Visible,
		&advertisedStartTime,
	); err != nil {
		return nil, err
	}

	event.Category = sportsapi.Event_Category(
		sportsapi.Event_Category_value[category],
	)
	event.AdvertisedStartTime = timestamppb.New(advertisedStartTime)
	event.Status = computeEventStatus(advertisedStartTime)

	return &event, nil
}

// parseFilter builds SQL filter query and its arguments from the provided
// filter object.
func parseFilter(
	req *sportsapi.ListEventsRequest,
) (filter string, args []any) {
	var w strings.Builder

	if len(req.GetCategory()) > 0 {
		_, _ = w.WriteString(" AND category IN (")

		for i, cat := range req.GetCategory() {
			_, _ = w.WriteString("?")
			args = append(
				args,
				sportsapi.Event_Category_name[int32(cat)],
			)

			// If it's the last element, close the parenthesis and break.
			// Otherwise, add a comma.
			if i == len(req.GetCategory())-1 {
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
var conflictingOrdering = map[sportsapi.ListEventsRequest_OrderBy]sportsapi.ListEventsRequest_OrderBy{
	sportsapi.ListEventsRequest_ADVERTISED_START_TIME_ASC:  sportsapi.ListEventsRequest_ADVERTISED_START_TIME_DESC,
	sportsapi.ListEventsRequest_ADVERTISED_START_TIME_DESC: sportsapi.ListEventsRequest_ADVERTISED_START_TIME_ASC,
	sportsapi.ListEventsRequest_NAME_ASC:                   sportsapi.ListEventsRequest_NAME_DESC,
	sportsapi.ListEventsRequest_NAME_DESC:                  sportsapi.ListEventsRequest_NAME_ASC,
	sportsapi.ListEventsRequest_COMPETITION_ASC:            sportsapi.ListEventsRequest_COMPETITION_DESC,
	sportsapi.ListEventsRequest_COMPETITION_DESC:           sportsapi.ListEventsRequest_COMPETITION_ASC,
}

func parseOrderBy(req *sportsapi.ListEventsRequest) (string, error) {
	var w strings.Builder

	if len(req.GetOrderBy()) == 0 {
		return "", nil
	}

	w.WriteString(" ORDER BY ")
	visited := map[sportsapi.ListEventsRequest_OrderBy]bool{}

	for i, order := range req.GetOrderBy() {
		if visited[conflictingOrdering[order]] {
			return "", status.Error(
				codes.InvalidArgument,
				"conflicting order by fields",
			)
		}

		switch order {
		case sportsapi.ListEventsRequest_ADVERTISED_START_TIME_ASC:
			_, _ = w.WriteString("advertised_start_time ASC")
		case sportsapi.ListEventsRequest_ADVERTISED_START_TIME_DESC:
			_, _ = w.WriteString("advertised_start_time DESC")
		case sportsapi.ListEventsRequest_NAME_ASC:
			_, _ = w.WriteString("name ASC")
		case sportsapi.ListEventsRequest_NAME_DESC:
			_, _ = w.WriteString("name DESC")
		case sportsapi.ListEventsRequest_COMPETITION_ASC:
			_, _ = w.WriteString("competition ASC")
		case sportsapi.ListEventsRequest_COMPETITION_DESC:
			_, _ = w.WriteString("competition DESC")
		}

		visited[order] = true

		// Add a comma if it's not the last element.
		if i != len(req.GetOrderBy())-1 {
			_, _ = w.WriteString(", ")
		}
	}

	return w.String(), nil
}

// computeEventStatus computes the status of a sports event based on its
// advertised start time.
func computeEventStatus(
	advertisedStartTime time.Time,
) sportsapi.Event_Status {
	if advertisedStartTime.After(time.Now()) {
		return sportsapi.Event_OPEN
	}
	return sportsapi.Event_CLOSED
}
