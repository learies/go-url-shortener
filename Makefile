BINARY_NAME=shortener

BUILD_DIR=cmd/server

BIN_DIR=bin

# Ensure BIN_DIR exists
create_bin_dir:
	mkdir -p $(BIN_DIR)

run: create_bin_dir
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(BUILD_DIR)/main.go && $(BIN_DIR)/$(BINARY_NAME)

build: create_bin_dir
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(BUILD_DIR)/main.go

clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)

start:
	$(BIN_DIR)/$(BINARY_NAME)
