package config

import (
  "errors"
  "os"
  "strings"
)

type Config struct {
  DatabaseURL       string
  Port              string
  JWTSecret         string
  AdminSeedEmail    string
  AdminSeedPassword string
}

func Load() (Config, error) {
  loadEnvFile(".env")
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

func loadEnvFile(path string) {
  data, err := os.ReadFile(path)
  if err != nil {
    return
  }
  lines := strings.Split(string(data), "\n")
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") {
      continue
    }
    if strings.HasPrefix(line, "export ") {
      line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
    }
    idx := strings.Index(line, "=")
    if idx <= 0 {
      continue
    }
    key := strings.TrimSpace(line[:idx])
    val := strings.TrimSpace(line[idx+1:])
    if key == "" {
      continue
    }
    if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") && len(val) >= 2 {
      val = strings.Trim(val, "\"")
    }
    if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") && len(val) >= 2 {
      val = strings.Trim(val, "'")
    }
    if _, exists := os.LookupEnv(key); exists {
      continue
    }
    _ = os.Setenv(key, val)
  }
}
