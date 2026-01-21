package config

import (
  "errors"
  "os"
)

type Config struct {
  DatabaseURL       string
  Port              string
  JWTSecret         string
  AdminSeedEmail    string
  AdminSeedPassword string
}

func Load() (Config, error) {
  cfg := Config{
    DatabaseURL:       os.Getenv("DATABASE_URL"),
    Port:              os.Getenv("PORT"),
    JWTSecret:         os.Getenv("JWT_SECRET"),
    AdminSeedEmail:    os.Getenv("ADMIN_SEED_EMAIL"),
    AdminSeedPassword: os.Getenv("ADMIN_SEED_PASSWORD"),
  }
  if cfg.Port == "" {
    cfg.Port = "8080"
  }
  if cfg.DatabaseURL == "" {
    return cfg, errors.New("DATABASE_URL is required")
  }
  if cfg.JWTSecret == "" {
    return cfg, errors.New("JWT_SECRET is required")
  }
  return cfg, nil
}
