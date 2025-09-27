# Entain BE Technical Test

## Why the initial setup was changed

I felt that I should change the initial setup of the technical test for the
following reasons. I hope you would understand my decisions and below I will do
my best to back those decisions up.

1. The original setup was arranged through multiple Go modules. That did not
   allow to import `git.neds.sh/matty/entain/api` module into
   `git.neds.sh/matty/entain/racing` unless I started using [Go workspaces](https://go.dev/doc/tutorial/workspaces).
   Instead, I simplified those modules to subpackages within a single module to
   ease importing. All of the subpackes in the proposed layout should end up
   being separate modules (possibly even repositories) and should be versioned
   separately as well.

1. Dealing with the item 1 allowed placing `*.proto` files into a centralized
   location and not repeat those file in multiple folder running into the risk
   of content de-syncing. In my experience, defining API specs (which `*.proto`
   files are) in a centralized location is an effective way to communicate API
   changes and delineate responsibilies especially in a multi-team environment.

1. I changed existing `ListRaces` RPC from `POST` to `GET` HTTP method in
   `api/racing/racing.proto` file. I am pretty sure that there will be opinions
   arguing about potential growth of the filter complexity and that `POST`
   should be a good fit. I do think, however, that `GET` is the more appropriate
   choice here, as it better aligns with the semantics of retrieving data
   without side effects. `GET` requests can be cached, bookmarked, and shared
   more easily, which is beneficial for a listing operation. I do think that
   current complexity of the filter is manageable within the constraints of
   `GET` requests.

1. I converted command line flags to use environment variables instead. I think
   that environment variables are a more common way to configure services in
   production environments. They are easier to manage in containerized
   environments and orchestration systems like Kubernetes. I also believe that
   it is in line with the [12-factor app principles](https://12factor.net) which
   are widely adopted in the industry. The same goes for structural logging
   which is accomplished through `log/slog` package.

1. I replaced abstractions (interfaces) around database queries in `racing`
   package with concrete implementations. The benefits of using real database
   connection in tests outweigh the benefits of using abstractions in this case.
   The tests are more realistic and better reflect the actual behavior of the
   application when interacting with a real database. It also addresses issues
   like SQL syntax errors and schema mismatches that might not be caught when
   using mocks or stubs.

1. I replaced `tools.go` files with `tool` directive in `go.mod` file. This is a
   more modern and cleaner way to manage tool dependencies in Go projects. It
   keeps the `go.mod` file as the single source of truth for all dependencies,
   making it easier to manage and update them.

1. I did a bit of tidying up of the existing code, introducing more readable
   functions, comments and types. I introduced `cmd` folder which is canonical
   in Go project setup.

1. I introduced `Makefile` to deal with testing, code generation based on
   `*.proto` files, etc. Makefiles is one of the build managers out there, it
   might not be the best, but it is surely widely spread enough. I felt it would
   be better to have some building management rather than none.

## API Gateway

API Gateway acts as a reverse proxy, routing requests from clients to the
appropriate microservices. It handles tasks such as request routing,
composition, and protocol (HTTP<->gRPC) translation.

### Running the API Gateway

To run the API Gateway, use the following command in a separate terminal
window/tab:

```bash
make run-gateway
```

The following environment variables can be used to configure the gateway:

- `LISTEN_ADDR` - address to listen on (default: `:8000`)
- `RACING_SERVICE_ADDR` - address of the racing service (default: `localhost:9000`)
- `DEBUG` - enable debug logging (default: `false`)

## Racing service

Racing service is a microservice that provides racing-related data and
functionality. The Swagger OpenAPI definitions of the service calls can be found
[here](./api/racing/racing.swagger.yaml).

### Running the service

To run the service, use the following command in a separate terminal window/tab:

```bash
make run-racing
```

The following environment variables can be used to configure the gateway:

- `LISTEN_ADDR` - address to listen on (default: `:9000`)
- `RACING_DB_PATH` - path to the racing database (default: `artefacts/racing.db`)
- `DEBUG` - enable debug logging (default: `false`)

### Calling the service through API Gateway

Once the gateway and racing service are running, you can call the service using
`curl` or any HTTP client in a separate terminal window/tab. Provided that the
default address of the gateway is still `:8000`, you can use the following
command:

```bash
curl -i -X GET http://localhost:8000/v1/races
```

### Listing races

You can use the `ListRaces` RPC to list all races. For example:

```bash
curl -i -X GET http://localhost:8000/v1/races
```

#### Filtering races

You can use `meetingId` query parameter to filter the races by meeting ID. You
can use this parameter multiple times to filter by multiple meeting IDs, for example:

```bash
curl -i -X GET "http://localhost:8000/v1/races?meetingId=1&meetingId=2"
```

You can also use `visibleOnly` query parameter to filter only visible races. For example:

```bash
curl -i -X GET "http://localhost:8000/v1/races?visibleOnly=true"
```

Please note that that if `visibleOnly` is set to false or not set at all, both
visible and non-visible races will be returned.

#### Ordering of races

You can use `orderBy` query parameter to order the races by different fields. The
possible values are:

- `ADVERTISED_START_TIME_ASC` - order by advertised start time in ascending order
- `ADVERTISED_START_TIME_DESC` - order by advertised start time in descending order
- `MEETING_ID_ASC` - order by meeting ID in ascending order
- `MEETING_ID_DESC` - order by meeting ID in descending order
- `NAME_ASC` - order by name in ascending order
- `NAME_DESC` - order by name in descending order
- `NUMBER_ASC` - order by number in ascending order
- `NUMBER_DESC` - order by number in descending order

You can use this parameter multiple times to order by multiple fields. The sequence
of the parameters defines the order of precedence. In the example below, the
races will be ordered first by advertised start time in ascending order, and then
by meeting ID in descending order.

```bash
curl -i -X GET "http://localhost:8000/v1/races?orderBy=ADVERTISED_START_TIME_ASC&orderBy=MEETING_ID_DESC"
```

Please note that if you specify conflicting ordering options (e.g.,
`ADVERTISED_START_TIME_ASC` and `ADVERTISED_START_TIME_DESC`), the service will
return an error.

### Getting a specific race

To get a specific race, you can use the `GetRace` RPC and specify the race ID at
the end of the URL. For example:

```bash
curl -i -X GET http://localhost:8000/v1/races/1
```

This will return the details of the race with ID 1.

## Testing

To run unit tests of all services, use the following command:

```bash
make test
```

## Code generation

To generate code from `*.proto` files, use the following command:

```bash
make generate
```

## Technical tasks

### Task 1

> Add another filter to the existing RPC, so we can call `ListRaces` asking for
> races that are visible only

The changes related to this task can be traced at the branch
[`visible-only-races-filter`](https://github.com/danilvpetrov/entain/tree/visible-only-races-filter).
The commits on this branch should be merged into the `main` branch.

### Task 2

> We'd like to see the races returned, ordered by their `advertised_start_time`
> Bonus points if you allow the consumer to specify an ORDER/SORT-BY they might be after.

The changes related to this task can be traced at the branch
[`race-result-ordering`](https://github.com/danilvpetrov/entain/tree/race-result-ordering).
The commits on this branch should be merged into the `main` branch.

### Task 3

> Our races require a new `status` field that is derived based on their
> `advertised_start_time`'s. The status is simply, `OPEN` or `CLOSED`. All races
> that have an `advertised_start_time` in the past should reflect `CLOSED`.

The changes related to this task can be traced at the branch
[`race-status`](https://github.com/danilvpetrov/entain/tree/race-status).
The commits on this branch should be merged into the `main` branch.
