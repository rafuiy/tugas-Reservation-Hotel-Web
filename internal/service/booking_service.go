package service

import (
  "context"
  "database/sql"
  "time"

  "github.com/jackc/pgconn"
  "github.com/jmoiron/sqlx"
)

type BookingCreateInput struct {
  UserID         int64
  AvailabilityID int64
  GuestName      string
  GuestPhone     string
  Notes          string
}

type BookingService struct {
  DB      *sqlx.DB
  Timeout time.Duration
}

func NewBookingService(db *sqlx.DB) *BookingService {
  return &BookingService{DB: db, Timeout: 5 * time.Second}
}

func (s *BookingService) Create(ctx context.Context, input BookingCreateInput) (int64, error) {
  if ctx == nil {
    ctx = context.Background()
  }
  ctx, cancel := context.WithTimeout(ctx, s.Timeout)
  defer cancel()

  tx, err := s.DB.BeginTxx(ctx, nil)
  if err != nil {
    return 0, err
  }
  defer func() {
    _ = tx.Rollback()
  }()

  var roomID int64
  var isOpen bool
  err = tx.QueryRowContext(ctx, `SELECT room_id, is_open FROM availability WHERE id=$1 FOR UPDATE`, input.AvailabilityID).Scan(&roomID, &isOpen)
  if err == sql.ErrNoRows {
    return 0, ErrNotFound
  }
  if err != nil {
    return 0, err
  }
  if !isOpen {
    return 0, ErrInvalid
  }

  var bookingID int64
  err = tx.QueryRowContext(ctx,
    `INSERT INTO bookings (user_id, room_id, availability_id, guest_name, guest_phone, notes, status)
     VALUES ($1,$2,$3,$4,$5,$6,'PENDING') RETURNING id`,
    input.UserID, roomID, input.AvailabilityID, input.GuestName, input.GuestPhone, input.Notes,
  ).Scan(&bookingID)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return 0, ErrConflict
    }
    return 0, err
  }

  if err := tx.Commit(); err != nil {
    return 0, err
  }
  return bookingID, nil
}
