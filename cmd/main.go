package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	todoRepository "todo-app/internal/todo/infrastructure/repository"
	todoHttp "todo-app/internal/todo/interface/http"
	todoUsecase "todo-app/internal/todo/usecase"

	authRepository "todo-app/internal/auth/infrastructure/repository"
	authHttp "todo-app/internal/auth/interface/http"
	authUsecase "todo-app/internal/auth/usecase"

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
	apiPort := os.Getenv("API_PORT")
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

	todoRepo := todoRepository.NewMongoTodoRepository(db)
	todoUc := todoUsecase.NewTodoUseCase(todoRepo)
	authRepo := authRepository.NewMemoryRepo()
	loginUc := authUsecase.NewLoginUsecase(authRepo)
	registerUc := authUsecase.NewRegisterUsecase(authRepo)
	jwtSecret := os.Getenv("JWT_SECRET")

	log.Println("Starting Huma server on :" + apiPort + "...\n")
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
	todoHttp.NewTodoHandler(api, todoUc)
	authHttp.NewHandler(api, registerUc, loginUc, []byte(jwtSecret))

	if err := netHttp.ListenAndServe(":"+apiPort, router); err != nil {
		log.Fatalf("Failed to start Huma server: %v", err)
	}
}
