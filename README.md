Learning Clean Architecture from ChatGPT and implementing it with Golang
```bash
todo-app/
├── cmd/                  # Entry point (main)
│   └── main.go
├── internal/
│   ├── domain/           # Entity + Interface definitions
│   │   ├── todo.go
│   │   └── repository.go
│   ├── usecase/          # UseCase implementations
│   │   └── todo_usecase.go
│   ├── interface/        # Interface layer (handler/controller)
│   │   └── http/
│   │       └── todo_handler.go
│   └── infrastructure/   # Actual implementations (DB, third-party, etc.)
│       └── repository/
│           └── mongo_todo_repository.go

# 1. Clone a .env file
cp .env.example .env

# 2. Run the server!
go run cmd/main.go
```