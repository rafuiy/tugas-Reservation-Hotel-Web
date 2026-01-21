package db

import (
  "time"

  "github.com/jmoiron/sqlx"
  _ "github.com/jackc/pgx/v5/stdlib"
)

func New(databaseURL string) (*sqlx.DB, error) {
  db, err := sqlx.Connect("pgx", databaseURL)
  if err != nil {
    return nil, err
  }
  db.SetMaxOpenConns(10)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(30 * time.Minute)
  return db, nil
}
