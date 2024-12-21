# Define the output binaries
CLI_BIN = bin/csv-processor-cli
WEB_BIN = bin/csv-processor-web

# Define the source directories
CLI_SRC = ./cmd
WEB_SRC = ./web

# Default target: build both binaries
all:  $(CLI_BIN) $(WEB_BIN) docker

# Build the CLI binary
$(CLI_BIN):
	go build -o $(CLI_BIN) $(CLI_SRC)/main.go
	chmod 755  $(CLI_BIN)

# Build the Web UI binary
$(WEB_BIN):
	go build -o $(WEB_BIN) $(WEB_SRC)/main.go
	chmod 755  $(WEB_BIN)


# Build the docker image
docker:
	# docker rmi csv-processor-web
	docker build -t csv-processor-web .

# Clean up build artifacts
clean:
	rm -f $(CLI_BIN) $(WEB_BIN)

# Run the CLI binary
run-cli: $(CLI_BIN)
	./$(CLI_BIN) --input input.csv --output output.csv

# Run the Web UI binary
run-web: $(WEB_BIN)
	./$(WEB_BIN)

run-docker:
	docker run -d -p 8080:8080 --name csv-web csv-processor-web