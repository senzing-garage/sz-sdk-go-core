# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

sz-sdk-go-core is a Go SDK that wraps the Senzing C SDK APIs using CGO. It implements the interfaces defined in [sz-sdk-go](https://github.com/senzing-garage/sz-sdk-go) to call Senzing's entity resolution C libraries directly.

## Prerequisites

Senzing C libraries must be installed before building or testing:

- `/opt/senzing/er/lib` - Senzing shared libraries
- `/opt/senzing/er/sdk/c` - Senzing C header files
- `/etc/opt/senzing` - Senzing configuration

See [How to Install Senzing for Go Development](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md)

## Common Commands

```bash
# Install development dependencies (one-time)
make dependencies-for-development

# Update Go dependencies
make dependencies

# Build
make clean build

# Run tests (requires setup first)
make clean setup test

# Run a single test
go test -v -run TestSzEngine_AddRecord ./szengine/...

# Run tests with coverage
make clean setup coverage

# Lint
make lint

# Auto-fix lint issues
make fix

# Run specific linter
make golangci-lint
make govulncheck
make cspell
```

## Environment Variables

- `LD_LIBRARY_PATH` - Path to Senzing libraries (default: `/opt/senzing/er/lib`)
- `SENZING_TOOLS_DATABASE_URL` - Database connection URL (default: `sqlite3://na:na@nowhere/tmp/sqlite/G2C.db`)
- `SENZING_LOG_LEVEL` - Log level (e.g., `INFO`, `TRACE`)

## Architecture

### Package Structure

Each package wraps a corresponding Senzing C library component:

- **szabstractfactory** - Factory for creating all SDK objects, implements `senzing.SzAbstractFactory`
- **szconfig** - Configuration management, wraps `libSzConfig.h`
- **szconfigmanager** - Configuration persistence, wraps `libSzConfigManager.h`
- **szdiagnostic** - System diagnostics, wraps `libSzDiagnostic.h`
- **szengine** - Core entity resolution engine, wraps `libSz.h`
- **szproduct** - Product information and licensing
- **helper** - Internal logging and messaging utilities

### CGO Bindings

Each package uses CGO directives to link against Senzing C libraries:

```go
/*
#include <stdlib.h>
#include "libSz.h"
#cgo linux CFLAGS: -g -I/opt/senzing/er/sdk/c
#cgo linux LDFLAGS: -L/opt/senzing/er/lib -lSz
*/
import "C"
```

### Interface Implementation Pattern

All packages implement interfaces from `github.com/senzing-garage/sz-sdk-go/senzing`. The typical struct pattern:

- Implements the interface methods that call C functions via CGO
- Supports tracing (via `isTrace` field)
- Supports observability (via `observers` subject)
- Includes lifecycle management (`Initialize`, `Destroy`, `Reinitialize`)

### Test Data

- `testdata/sqlite/G2C.db` - SQLite database template for tests
- Tests copy this to `/tmp/sqlite/G2C.db` during setup

## Code Style

- Uses `gofumpt` for formatting
- Uses `golangci-lint` with extensive linter configuration at `.github/linters/.golangci.yaml`
- Error wrapping uses `github.com/senzing-garage/go-helpers/wraperror`
- Tests use `github.com/stretchr/testify/require`
- Tests run sequentially (`-p 1`) due to shared C library state

## Error Codes

Error identifiers follow format `senzing-PPPPnnnn`:

- `6001` - szconfig
- `6002` - szconfigmanager
- `6003` - szdiagnostic
- `6004` - szengine
- `6006` - szproduct
