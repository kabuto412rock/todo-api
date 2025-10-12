package usecase

import (
	"testing"
	"time"
	"todo-app/internal/todo/infrastructure/repository"

	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	repo := repository.NewMemoryTodoRepository()
	uc := NewTodoUseCase(repo)

	title := "Learn Clean Architecture"
	dueDate := parseDate("2025-07-01")

	err := uc.CreateTodo(title, dueDate, false)
	assert.NoError(t, err)

	todos, total, err := uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, int64(1), total)
	assert.Equal(t, title, todos[0].Title)
	assert.Equal(t, false, todos[0].Done)

	assert.NoError(t, err)
	assert.Equal(t, dueDate, todos[0].DueDate)
}

func TestGetAllTodos(t *testing.T) {
	repo := repository.NewMemoryTodoRepository()
	uc := NewTodoUseCase(repo)

	todos, total, err := uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
	assert.Equal(t, int64(0), total)

	title := "Learn Clean Architecture"
	dueDate := parseDate("2025-07-01")
	uc.CreateTodo(title, dueDate, false)
	todos, total, err = uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, int64(1), total)
	assert.Equal(t, title, todos[0].Title)
	assert.Equal(t, false, todos[0].Done)
	assert.Equal(t, dueDate, todos[0].DueDate)

	todos, total, err = uc.GetAllTodos(1, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
	assert.Equal(t, int64(1), total)
	todos, total, err = uc.GetAllTodos(0, 10, "NonExistingTitle")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
	assert.Equal(t, int64(0), total)

	todos, total, err = uc.GetAllTodos(0, 10, "Clean")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, int64(1), total)
}

func TestDeleteTodo(t *testing.T) {
	repo := repository.NewMemoryTodoRepository()
	uc := NewTodoUseCase(repo)

	// Create a todo to delete
	title := "Learn Clean Architecture"
	dueDate := parseDate("2025-07-01")
	err := uc.CreateTodo(title, dueDate, false)
	assert.NoError(t, err)

	todos, total, err := uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, int64(1), total)

	// Delete the todo
	err = uc.DeleteTodo(todos[0].ID)
	assert.NoError(t, err)

	// Verify the todo is deleted
	todos, total, err = uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
	assert.Equal(t, int64(0), total)
}

func TestGetTodoByID(t *testing.T) {
	repo := repository.NewMemoryTodoRepository()
	uc := NewTodoUseCase(repo)

	// Create a todo to find
	title := "Learn Clean Architecture"
	dueDate := parseDate("2025-07-01")
	done := true
	err := uc.CreateTodo(title, dueDate, done)
	assert.NoError(t, err)

	todos, total, err := uc.GetAllTodos(0, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, int64(1), total)

	// Find the todo by ID
	foundTodo, err := uc.GetTodoByID(todos[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundTodo)
	assert.Equal(t, title, foundTodo.Title)
	assert.Equal(t, done, foundTodo.Done)
}

func TestUpdateTodo(t *testing.T) {
	repo := repository.NewMemoryTodoRepository()
	uc := NewTodoUseCase(repo)

	// Create a todo to update
	title := "Learn Clean Architecture"
	dueDate := parseDate("2025-07-01")
	done := false
	err := uc.CreateTodo(title, dueDate, done)
	assert.NoError(t, err)

	todos, total, err := uc.GetAllTodos(0, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, false, todos[0].Done)
	assert.Equal(t, int64(1), total)

	// Update the todo
	todos[0].Title = "Learn Clean Architecture Updated"
	err = uc.UpdateTodo(todos[0].ID, todos[0].Title, todos[0].DueDate, true)
	assert.NoError(t, err)

	// Verify the todo is updated
	updatedTodo, err := uc.GetTodoByID(todos[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "Learn Clean Architecture Updated", updatedTodo.Title)
	assert.Equal(t, true, updatedTodo.Done)
}
func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
