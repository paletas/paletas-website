# Define the name of the binary
WEBAPP_BINARY_NAME=paletas.webapp

# Define the source files
SRC_FILES=$(wildcard cmd/server/**/*.go) $(wildcard internal/**/*.go)

# Define the build flags
BUILD_FLAGS=-ldflags="-s -w"

# Define the run flags
RUN_FLAGS=--port=8080

# Define the build target
build:
	go build $(BUILD_FLAGS) -o $(WEBAPP_BINARY_NAME) $(SRC_FILES)

# Define the run target
run: clean build
	chmod +x $(WEBAPP_BINARY_NAME)
	./$(WEBAPP_BINARY_NAME) $(RUN_FLAGS)

# Define the clean target
clean:
	rm -f $(WEBAPP_BINARY_NAME)

# Define the test target
test:
	go test -v ./...