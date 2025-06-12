package domain

type TodoRepository interface {
	Save(todo *Todo) error
	FindAll() ([]*Todo, error)
}
