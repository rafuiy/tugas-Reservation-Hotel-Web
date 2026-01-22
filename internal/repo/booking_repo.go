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
  err := r.DB.GetContext(ctx, &booking, `
    SELECT id, user_id, room_id, guest_name, guest_phone, notes, status,
           COALESCE(start_date::text, '') AS start_date, COALESCE(end_date::text, '') AS end_date,
           COALESCE(total_days, 0) AS total_days, COALESCE(total_amount, 0) AS total_amount,
           created_at
    FROM bookings
    WHERE id=$1
  `, id)
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

  query := `SELECT b.id, b.user_id, b.room_id, b.guest_name, b.guest_phone, b.notes, b.status,
                   COALESCE(b.start_date::text, '') AS start_date, COALESCE(b.end_date::text, '') AS end_date,
                   COALESCE(b.total_days, 0) AS total_days, COALESCE(b.total_amount, 0) AS total_amount,
                   b.created_at,
                   r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price,
                   COALESCE(r.image_url, '') AS room_image
            FROM bookings b
            JOIN rooms r ON r.id = b.room_id
            WHERE b.user_id=$1
            ORDER BY b.created_at DESC`

  var rows []model.BookingView
  err := r.DB.SelectContext(ctx, &rows, query, userID)
  return rows, err
}

func (r *BookingRepo) ListAdmin(ctx context.Context, status string) ([]model.BookingView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT b.id, b.user_id, b.room_id, b.guest_name, b.guest_phone, b.notes, b.status,
                   COALESCE(b.start_date::text, '') AS start_date, COALESCE(b.end_date::text, '') AS end_date,
                   COALESCE(b.total_days, 0) AS total_days, COALESCE(b.total_amount, 0) AS total_amount,
                   b.created_at,
                   r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price,
                   COALESCE(r.image_url, '') AS room_image
            FROM bookings b
            JOIN rooms r ON r.id = b.room_id
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

func (r *BookingRepo) ListSchedule(ctx context.Context, roomID *int64) ([]model.BookingView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  query := `SELECT b.id, b.user_id, b.room_id, b.guest_name, b.guest_phone, b.notes, b.status,
                   COALESCE(b.start_date::text, '') AS start_date, COALESCE(b.end_date::text, '') AS end_date,
                   COALESCE(b.total_days, 0) AS total_days, COALESCE(b.total_amount, 0) AS total_amount,
                   b.created_at,
                   r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price,
                   COALESCE(r.image_url, '') AS room_image
            FROM bookings b
            JOIN rooms r ON r.id = b.room_id
            WHERE b.status IN ('PENDING','APPROVED')`
  args := []interface{}{}
  if roomID != nil {
    args = append(args, *roomID)
    query += fmt.Sprintf(" AND b.room_id = $%d", len(args))
  }
  query += " ORDER BY b.start_date, b.end_date"

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
