BUILD_DIR	:= bin
NAME     	:= test

.PHONY: build
build:
		@go build -o $(BUILD_DIR)/$(NAME) ./src/

.PHONY: run
run:
		@./$(BUILD_DIR)/$(NAME)

.PHONY: fmt
fmt:
		@go fmt

.PHONY: test
test: fmt build run

.PHONY: clean
clean:
		rm -rf ./bin
