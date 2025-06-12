package http

import (
	"net/http"
	"todo-app/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	uc *usecase.TodoUseCase
}

func NewTodoHandler(r *gin.Engine, uc *usecase.TodoUseCase) {
	handler := &TodoHandler{}

	r.POST("/todos", handler.Create)
	r.GET("/todos", handler.List)
}

func (h *TodoHandler) Create(c *gin.Context) {
	var body struct {
		Title   string `json:"title"`
		DueDate string `json:"due_date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := h.uc.CreateTodo(body.Title, body.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "todo created"})
}

func (h *TodoHandler) List(c *gin.Context) {
	todos, _ := h.uc.GetAllTodos()
	c.JSON(http.StatusOK, todos)
}
