BINARY_NAME=imageboard
BUILD_PATH=bin/$(BINARY_NAME)
MAIN_PATH=$(BINARY_NAME)/main.go

.PHONY: setup clean build run dev all

setup:
	@echo "Setting up environment..."
	@go mod download
	@echo "Environment setup complete."

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Cleanup complete."

build:
	@echo "Building..."
	@go build -o $(BUILD_PATH) $(MAIN_PATH) || true
	@echo "Build complete."

run:
	@if [ ! -f $(BUILD_PATH) ]; then echo "Binary not found. Building binary..."; $(MAKE) -s build; fi
	@echo "Running..."
	@$(BUILD_PATH) || true

dev:
	@echo "Running in development mode..."
	@go run $(MAIN_PATH) || true

all: setup clean build run

.SILENT: