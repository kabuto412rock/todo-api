package repository

import (
	"time"
	"todo-app/internal/domain"
)

type MongoTodoRepository struct {
	// 例如: client *mongo.Client
}

func NewMongoTodoRepository() *MongoTodoRepository {
	return &MongoTodoRepository{}
}

func (r *MongoTodoRepository) Save(todo *domain.Todo) error {
	// 這裡會呼叫MongoDB的Driver來實作實際儲存
	// 假設為簡化範例，這裡不實作Mongo實作
	return nil
}

func (r *MongoTodoRepository) FindAll() ([]*domain.Todo, error) {

	return []*domain.Todo{
		{
			ID:      "1",
			Title:   "Sample Todo",
			DueDate: time.Now(),
		},
	}, nil
}
