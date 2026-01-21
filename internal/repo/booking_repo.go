package repo

import (
  "context"
  "database/sql"
  "fmt"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type BookingRepo struct {
  DB *sqlx.DB
}

func NewBookingRepo(db *sqlx.DB) *BookingRepo {
  return &BookingRepo{DB: db}
}

func (r *BookingRepo) GetByID(ctx context.Context, id int64) (*model.Booking, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var booking model.Booking
  err := r.DB.GetContext(ctx, &booking, `SELECT * FROM bookings WHERE id=$1`, id)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &booking, nil
}

func (r *BookingRepo) ListByUser(ctx context.Context, userID int64) ([]model.BookingView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT b.id, b.user_id, b.room_id, b.availability_id, b.guest_name, b.guest_phone, b.notes, b.status, b.created_at,
                   r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price,
                   a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end
            FROM bookings b
            JOIN rooms r ON r.id = b.room_id
            JOIN availability a ON a.id = b.availability_id
            WHERE b.user_id=$1
            ORDER BY b.created_at DESC`

  var rows []model.BookingView
  err := r.DB.SelectContext(ctx, &rows, query, userID)
  return rows, err
}

func (r *BookingRepo) ListAdmin(ctx context.Context, status string) ([]model.BookingView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT b.id, b.user_id, b.room_id, b.availability_id, b.guest_name, b.guest_phone, b.notes, b.status, b.created_at,
                   r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price,
                   a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end
            FROM bookings b
            JOIN rooms r ON r.id = b.room_id
            JOIN availability a ON a.id = b.availability_id
            WHERE 1=1`
  args := []interface{}{}
  if status != "" {
    args = append(args, status)
    query += fmt.Sprintf(" AND b.status = $%d", len(args))
  }
  query += " ORDER BY b.created_at DESC"

  var rows []model.BookingView
  err := r.DB.SelectContext(ctx, &rows, query, args...)
  return rows, err
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  res, err := r.DB.ExecContext(ctx, `UPDATE bookings SET status=$1 WHERE id=$2`, status, id)
  if err != nil {
    return err
  }
  affected, _ := res.RowsAffected()
  if affected == 0 {
    return sql.ErrNoRows
  }
  return nil
}
