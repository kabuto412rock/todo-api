package usecase

import (
	"time"
	"todo-app/internal/todo/domain"

	"github.com/google/uuid"
)

type TodoUseCase struct {
	repo domain.TodoRepository
}

func NewTodoUseCase(repo domain.TodoRepository) *TodoUseCase {
	return &TodoUseCase{
		repo: repo,
	}
}

func (uc *TodoUseCase) CreateTodo(title string, dueDate string) error {
	// 處理業務邏輯
	// 注意: UseCase 不知道資料從哪裡來，也不管要存去哪裡
	dueTime := parseDate(dueDate)
	todo := &domain.Todo{
		ID:      generateID(),
		Title:   title,
		DueDate: dueTime,
	}
	return uc.repo.Save(todo)
}

func (uc *TodoUseCase) GetAllTodos() ([]*domain.Todo, error) {
	return uc.repo.FindAll()
}

func (uc *TodoUseCase) DeleteTodo(id string) error {
	return uc.repo.DeleteByID(id)
}

func (uc *TodoUseCase) GetTodoByID(id string) (*domain.Todo, error) {
	todo, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (uc *TodoUseCase) UpdateTodo(id string, title, dueDate string) error {
	dueTime := parseDate(dueDate)
	todo := &domain.Todo{
		ID:      id,
		Title:   title,
		DueDate: dueTime,
	}
	return uc.repo.UpdateByID(todo)
}

func generateID() string {
	return uuid.New().String()
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
