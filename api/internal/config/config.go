package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPass      string
	ServerPort  int
	JWTSecret   string
	CORSOrigins string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvInt("DB_PORT", 3306),
		DBName:      getEnv("DB_NAME", "pos"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPass:      getEnv("DB_PASS", ""),
		ServerPort:  getEnvInt("SERVER_PORT", 8080),
		JWTSecret:   getEnv("JWT_SECRET", "pos-secret-change-me"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:8000,http://localhost:8080"),
	}

	if cfg.JWTSecret == "" || cfg.JWTSecret == "pos-secret-change-me" {
		log.Fatal("JWT_SECRET environment variable must be set to a strong random value")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
