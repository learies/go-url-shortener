# Define the binary name
BINARY_NAME=shortener

# Define the build directory
BUILD_DIR=cmd/shortener

.PHONY: all clean run

# Default target
all: build

# Build the Go binary
build:
	cd $(BUILD_DIR) && go build -buildvcs=false -o $(BINARY_NAME)

# Run the Go binary
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean the build artifacts
clean:
	cd $(BUILD_DIR) && rm -f $(BINARY_NAME)
