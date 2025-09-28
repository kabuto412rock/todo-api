package repository

import (
	"context"
	"errors"
	"time"
	"todo-app/internal/auth/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAuthRepository implements domain.AuthRepository using MongoDB.
type MongoAuthRepository struct {
	collection *mongo.Collection
}

// NewMongoAuthRepository creates a new auth repository backed by the given DB.
// It also ensures a unique index on username.
func NewMongoAuthRepository(db *mongo.Database) *MongoAuthRepository {
	coll := db.Collection("auth_users")
	// Ensure unique index on username
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_username"),
	})
	return &MongoAuthRepository{collection: coll}
}

func (r *MongoAuthRepository) CreateUser(user domain.AuthUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":           user.Username, // use username as _id for natural uniqueness
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"createdAt":     time.Now(),
		"updatedAt":     time.Now(),
	})
	if err != nil {
		// Normalize duplicate key error to a simple error string as memory repo
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user exists")
		}
		return err
	}
	return nil
}

func (r *MongoAuthRepository) GetUserByUsername(username string) (domain.AuthUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc struct {
		Username     string `bson:"username"`
		PasswordHash string `bson:"password_hash"`
	}
	err := r.collection.FindOne(ctx, bson.M{"_id": username}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.AuthUser{}, errors.New("user [" + username + "] not found")
		}
		return domain.AuthUser{}, err
	}
	return domain.AuthUser{Username: doc.Username, PasswordHash: doc.PasswordHash}, nil
}
