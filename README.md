# reg-room-service

## Overview

A backend service that provides room management.

Implemented in go.

Command line arguments
```-config <path-to-config-file> [-migrate-database]```

## Configuration File

There is a template configuration file under `docs/config.example.yaml`. Copy it to `config.yaml` in the service
root (or wherever your `-config` argument points), and edit it to match your requirements.

The sensitive values in the configuration file can also be specified via environment variables, so they can be
configured using a kubernetes secret or vault integration. If set, the environment variables override any
values in the configuration file, which are then allowed to be missing or empty.

| Environment variable   | Overrides configuration value |
|------------------------|-------------------------------|
| REG_SECRET_DB_PASSWORD | database.password             |
| REG_SECRET_API_TOKEN   | database.password             |

## Installation

This service uses go modules to provide dependency management, see `go.mod`.

If you place this repository OUTSIDE of your gopath, `go build cmd/main.go` and
`go test ./...` will download all required dependencies by default.

## Test Coverage

In order to collect full test coverage, set go tool arguments to `-covermode=atomic -coverpkg=./internal/...`,
or manually run
```
go test -covermode=atomic -coverpkg=./internal/... ./...
```

