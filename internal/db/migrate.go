package db

import (
  "context"
  "os"
  "path/filepath"
  "sort"
  "strings"
  "time"

  "github.com/jmoiron/sqlx"
)

func Migrate(database *sqlx.DB, dir string) error {
  if dir == "" {
    dir = "migrations"
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  if _, err := database.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS schema_migrations (
      filename TEXT PRIMARY KEY,
      applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
    )
  `); err != nil {
    return err
  }

  entries, err := os.ReadDir(dir)
  if err != nil {
    return err
  }

  files := make([]string, 0, len(entries))
  for _, entry := range entries {
    if entry.IsDir() {
      continue
    }
    name := entry.Name()
    if strings.HasSuffix(name, ".sql") {
      files = append(files, name)
    }
  }
  sort.Strings(files)

  applied := map[string]bool{}
  rows, err := database.QueryxContext(ctx, `SELECT filename FROM schema_migrations`)
  if err != nil {
    return err
  }
  for rows.Next() {
    var name string
    if err := rows.Scan(&name); err != nil {
      _ = rows.Close()
      return err
    }
    applied[name] = true
  }
  _ = rows.Close()

  for _, name := range files {
    if applied[name] {
      continue
    }
    path := filepath.Join(dir, name)
    content, err := os.ReadFile(path)
    if err != nil {
      return err
    }
    if strings.TrimSpace(string(content)) == "" {
      continue
    }

    tx, err := database.BeginTxx(ctx, nil)
    if err != nil {
      return err
    }
    if _, err := tx.ExecContext(ctx, string(content)); err != nil {
      _ = tx.Rollback()
      return err
    }
    if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, name); err != nil {
      _ = tx.Rollback()
      return err
    }
    if err := tx.Commit(); err != nil {
      return err
    }
  }

  if err := ensureBookingDateColumns(ctx, database); err != nil {
    return err
  }

  return nil
}

func ensureBookingDateColumns(ctx context.Context, database *sqlx.DB) error {
  _, err := database.ExecContext(ctx, `
    ALTER TABLE bookings
      ADD COLUMN IF NOT EXISTS start_date DATE,
      ADD COLUMN IF NOT EXISTS end_date DATE,
      ADD COLUMN IF NOT EXISTS total_days INTEGER,
      ADD COLUMN IF NOT EXISTS total_amount BIGINT
  `)
  return err
}
