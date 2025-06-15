package main

import (
	"context"
	"log"
	"os"
	"todo-app/internal/infrastructure/repository"
	"todo-app/internal/interface/http"
	"todo-app/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 載入 .env 檔案，路徑為 cmd 資料夾的上層
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	dbClientURI := os.Getenv("MONGO_DB_URI")

	// 直接用 URI 建立 MongoDB 連線選項
	clientOptions := options.Client().ApplyURI(dbClientURI)

	ctx := context.TODO()
	log.Printf("Connecting to MongoDB at %s...\n", dbClientURI)
	client, err := mongo.Connect(ctx, clientOptions)

	// 建立 MongoDB 連線
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
	}()
	// ping 一下看看能不能通
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	db := client.Database(dbName)

	r := gin.Default()
	repo := repository.NewMongoTodoRepository(db)
	uc := usecase.NewTodoUseCase(repo)
	http.NewTodoHandler(r, uc)

	r.Run(":8080")
}
