package main

import (
	"todo-app/internal/infrastructure/repository"
	"todo-app/internal/interface/http"
	"todo-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	repo := repository.NewMongoTodoRepository()
	uc := usecase.NewTodoUseCase(repo)
	http.NewTodoHandler(r, uc)

	r.Run(":8080")
}
