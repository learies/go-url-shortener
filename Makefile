# Define the binary name
BINARY_NAME=shortener

# Define the build directory
BUILD_DIR=cmd/shortener

.PHONY: all clean run run_with_flag

# Default target
all: build

# Build the Go binary
build:
	cd $(BUILD_DIR) && go build -buildvcs=false -o $(BINARY_NAME)

# Run the Go binary
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run the Go binary with the file storage flag
run_with_flag:
	cd $(BUILD_DIR) && go run . -f ../../urls.json

# Clean the build artifacts
clean:
	cd $(BUILD_DIR) && rm -f $(BINARY_NAME)
