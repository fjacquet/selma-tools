# Define the output binaries
CLI_BIN = csv-processor-cli
WEB_BIN = csv-processor-web

# Define the source directories
CLI_SRC = ./cmd
WEB_SRC = ./web

# Default target: build both binaries
all: $(CLI_BIN) $(WEB_BIN) docker

test:
	go test ./cmd/
	go test ./internal/csvprocessor

# Build the CLI binary
$(CLI_BIN):
	go build -o bin/$(CLI_BIN) $(CLI_SRC)/main.go
	chmod 755 bin/$(CLI_BIN)

# Build the Web UI binary
$(WEB_BIN):
	go build -o bin/$(WEB_BIN) $(WEB_SRC)/main.go
	chmod 755 bin/$(WEB_BIN)

# Build the Docker image
docker:
	@if [ -n "$(shell docker images -q $(WEB_BIN) 2> /dev/null)" ]; then \
		docker image rm -f $(WEB_BIN); \
	fi
	docker build -t $(WEB_BIN) .

.PHONY: docker clean run-cli run-web run-docker

# Clean up build artifacts
clean:
	rm -f bin/$(CLI_BIN) bin/$(WEB_BIN)
	@if [ -n "$(shell docker images -q $(WEB_BIN) 2> /dev/null)" ]; then \
		docker image rm -f $(WEB_BIN); \
	fi

# Run the CLI binary
run-cli: $(CLI_BIN)
	./bin/$(CLI_BIN) --input testdata/input.csv --output testdata/output.csv

# Run the Web UI binary
run-web: $(WEB_BIN)
	./bin/$(WEB_BIN)

# Run the Docker container
run-docker:
	docker run -d -p 8080:8080 --name csv-web $(WEB_BIN)
