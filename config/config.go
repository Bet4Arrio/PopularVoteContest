package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Port string
	Env  string
}

// DBConfig selects which database backend to use.
// DB_TYPE: "mongo" (default) | "sqlite" | "postgres"
// DB_DSN:  file path for SQLite, connection string for Postgres (unused for mongo)
type DBConfig struct {
	Type string
	DSN  string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  int // minutes
	RefreshTokenTTL int // hours
}

func Load() *Config {
	_ = godotenv.Load()

	accessTTL, _ := strconv.Atoi(getEnv("JWT_ACCESS_TTL", "15"))
	refreshTTL, _ := strconv.Atoi(getEnv("JWT_REFRESH_TTL", "168"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "3000"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			DSN: getEnv("DB_DSN", ""),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "super-secret-change-in-prod"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
