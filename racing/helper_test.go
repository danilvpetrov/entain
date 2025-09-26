package racing_test

import (
	"database/sql"
	"net"
	"testing"

	racingapi "github.com/danilvpetrov/entain/api/racing"
	. "github.com/danilvpetrov/entain/racing"
	_ "github.com/mattn/go-sqlite3" // underscore import for the SQLite driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// setupDatabase is a test helper that sets up a test database, seeds it with
// test data, and returns a connection to it.
func setupDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := ApplySchema(t.Context(), db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	if err := SeedTestData(t.Context(), db); err != nil {
		t.Fatal(err)
	}

	return db
}

// setupServer is test helper that sets up a gRPC server for testing and returns
// a client connected to it.
func setupServer(
	t *testing.T,
	s racingapi.RacingServer,
) racingapi.RacingClient {
	t.Helper()

	server := grpc.NewServer()
	racingapi.RegisterRacingServer(server, s)

	// Listen on a random port.
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(server.GracefulStop)

	go func() {
		_ = server.Serve(l)
	}()

	conn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	return racingapi.NewRacingClient(conn)
}
