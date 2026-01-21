package service

import (
  "context"
  "database/sql"
  "time"

  "github.com/jackc/pgconn"
  "github.com/jmoiron/sqlx"
)

type PaymentCreateInput struct {
  BookingID int64
  UserID    int64
  Method    string
}

type PaymentService struct {
  DB      *sqlx.DB
  Timeout time.Duration
}

func NewPaymentService(db *sqlx.DB) *PaymentService {
  return &PaymentService{DB: db, Timeout: 5 * time.Second}
}

func (s *PaymentService) Create(ctx context.Context, input PaymentCreateInput) (int64, error) {
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

  var bookingUserID int64
  var bookingStatus string
  var pricePerSlot int64
  err = tx.QueryRowContext(ctx, `SELECT b.user_id, b.status, r.price_per_slot
     FROM bookings b
     JOIN rooms r ON r.id = b.room_id
     WHERE b.id=$1 FOR UPDATE`, input.BookingID).Scan(&bookingUserID, &bookingStatus, &pricePerSlot)
  if err == sql.ErrNoRows {
    return 0, ErrNotFound
  }
  if err != nil {
    return 0, err
  }
  if bookingUserID != input.UserID {
    return 0, ErrForbidden
  }
  if bookingStatus != "APPROVED" {
    return 0, ErrInvalid
  }

  var paymentID int64
  err = tx.QueryRowContext(ctx,
    `INSERT INTO payments (booking_id, amount, method, status) VALUES ($1,$2,$3,'UNPAID') RETURNING id`,
    input.BookingID, pricePerSlot, input.Method,
  ).Scan(&paymentID)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return 0, ErrConflict
    }
    return 0, err
  }

  if err := tx.Commit(); err != nil {
    return 0, err
  }
  return paymentID, nil
}

func (s *PaymentService) MarkPaid(ctx context.Context, paymentID int64, userID int64) error {
  if ctx == nil {
    ctx = context.Background()
  }
  ctx, cancel := context.WithTimeout(ctx, s.Timeout)
  defer cancel()

  res, err := s.DB.ExecContext(ctx, `UPDATE payments p
     SET status='PAID', paid_at=now()
     FROM bookings b
     WHERE p.id=$1 AND p.booking_id=b.id AND b.user_id=$2`, paymentID, userID)
  if err != nil {
    return err
  }
  affected, _ := res.RowsAffected()
  if affected == 0 {
    return sql.ErrNoRows
  }
  return nil
}
