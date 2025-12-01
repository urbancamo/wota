
# WOTA - Wainwrights on the Air

A collection of Go tools for working with SOTA (Summits on the Air) and WOTA (Wainwrights on the Air) data.

## Overview

This project provides several utilities for fetching SOTA alerts and spots, converting them to WOTA format, and managing activation/chase logs via a web interface.

## Prerequisites

- Go 1.21 or later
- MySQL database (for storing WOTA data)

## Installation

### Using Go Install

Install all tools directly to your `$GOPATH/bin`:

```bash
go install github.com/urbancamo/wota/sotaalerts@latest
go install github.com/urbancamo/wota/sotaspots@latest
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/urbancamo/wota.git
cd wota
```

2. Download dependencies:
```bash
go mod download
```

3. Build all executables:
```bash
make build
```

Executables will be placed in the `bin/` directory.

4. Or install to your `$GOPATH/bin`:
```bash
make install
```

## Available Tools

### Command-Line Tools

- **sotaalerts** - Fetches SOTA alerts and converts Lake District summits to WOTA format
- **sotaspots** - Fetches SOTA spots and converts Lake District summits to WOTA format

### CGI Web Tools

- **sota-uploader-create.cgi** - Web form for creating SOTA/WOTA log entries
- **sota-uploader-check.cgi** - CSV validation and preview tool
- **sota-uploader-process.cgi** - Processes uploaded CSV files

### Testing/Development Tools

- **cgitest** - CGI testing utility
- **dbtest** - Database connection testing
- **testcsv** - CSV parsing validation

## Building

### Build All Tools
```bash
make all
```

### Build Individual Tools
```bash
make build-sotaalerts
make build-sotaspots
make build-uploader-check
```

### Clean Build Artifacts
```bash
make clean
```

### Run Tests
```bash
make test
```

### Show All Available Targets
```bash
make help
```

## Testing

### Run All Tests

Using Make:
```bash
make test
```

Using Go directly:
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Tests for a Specific Package
```bash
go test ./sotautils
go test -v ./sotauploader/csv
```

### Run a Specific Test Function
```bash
go test -v ./sotautils -run TestFunctionName
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Project Structure

```
wota/
├── sotaalerts/         # SOTA alerts fetcher
├── sotaspots/          # SOTA spots fetcher
├── sotauploader/       # Web-based CSV uploader
│   ├── check/          # Form validation
│   ├── create/         # Form creation
│   ├── process/        # Form processing
│   ├── csv/            # CSV parsing
│   ├── db/             # Database access
│   └── utils/          # Utilities
├── sotautils/          # Shared SOTA/WOTA utilities
├── domain/             # Domain models
└── cgitest/            # CGI testing tools
```

## Band Plan Data

The `sotautils` package includes comprehensive UK amateur radio band plan data based on RSGB 2024/2025 allocations:

```go
import "github.com/urbancamo/wota/sotautils"

// Get band from frequency
band := sotautils.GetBandFromFrequencyMHz(14.200)  // Returns "20m"

// Check if frequency is in a specific band
inBand := sotautils.IsFrequencyInBand(7100, "40m")  // Returns true

// Get all available bands
bands := sotautils.GetAllBandNames()

// Access band plan data directly
bandInfo := sotautils.RSGBBandPlan["20m"]
fmt.Printf("%s: %.1f - %.1f kHz\n", bandInfo.Name, bandInfo.LowerFreq, bandInfo.UpperFreq)
```

Run the example program to see all available bands:
```bash
./bin/bandplan_example
```

### Supported Bands

**LF/MF:** 2190m, 630m, 160m
**HF:** 80m, 60m, 40m, 30m, 20m, 17m, 15m, 12m, 10m
**VHF:** 6m, 4m, 2m
**UHF:** 70cm
**Microwave:** 23cm, 13cm, 9cm, 6cm, 3cm

## Development

This project uses Go modules. Import paths follow the pattern:

```go
import "github.com/urbancamo/wota/sotautils"
```

## License

See LICENSE file for details.
