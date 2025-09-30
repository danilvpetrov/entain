# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.7.0] - 2025-09-30

### Added

- OTEL tracing support was added for gateway, racing and sports services.
  The tracing can be enabled by setting the `OTEL_EXPORTER_OTLP_ENDPOINT`
  environment variable to the desired OTLP collector endpoint. For more details,
  please refer to [tracing in README.md](./README.md#tracing).

## [v0.6.0] - 2025-09-29

### Added

- Added `sports` service that serves sports events similar to the racing service.
  The service includes `ListEvents` and `GetEvent` RPCs to list all
  sports events and get a specific event by its ID. For more details, please refer to
  [sports service section in README.md](./README.md#sports-service).

## [v0.5.0] - 2025-09-27

### Added

- Added `GetRace` RPC to the racing service to retrieve a specific race by its
  ID. For more details, please refer to [getting a specific race in README.md](./README.md#getting-a-specific-race).

## [v0.4.0] - 2025-09-27

### Added

- Added `status` field to the result of `ListRaces` RPC in the racing service.
  The field can have values `OPEN` or `CLOSED` depending on whether the race is
  currently open for betting or not.

## [v0.3.0] - 2025-09-27

### Added

- Added ordering functionality to the `ListRaces` RPC in the racing service,
  allowing clients to specify the order of returned races. For more details,
  please refer to [race ordering in README.md](./README.md#ordering-of-races).

## [v0.2.1] - 2025-09-27

### Fixed

- Fixed internal database indexing for the `visibleOnly` filter in the
  `ListRaces` RPC in the racing service.

## [v0.2.0] - 2025-09-27

### Added

- Added `visibleOnly` filter to the `ListRaces` RPC in the racing service which
  allows clients to request only visible races.

## [v0.1.0] - 2025-09-26

### Added

- Initial implementation of the racing service
- API Gateway to route requests to the racing service
- Swagger OpenAPI definitions for the racing service

### Changed

- Updated README.md with instructions for running the services and making requests
- Improved code structure and organization
- Added Makefile for building and running the services

[unreleased]: https://github.com/danilvpetrov/entain/compare/v0.7.0...HEAD
[v0.7.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.7.0
[v0.6.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.6.0
[v0.5.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.5.0
[v0.4.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.4.0
[v0.3.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.3.0
[v0.2.1]: https://github.com/danilvpetrov/entain/releases/tag/v0.2.1
[v0.2.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.2.0
[v0.1.0]: https://github.com/danilvpetrov/entain/releases/tag/v0.1.0
