package repo

import (
  "context"
  "database/sql"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type UserRepo struct {
  DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
  return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (int64, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var id int64
  err := r.DB.QueryRowContext(ctx,
    `INSERT INTO users (name, email, password_hash, phone, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
    user.Name, user.Email, user.PasswordHash, user.Phone, user.Role,
  ).Scan(&id)
  return id, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var user model.User
  err := r.DB.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var user model.User
  err := r.DB.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var users []model.User
  err := r.DB.SelectContext(ctx, &users, `SELECT id, name, email, phone, role, created_at FROM users ORDER BY id DESC`)
  return users, err
}
