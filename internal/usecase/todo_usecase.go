package usecase

import (
	"time"
	"todo-app/internal/domain"

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
	dueTime, err := parseDate(dueDate)
	if err != nil {
		return err
	}
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
func generateID() string {
	return uuid.New().String()
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, dateStr) // 這只是示例，實際應該解析 dateStr
}
