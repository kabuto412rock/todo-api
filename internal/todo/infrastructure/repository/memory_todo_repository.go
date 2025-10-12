package repository

import (
	"errors"
	"sort"
	"sync"
	"time"
	"todo-app/internal/todo/domain"
)

type MemoryTodoRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Todo
}

func NewMemoryTodoRepository() *MemoryTodoRepository {
	return &MemoryTodoRepository{items: map[string]*domain.Todo{}}
}

func (r *MemoryTodoRepository) Save(todo *domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// emulate createdAt ordering by using time
	if todo.DueDate.IsZero() {
		// keep same semantics; just allow empty
	}
	r.items[todo.ID] = &domain.Todo{ID: todo.ID, Title: todo.Title, DueDate: todo.DueDate, Done: todo.Done}
	return nil
}

func (r *MemoryTodoRepository) FindAll(page, limit int) (list []*domain.Todo, total int64, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Todo, 0, len(r.items))
	for _, v := range r.items {
		copy := *v
		res = append(res, &copy)
	}
	if page < 0 || limit <= 0 {
		return []*domain.Todo{}, int64(len(r.items)), nil
	}
	start := page * limit
	if start >= len(res) {
		return []*domain.Todo{}, int64(len(r.items)), nil
	}
	end := start + limit
	if end > len(res) {
		end = len(res)
	}
	res = res[start:end]
	// stable order just by ID for determinism
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res, int64(len(r.items)), nil
}

func (r *MemoryTodoRepository) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return errors.New("there is no document with the given ID")
	}
	delete(r.items, id)
	return nil
}

func (r *MemoryTodoRepository) FindByID(id string) (*domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.items[id]
	if !ok {
		return nil, errors.New("there is no document with the given ID")
	}
	copy := *v
	return &copy, nil
}

func (r *MemoryTodoRepository) UpdateByID(todo *domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[todo.ID]; !ok {
		return nil
	}
	r.items[todo.ID] = &domain.Todo{ID: todo.ID, Title: todo.Title, DueDate: todo.DueDate, Done: todo.Done}
	return nil
}

// helper to seed
func (r *MemoryTodoRepository) seed(title string) *domain.Todo {
	t := &domain.Todo{ID: time.Now().Format("20060102150405.000000"), Title: title, DueDate: time.Now().Add(24 * time.Hour)}
	r.Save(t)
	return t
}
