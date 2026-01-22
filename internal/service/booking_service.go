package service

import (
  "context"
  "database/sql"
  "time"

  "github.com/jackc/pgconn"
  "github.com/jmoiron/sqlx"
)

type BookingCreateInput struct {
  UserID    int64
  RoomID    int64
  StartDate string
  EndDate   string
  GuestName string
  GuestPhone string
  Notes     string
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

  start, err := time.Parse("2006-01-02", input.StartDate)
  if err != nil {
    return 0, ErrInvalid
  }
  end, err := time.Parse("2006-01-02", input.EndDate)
  if err != nil {
    return 0, ErrInvalid
  }
  if end.Before(start) {
    return 0, ErrInvalid
  }
  days := int(end.Sub(start).Hours()/24) + 1
  if days <= 0 {
    return 0, ErrInvalid
  }

  var roomPrice int64
  var roomStatus string
  err = tx.QueryRowContext(ctx, `SELECT price_per_slot, status FROM rooms WHERE id=$1 FOR UPDATE`, input.RoomID).
    Scan(&roomPrice, &roomStatus)
  if err == sql.ErrNoRows {
    return 0, ErrNotFound
  }
  if err != nil {
    return 0, err
  }
  if roomStatus != "ACTIVE" {
    return 0, ErrInvalid
  }

  var overlap bool
  err = tx.QueryRowContext(ctx, `SELECT EXISTS (
      SELECT 1 FROM bookings
      WHERE room_id=$1
        AND status IN ('PENDING','APPROVED')
        AND NOT (end_date < $2 OR start_date > $3)
    )`, input.RoomID, input.StartDate, input.EndDate).Scan(&overlap)
  if err != nil {
    return 0, err
  }
  if overlap {
    return 0, ErrConflict
  }

  var bookingID int64
  totalAmount := roomPrice * int64(days)
  err = tx.QueryRowContext(ctx,
    `INSERT INTO bookings (user_id, room_id, guest_name, guest_phone, notes, status, start_date, end_date, total_days, total_amount)
     VALUES ($1,$2,$3,$4,$5,'PENDING',$6,$7,$8,$9) RETURNING id`,
    input.UserID, input.RoomID, input.GuestName, input.GuestPhone, input.Notes,
    input.StartDate, input.EndDate, days, totalAmount,
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
