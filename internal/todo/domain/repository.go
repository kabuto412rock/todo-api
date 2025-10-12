package domain

type TodoRepository interface {
	Save(todo *Todo) error
	FindAll(page, limit int) (list []*Todo, total int64, err error)
	DeleteByID(id string) error
	FindByID(id string) (*Todo, error)
	UpdateByID(todo *Todo) error
}
