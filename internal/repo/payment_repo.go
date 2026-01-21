package repo

import (
  "context"
  "database/sql"
  "fmt"
  "time"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type PaymentRepo struct {
  DB *sqlx.DB
}

func NewPaymentRepo(db *sqlx.DB) *PaymentRepo {
  return &PaymentRepo{DB: db}
}

func (r *PaymentRepo) GetByID(ctx context.Context, id int64) (*model.Payment, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var p model.Payment
  err := r.DB.GetContext(ctx, &p, `SELECT * FROM payments WHERE id=$1`, id)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &p, nil
}

func (r *PaymentRepo) ListByUser(ctx context.Context, userID int64) ([]model.PaymentView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT p.id, p.booking_id, p.amount, p.method, p.status, p.paid_at, p.reference, p.created_at,
                   b.user_id, b.status AS booking_status,
                   r.name AS room_name, r.room_no,
                   a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end
            FROM payments p
            JOIN bookings b ON b.id = p.booking_id
            JOIN rooms r ON r.id = b.room_id
            JOIN availability a ON a.id = b.availability_id
            WHERE b.user_id=$1
            ORDER BY p.created_at DESC`

  var rows []model.PaymentView
  err := r.DB.SelectContext(ctx, &rows, query, userID)
  return rows, err
}

func (r *PaymentRepo) ListAdmin(ctx context.Context, status string) ([]model.PaymentView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT p.id, p.booking_id, p.amount, p.method, p.status, p.paid_at, p.reference, p.created_at,
                   b.user_id, b.status AS booking_status,
                   r.name AS room_name, r.room_no,
                   a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end
            FROM payments p
            JOIN bookings b ON b.id = p.booking_id
            JOIN rooms r ON r.id = b.room_id
            JOIN availability a ON a.id = b.availability_id
            WHERE 1=1`
  args := []interface{}{}
  if status != "" {
    args = append(args, status)
    query += fmt.Sprintf(" AND p.status = $%d", len(args))
  }
  query += " ORDER BY p.created_at DESC"

  var rows []model.PaymentView
  err := r.DB.SelectContext(ctx, &rows, query, args...)
  return rows, err
}

func (r *PaymentRepo) UpdateStatus(ctx context.Context, id int64, status string, paidAt *time.Time) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  res, err := r.DB.ExecContext(ctx, `UPDATE payments SET status=$1, paid_at=$2 WHERE id=$3`, status, paidAt, id)
  if err != nil {
    return err
  }
  affected, _ := res.RowsAffected()
  if affected == 0 {
    return sql.ErrNoRows
  }
  return nil
}
