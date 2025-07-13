package domain

type TodoRepository interface {
	Save(todo *Todo) error
	FindAll() ([]*Todo, error)
	DeleteByID(id string) error
	FindByID(id string) (*Todo, error)
	UpdateByID(todo *Todo) error
}
