package repository

import (
	"context"
	"errors"
	"time"
	"todo-app/internal/todo/domain"

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

func (r *MongoTodoRepository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("there is no document with the given ID")
	}
	return err
}

func (r *MongoTodoRepository) FindByID(id string) (*domain.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var item struct {
		ID      string    `bson:"_id"`
		Title   string    `bson:"title"`
		DueDate time.Time `bson:"dueDate"`
	}
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("there is no document with the given ID")
		}
		return nil, err
	}
	return &domain.Todo{
		ID:      item.ID,
		Title:   item.Title,
		DueDate: item.DueDate,
	}, nil
}

func (r *MongoTodoRepository) UpdateByID(todo *domain.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": todo.ID}, bson.M{
		"$set": bson.M{
			"title":   todo.Title,
			"dueDate": todo.DueDate,
		},
	})
	if err == mongo.ErrNoDocuments {
		return nil // No error if no document found
	}
	return err
}
