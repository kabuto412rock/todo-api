package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	"todo-app/internal/infrastructure/repository"
	"todo-app/internal/interface/http"
	"todo-app/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
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

	repo := repository.NewMongoTodoRepository(db)
	uc := usecase.NewTodoUseCase(repo)

	humaPort := "8081"
	log.Println("Starting Huma server on :" + humaPort + "...\n")
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
	http.NewTodoHandler(api, uc)
	if err := netHttp.ListenAndServe(":"+humaPort, router); err != nil {
		log.Fatalf("Failed to start Huma server: %v", err)
	}
}
