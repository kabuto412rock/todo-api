package repository

import (
	"context"
	"errors"
	"time"
	"todo-app/internal/todo/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"_id":       todo.ID,
		"title":     todo.Title,
		"dueDate":   todo.DueDate,
		"done":      todo.Done,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})
	return err
}

func (r *MongoTodoRepository) FindAll(page, limit int) (list []*domain.Todo, total int64, err error) {
	skip := int64(page * limit)
	qLimit := int64(limit)
	filter := bson.M{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, filter, &options.FindOptions{
		Sort:  bson.D{{Key: "createdAt", Value: -1}}, // Sort by createdAt in descending order
		Skip:  &skip,
		Limit: &qLimit,
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	total, err = r.collection.CountDocuments(ctx2, filter)
	if err != nil {
		return nil, 0, err
	}

	var todos []*domain.Todo
	for cursor.Next(ctx) {
		var item struct {
			ID      string    `bson:"_id"`
			Title   string    `bson:"title"`
			DueDate time.Time `bson:"dueDate"`
			Done    bool      `bson:"done"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, total, err
		}
		todos = append(todos, &domain.Todo{
			ID:      item.ID,
			Title:   item.Title,
			DueDate: item.DueDate,
			Done:    item.Done,
		})
	}
	return todos, total, nil
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
		Done    bool      `bson:"done"`
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
		Done:    item.Done,
	}, nil
}

func (r *MongoTodoRepository) UpdateByID(todo *domain.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": todo.ID}, bson.M{
		"$set": bson.M{
			"title":     todo.Title,
			"dueDate":   todo.DueDate,
			"updatedAt": time.Now(),
			"done":      todo.Done,
		},
	})
	if err == mongo.ErrNoDocuments {
		return nil // No error if no document found
	}
	return err
}
