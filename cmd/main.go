package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	"regexp"
	"todo-app/internal/api/middleware"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	dbClientURI := os.Getenv("MONGO_DB_URI")
	apiPort := os.Getenv("API_PORT")
	clientOptions := options.Client().ApplyURI(dbClientURI)

	ctx := context.TODO()
	maskDBURI := MaskDBURI(dbClientURI)
	log.Printf("Connecting to MongoDB at %s...\n", maskDBURI)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
	}()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	db := client.Database(dbName)

	todoRepo := todoRepository.NewMongoTodoRepository(db)
	todoUc := todoUsecase.NewTodoUseCase(todoRepo)
	jwtSecret := os.Getenv("JWT_SECRET")
	authRepo := authRepository.NewMemoryRepo()
	tokenGen := &authRepository.JWTTokenGenerator{Secret: jwtSecret}
	loginUc := authUsecase.NewLoginUsecase(authRepo, tokenGen)
	registerUc := authUsecase.NewRegisterUsecase(authRepo)

	log.Println("Starting Huma server on :" + apiPort + "...\n")
	router := chi.NewRouter()
	config := huma.DefaultConfig("Todo API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"myAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	api := humachi.New(router, config)
	api.UseMiddleware(middleware.NewAuthMiddleware(api, []byte(jwtSecret)))
	todoHttp.NewTodoHandler(api, todoUc)
	authHttp.NewHandler(api, registerUc, loginUc)

	if err := netHttp.ListenAndServe(":"+apiPort, router); err != nil {
		log.Fatalf("Failed to start Huma server: %v", err)
	}
}

func MaskDBURI(uri string) string {
	re := regexp.MustCompile(`(?i)(mongodb(?:\+srv)?://)(.*?):(.*?)@(.*?)(/|$)`)
	return re.ReplaceAllStringFunc(uri, func(m string) string {
		sub := re.FindStringSubmatch(m)
		if len(sub) < 5 {
			return m
		}
		user := sub[2]
		if len(user) > 2 {
			user = user[:1] + "****" + user[len(user)-1:]
		} else {
			user = "****"
		}
		pass := "****"
		host := sub[4]
		hostParts := regexp.MustCompile(`\.`).Split(host, -1)
		if len(hostParts) > 2 {
			host = hostParts[0][:1] + "****." + hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
		} else if len(host) > 4 {
			host = host[:2] + "****" + host[len(host)-2:]
		} else {
			host = "****"
		}
		return sub[1] + user + ":" + pass + "@" + host + sub[5]
	})
}
