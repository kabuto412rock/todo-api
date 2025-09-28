package main

import (
	"context"
	"log"
	netHttp "net/http"
	"regexp"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	authRepo "todo-app/internal/auth/infrastructure/repository"
	"todo-app/internal/config"
	"todo-app/internal/server"
	todoRepo "todo-app/internal/todo/infrastructure/repository"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load()

	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	maskDBURI := MaskDBURI(cfg.MongoURI)
	log.Printf("Connecting to MongoDB at %s...\n", maskDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	db := client.Database(cfg.MongoDB)

	deps := server.Deps{
		JWTSecret: cfg.JWTSecret,
		AuthRepo:  authRepo.NewMemoryRepo(),
		TokenGen:  &authRepo.JWTTokenGenerator{Secret: cfg.JWTSecret},
		TodoRepo:  todoRepo.NewMongoTodoRepository(db),
	}

	h := server.NewRouter(deps)
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := netHttp.ListenAndServe(cfg.ServerAddress, h); err != nil {
		log.Fatalf("server error: %v", err)
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
