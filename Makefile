# Makefile for Go Todo App

APP_NAME=todo-app
CMD_DIR=./cmd
BIN_DIR=bin

.PHONY: all build run test clean fmt vet

all: build

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run: build
	./$(BIN_DIR)/$(APP_NAME)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BIN_DIR)/$(APP_NAME)
