SHELL := /bin/sh

APP_NAME := todo-app
CMD_DIR  := ./cmd
BIN_DIR  := bin

.PHONY: all build run test fmt vet clean dev

all: build

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run:
	air -c .air.toml

run1: build
	exec ./$(BIN_DIR)/$(APP_NAME)

dev: fmt vet test run

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BIN_DIR)/$(APP_NAME)
	rm -rf tmp
