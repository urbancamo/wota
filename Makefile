.PHONY: all clean install test build-sotaalerts build-sotaspots build-cgitest build-uploader-check build-uploader-create build-uploader-process build-dbtest build-testcsv build-examples help

# Build output directory
BUILD_DIR := bin

# Go build flags
LDFLAGS := -ldflags "-s -w"
GOFLAGS := -trimpath

# Module name
MODULE := github.com/urbancamo/wota

# Default target
all: clean build

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Clean and build all executables"
	@echo "  build            - Build all executables"
	@echo "  clean            - Remove build artifacts"
	@echo "  install          - Install all executables to GOPATH/bin"
	@echo "  test             - Run tests"
	@echo "  build-sotaalerts - Build SOTA alerts fetcher"
	@echo "  build-sotaspots  - Build SOTA spots fetcher"
	@echo "  build-cgitest    - Build CGI test tool"
	@echo "  build-uploader-* - Build SOTA uploader CGI tools"
	@echo "  build-dbtest     - Build database test tool"
	@echo "  build-testcsv    - Build CSV parser test tool"

# Build all executables
build: build-sotaalerts build-sotaspots build-cgitest build-uploader-check build-uploader-create build-uploader-process build-dbtest build-testcsv build-examples
	@echo "Build complete!"

# Build individual executables
build-sotaalerts: $(BUILD_DIR)
	@echo "Building sotaalerts..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/sotaalerts ./sotaalerts

build-sotaspots: $(BUILD_DIR)
	@echo "Building sotaspots..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/sotaspots ./sotaspots

build-uploader-check: $(BUILD_DIR)
	@echo "Building SOTA uploader check form..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/sota-uploader-check.cgi ./sotauploader/check

build-uploader-create: $(BUILD_DIR)
	@echo "Building SOTA uploader create form..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/sota-uploader-create.cgi ./sotauploader/create

build-uploader-process: $(BUILD_DIR)
	@echo "Building SOTA uploader process..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/sota-uploader-process.cgi ./sotauploader/process

build-dbtest: $(BUILD_DIR)
	@echo "Building database test tool..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/dbtest ./sotauploader/dbtest

build-testcsv: $(BUILD_DIR)
	@echo "Building CSV parser test tool..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/testcsv ./sotauploader/test

build-examples: $(BUILD_DIR)
	@echo "Building example programs..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/bandplan_example ./examples

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

# Install executables to GOPATH/bin
install: build
	@echo "Installing executables..."
	@go install $(LDFLAGS) ./sotaalerts
	@go install $(LDFLAGS) ./sotaspots
	@go install $(LDFLAGS) ./cgitest
	@go install $(LDFLAGS) ./sotauploader/check
	@go install $(LDFLAGS) ./sotauploader/create
	@go install $(LDFLAGS) ./sotauploader/process
	@go install $(LDFLAGS) ./sotauploader/dbtest
	@go install $(LDFLAGS) ./sotauploader/test
	@echo "Install complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

# Verify module
verify:
	@echo "Verifying module..."
	@go mod verify
