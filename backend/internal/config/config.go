package config

import (
  "os"
)

type Config struct {
  Port string
  DB_DSN string
  AllowedOrigins string
}

func Load() Config {
  return Config{
    Port: getenv("PORT", "8080"),
    DB_DSN: getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/mascotas?sslmode=disable"),
    AllowedOrigins: getenv("ALLOWED_ORIGINS", "http://localhost:3000"),
  }
}

func getenv(key, def string) string {
  if v, ok := os.LookupEnv(key); ok && v != "" {
    return v
  }
  return def
}

