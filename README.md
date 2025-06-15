從ChatGPT學習Clean Architecture並用golang實踐
```bash
todo-app/
├── cmd/                  # 啟動入口（main）
│   └── main.go
├── internal/
│   ├── domain/           # Entity + Interface 定義
│   │   ├── todo.go
│   │   └── repository.go
│   ├── usecase/          # UseCase 實作
│   │   └── todo_usecase.go
│   ├── interface/        # 交互層 (handler/controller)
│   │   └── http/
│   │       └── todo_handler.go
│   └── infrastructure/   # 實際實作（DB、第三方等）
│       └── repository/
│           └── mongo_todo_repository.go

# 1. clone a .env
cp .env.example .env

# 2. make server run!
go run cmd/main.go
```