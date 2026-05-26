package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:      getEnv("APP_PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret"),
		DBPath:    getEnv("DB_PATH", "data/airdrop.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Printf("Using default for %s", key)
	return fallback
}
