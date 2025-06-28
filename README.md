Learning Clean Architecture from ChatGPT and implementing it with Golang
```bash
todo-app/
├── cmd
│   └── main.go                 # Entry point (main)
├── internal
│   ├── domain                  # Entity + Interface definitions
│   │   ├── repository.go
│   │   └── todo.go
│   ├── infrastructure          # Actual implementations (DB, third-party, etc.)
│   │   └── repository
│   │       └── mongo_todo_repository.go
│   ├── interface               # Interface layer (handler/controller)
│   │   └── http
│   │       ├── dto.go
│   │       └── todo_handler.go
│   └── usecase                 # UseCase implementations
│       ├── todo_usecase.go
│       └── todo_usecase_test.go
```
## Description

A simple Todo REST API built with Go, following Clean Architecture principles.

## Prerequisites

- Go 1.18 or higher
- MongoDB instance
- [Optional] Docker for running MongoDB locally

## Environment Variables

Create a `.env` file at the project root and set the following variables:

```bash
MONGO_URI="mongodb://localhost:27017"
DB_NAME="todoapp"
PORT="8080"
```

## Installation

```bash
go mod tidy
```

## Running the server

```bash
cp .env.example .env
go run cmd/main.go
```

## API Endpoints

| Method | Route      | Description             |
| ------ | ---------- | ----------------------- |
| GET    | /todos     | List all todos          |
| GET    | /todos/:id | Get a single todo       |
| POST   | /todos     | Create a new todo       |
| PUT    | /todos/:id | Update an existing todo |
| DELETE | /todos/:id | Delete a todo           |

## Running Tests

```bash
go test ./internal/usecase
```

