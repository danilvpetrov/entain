package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

// testdataFile is the path to the JSON file where test data is stored.
const testdataFile = "sports/testdata/testdata.json"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Event represents a sports event with its name, category, and competition.
type Event struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Competition string `json:"competition"`
}

// run is the main execution function that orchestrates fetching and storing
// in-play sports events from the Ladbrokes API.
//
// This code is not pretty. It is intended to be run manually when additional
// test data is required. Please avoid using it in a production setup.
// It retrieves the list of sports categories, fetches in-play events for each
// category, and stores thouse events in a JSON file, deduplicating them
// against any existing events in [testdataFile].
func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// seen keeps track of events that have already been recorded to avoid
	// duplicates.
	seen := map[string]bool{}

	raw, err := os.ReadFile(testdataFile)

	var events []Event
	if err == nil {
		if err := json.Unmarshal(raw, &events); err != nil {
			return err
		}

		for _, event := range events {
			seen[event.Name] = true
		}
	}

	categories, err := retrieveCategories(ctx)
	if err != nil {
		return err
	}

	var inPlayEvents []Event
	for _, cat := range categories {
		catEvents, err := getInPlayEvents(ctx, cat)
		if err != nil {
			return err
		}

		inPlayEvents = append(inPlayEvents, catEvents...)
	}

	for _, inPlayEvent := range inPlayEvents {
		if !seen[inPlayEvent.Name] {
			events = append(events, inPlayEvent)
			seen[inPlayEvent.Name] = true
		}
	}

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(
		"sports/testdata/testdata.json",
		data,
		0o600,
	); err != nil {
		return err
	}

	return nil
}

// getInPlayEvents fetches in-play events for a given sports category from the
// Ladbrokes API.
func getInPlayEvents(ctx context.Context, category string) ([]Event, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		`https://api.ladbrokes.com.au/gql/router?variables=%7B%22category%22%3A%22`+
			category+
			`%22%2C%22excludeCategoryIds%22%3A%5B%5D%2C%22includeGroupedInPlayEvents%22%3Afalse`+
			`%2C%22includeInPlayEvents%22%3Atrue%2C%22includeInPlayFilters%22%3Atrue%7D&operation`+
			`Name=SportingInPlayEvents&extensions=%7B%22persistedQuery%22%3A%7B%22version%22`+
			`%3A1%2C%22sha256Hash%22%3A%22d1ec528fea5c0821d6d2cd03ba22b75ec35f2bd63a72bb60afac5a67c65dc691%22%7D%7D`,
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	response := struct {
		Data struct {
			InPlayEvents struct {
				Events struct {
					Nodes []struct {
						Name          string `json:"name"`
						Status        string `json:"status"`
						SportCategory struct {
							Category string `json:"category"`
						} `json:"sportCategory"`
						Competition struct {
							Name string `json:"name"`
						} `json:"competition"`
					} `json:"nodes"`
				} `json:"events"`
			} `json:"inPlayEvents"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	events := make(
		[]Event,
		0,
		len(response.Data.InPlayEvents.Events.Nodes),
	)
	for _, event := range response.Data.InPlayEvents.Events.Nodes {
		events = append(
			events,
			Event{
				Name:        event.Name,
				Category:    event.SportCategory.Category,
				Competition: event.Competition.Name,
			},
		)
	}

	return events, nil
}

// retrieveCategories fetches the list of sports categories from the Ladbrokes
// API.
func retrieveCategories(
	ctx context.Context,
) ([]string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		`https://api.ladbrokes.com.au/gql/router?variables=%7B%22marketControlExclude`+
			`%22%3Afalse%7D&operationName=SportingCategories&extensions=%7B%22persistedQuery`+
			`%22%3A%7B%22version%22%3A1%2C%22sha256Hash%22%3A%2215f`+
			`e182fc110ba95c04a688ae603f3f70553fa412295d2cc1114fca5f6aec4f1%22%7D%7D`,
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	response := struct {
		Data struct {
			Categories []struct {
				Category string `json:"category"`
			} `json:"categories"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	categories := make([]string, 0, len(response.Data.Categories))
	for _, cat := range response.Data.Categories {
		categories = append(categories, cat.Category)
	}

	return categories, nil
}
