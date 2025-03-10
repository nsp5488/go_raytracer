APP_NAME = go_raytracer
SRC_DIR  = .
BUILD_DIR = ./build

.PHONY: all build clean run

all: build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

profile: build
	@$(BUILD_DIR)/$(APP_NAME) -cpuprofile=ray_trace.prof
