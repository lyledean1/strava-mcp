# Define variables
BINARY_NAME=strava
GOBIN=$(shell go env GOPATH)/bin
GO_FILES=$(shell find . -name "*.go" -type f)

# Default goal
.PHONY: all
all: lint vet build

.PHONY: build-api
build-api: $(GO_FILES)
	go build -o bin/$(BINARY_NAME)-api ./cmd/strava-api

# Build the application
.PHONY: build
build: $(GO_FILES)
	go build -o bin/$(BINARY_NAME)-mcp ./cmd/strava-mcp

# Install binary to GOBIN directory
.PHONY: install
install: build
	cp bin/$(BINARY_NAME) $(GOBIN)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(GOBIN)"

# Run go vet to check for suspicious constructs
.PHONY: vet
vet:
	go vet ./...
	@echo "go vet passed"

# Run linter (requires golangci-lint to be installed)
.PHONY: lint
lint:
	golangci-lint run
	@echo "Linting passed"

# Install golangci-lint if not present
.PHONY: install-lint
install-lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN))

# Run tests
.PHONY: test
test:
	go test -v ./...

# Check code quality (vet + lint + test)
.PHONY: check
check: vet lint test
	@echo "All checks passed"

# Clean build artifacts
.PHONY: clean
clean:
	rm -f bin/$(BINARY_NAME)*
	@echo "Cleaned up build artifacts"

# Run the application
.PHONY: run
run: build
	./bin/$(BINARY_NAME)-mcp

install-mdbook:
	cargo install mdbook

s serve:
	cd book && mdbook serve

build-book:
	cd book && mdbook build

# Display help information
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make          : Lint, vet, and build the application"
	@echo "  make build    : Build the application"
	@echo "  make build-api: Build the API application"
	@echo "  make vet      : Run go vet"
	@echo "  make lint     : Run golangci-lint"
	@echo "  make test     : Run tests"
	@echo "  make check    : Run vet, lint, and tests"
	@echo "  make install  : Build and install to Go bin directory"
	@echo "  make install-lint : Install golangci-lint"
	@echo "  make clean    : Remove build artifacts"
	@echo "  make run      : Build and run the application"
	@echo "  make help     : Display this help"