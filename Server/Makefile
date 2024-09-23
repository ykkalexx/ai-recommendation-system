# Variables
APP_NAME = recommendation-system
MAIN_FILE = ./cmd/server/main.go
BUILD_DIR = ./bin
BUILD_FILE = $(BUILD_DIR)/$(APP_NAME)
CONFIG_FILE = ./config/config.yaml

# Default target
all: 
	run

# Run the application
run:
	go run $(MAIN_FILE)

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_FILE) $(MAIN_FILE)

# Test the application
test:
	go test ./...

# Clean up build files
clean:
	rm -rf $(BUILD_DIR)

# Load configuration
config:
	@echo "Loading configuration from $(CONFIG_FILE)"
	@cat $(CONFIG_FILE)

# Help
help:
	@echo "Makefile commands:"
	@echo "  run       - Run the application"
	@echo "  build     - Build the application"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean up build files"
	@echo "  config    - Load configuration"
	@echo "  help      - Show this help message"