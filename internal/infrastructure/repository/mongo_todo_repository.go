package repository

import (
	"context"
	"time"
	"todo-app/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTodoRepository struct {
	collection *mongo.Collection
}

func NewMongoTodoRepository(db *mongo.Database) *MongoTodoRepository {

	return &MongoTodoRepository{
		collection: db.Collection("todos"),
	}
}

func (r *MongoTodoRepository) Save(todo *domain.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":     todo.ID,
		"title":   todo.Title,
		"dueDate": todo.DueDate,
	})
	return err
}

func (r *MongoTodoRepository) FindAll() ([]*domain.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*domain.Todo
	for cursor.Next(ctx) {
		var item struct {
			ID      string    `bson:"_id"`
			Title   string    `bson:"title"`
			DueDate time.Time `bson:"dueDate"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		todos = append(todos, &domain.Todo{
			ID:      item.ID,
			Title:   item.Title,
			DueDate: item.DueDate,
		})
	}
	return todos, nil
}
