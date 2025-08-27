APP_NAME := http-task-executor

BIN_DIR := bin
BUILD_PATH := $(BIN_DIR)/$(APP_NAME)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BUILD_PATH):
	mkdir -p $(BUILD_PATH)

.PHONY: build
build: $(BIN_DIR)
	go build -o $(BUILD_PATH) ./...

TASK_SERVICE_CONFIG ?= ./task-service/config/local.yaml
TASK_EXECUTOR_CONFIG ?= ./task-executor/config/local.yaml

.PHONY: run-service
run-service:
	./$(BUILD_PATH)/task-service.exe --config=$(TASK_SERVICE_CONFIG)

.PHONY: run-executor
run-executor:
	./$(BUILD_PATH)/task-executor.exe --config=$(TASK_EXECUTOR_CONFIG)


.PHONY: run
run: build
	bash -c 'make run-service & make run-executor'


.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)