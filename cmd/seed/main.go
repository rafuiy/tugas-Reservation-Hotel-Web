package main

import (
  "log"
  "os"

  "project_ap3/internal/db"
  "project_ap3/internal/model"
  "project_ap3/internal/repo"
  "project_ap3/internal/util"
)

func main() {
  databaseURL := os.Getenv("DATABASE_URL")
  email := os.Getenv("ADMIN_SEED_EMAIL")
  password := os.Getenv("ADMIN_SEED_PASSWORD")

  if databaseURL == "" || email == "" || password == "" {
    log.Fatal("DATABASE_URL, ADMIN_SEED_EMAIL, and ADMIN_SEED_PASSWORD are required")
  }

  database, err := db.New(databaseURL)
  if err != nil {
    log.Fatal(err)
  }
  defer database.Close()

  users := repo.NewUserRepo(database)
  existing, err := users.GetByEmail(nil, email)
  if err != nil {
    log.Fatal(err)
  }
  if existing != nil {
    log.Println("admin already exists")
    return
  }

  hash, err := util.HashPassword(password)
  if err != nil {
    log.Fatal(err)
  }

  admin := &model.User{
    Name:         "Admin",
    Email:        email,
    PasswordHash: hash,
    Phone:        "",
    Role:         "ADMIN",
  }
  if _, err := users.Create(nil, admin); err != nil {
    log.Fatal(err)
  }

  log.Println("admin seeded")
}
