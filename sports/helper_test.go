package sports_test

import (
	"database/sql"
	"net"
	"testing"

	sportsapi "github.com/danilvpetrov/entain/api/sports"
	. "github.com/danilvpetrov/entain/sports"
	_ "github.com/mattn/go-sqlite3" // underscore import for the SQLite driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// setupDatabase is a test helper that sets up a test database, seeds it with
// test data, and returns a connection to it along with the number of seeded
// records.
func setupDatabase(t *testing.T) (_ *sql.DB, numOfRecords int) {
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

	numRecords, err := SeedTestData(
		t.Context(),
		db,
		"testdata/testdata.json",
	)
	if err != nil {
		t.Fatal(err)
	}

	return db, numRecords
}

// setupServer is test helper that sets up a gRPC server for testing and returns
// a client connected to it.
func setupServer(
	t *testing.T,
	s sportsapi.SportsServer,
) sportsapi.SportsClient {
	t.Helper()

	server := grpc.NewServer()
	sportsapi.RegisterSportsServer(server, s)

	listenCfg := net.ListenConfig{}
	// Listen on a random port.
	listener, err := listenCfg.Listen(t.Context(), "tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(server.GracefulStop)

	go func() {
		_ = server.Serve(listener)
	}()

	conn, err := grpc.NewClient(
		listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	return sportsapi.NewSportsClient(conn)
}
