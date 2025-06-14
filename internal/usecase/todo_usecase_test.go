package usecase

import (
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
