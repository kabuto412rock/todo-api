package config

import (
	"log"
	"os"
)

type Config struct {
	APIPort   string
	MongoURI  string
	MongoDB   string
	JWTSecret string
}

func Load() Config {
	return Config{
		APIPort:   getOr("API_PORT", "8080"),
		MongoURI:  must("MONGO_DB_URI"),
		MongoDB:   must("MONGO_DB_NAME"),
		JWTSecret: must("JWT_SECRET"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env %s", k)
	}
	return v
}

func getOr(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
