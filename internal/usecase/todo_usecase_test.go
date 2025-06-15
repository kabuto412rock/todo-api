package usecase

import (
	"fmt"
	"testing"
	"time"
	"todo-app/internal/domain"

	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	stored []*domain.Todo
}

func (m *mockTodoRepository) Save(todo *domain.Todo) error {
	m.stored = append(m.stored, todo)
	return nil
}

func (m *mockTodoRepository) FindAll() ([]*domain.Todo, error) {
	return m.stored, nil
}
func (m *mockTodoRepository) DeleteByID(id string) error {
	for i, todo := range m.stored {
		if todo.ID == id {
			m.stored = append(m.stored[:i], m.stored[i+1:]...)
			return nil
		}
	}
	return nil // or an error if not found
}

func (m *mockTodoRepository) FindByID(id string) (*domain.Todo, error) {
	for _, todo := range m.stored {
		if todo.ID == id {
			return todo, nil
		}
	}
	return nil, fmt.Errorf("todo with ID %s not found", id)
}

func (m *mockTodoRepository) UpdateByID(todo *domain.Todo) error {
	for i, existingTodo := range m.stored {
		if existingTodo.ID == todo.ID {
			m.stored[i] = todo
			return nil
		}
	}
	return fmt.Errorf("todo with ID %s not found", todo.ID)
}
func TestCreateTodo(t *testing.T) {
	repo := &mockTodoRepository{}
	uc := NewTodoUseCase(repo)

	title := "Learn Clean Architecture"
	dueDate := "2025-07-01"

	err := uc.CreateTodo(title, dueDate)
	assert.NoError(t, err)

	todos, err := uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
	assert.Equal(t, title, todos[0].Title)

	parsedDue, err := time.Parse("2006-01-02", dueDate)
	assert.NoError(t, err)
	assert.Equal(t, parsedDue.Format("2006-01-02"), todos[0].DueDate.Format("2006-01-02"))
}

func TestGetAllTodos(t *testing.T) {
	repo := &mockTodoRepository{}
	uc := NewTodoUseCase(repo)

	todos, err := uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestDeleteTodo(t *testing.T) {
	repo := &mockTodoRepository{}
	uc := NewTodoUseCase(repo)

	// Create a todo to delete
	title := "Learn Clean Architecture"
	dueDate := "2025-07-01"
	err := uc.CreateTodo(title, dueDate)
	assert.NoError(t, err)

	todos, err := uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))

	// Delete the todo
	err = uc.DeleteTodo(todos[0].ID)
	assert.NoError(t, err)

	// Verify the todo is deleted
	todos, err = uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestGetTodoByID(t *testing.T) {
	repo := &mockTodoRepository{}
	uc := NewTodoUseCase(repo)

	// Create a todo to find
	title := "Learn Clean Architecture"
	dueDate := "2025-07-01"
	err := uc.CreateTodo(title, dueDate)
	assert.NoError(t, err)

	todos, err := uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))

	// Find the todo by ID
	foundTodo, err := uc.GetTodoByID(todos[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundTodo)
	assert.Equal(t, title, foundTodo.Title)
}

func TestUpdateTodo(t *testing.T) {
	repo := &mockTodoRepository{}
	uc := NewTodoUseCase(repo)

	// Create a todo to update
	title := "Learn Clean Architecture"
	dueDate := "2025-07-01"
	err := uc.CreateTodo(title, dueDate)
	assert.NoError(t, err)

	todos, err := uc.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))

	// Update the todo
	todos[0].Title = "Learn Clean Architecture Updated"
	err = uc.UpdateTodo(todos[0].ID, todos[0].Title, todos[0].DueDate.Format("2006-01-02"))
	assert.NoError(t, err)

	// Verify the todo is updated
	updatedTodo, err := uc.GetTodoByID(todos[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, "Learn Clean Architecture Updated", updatedTodo.Title)
}
