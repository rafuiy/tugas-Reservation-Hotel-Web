package repo

import (
  "context"
  "database/sql"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type InvoiceRepo struct {
  DB *sqlx.DB
}

func NewInvoiceRepo(db *sqlx.DB) *InvoiceRepo {
  return &InvoiceRepo{DB: db}
}

func (r *InvoiceRepo) GetByBookingID(ctx context.Context, bookingID int64) (*model.Invoice, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  type row struct {
    BookingID     int64          `db:"booking_id"`
    UserID        int64          `db:"user_id"`
    BookingStatus string         `db:"booking_status"`
    GuestName     string         `db:"guest_name"`
    GuestPhone    string         `db:"guest_phone"`
    Notes         string         `db:"notes"`
    RoomName      string         `db:"room_name"`
    RoomNo        string         `db:"room_no"`
    RoomType      string         `db:"room_type"`
    RoomCapacity  int            `db:"room_capacity"`
    PricePerSlot  int64          `db:"price_per_slot"`
    Date          string         `db:"date"`
    TimeStart     string         `db:"time_start"`
    TimeEnd       string         `db:"time_end"`
    PaymentID     sql.NullInt64  `db:"payment_id"`
    PaymentStatus sql.NullString `db:"payment_status"`
    PaymentMethod sql.NullString `db:"payment_method"`
    PaymentAmount sql.NullInt64  `db:"payment_amount"`
    PaidAt        sql.NullTime   `db:"paid_at"`
  }

  var rrow row
  err := r.DB.GetContext(ctx, &rrow, `SELECT b.id AS booking_id, b.user_id, b.status AS booking_status,
           b.guest_name, b.guest_phone, b.notes,
           r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot,
           a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end,
           p.id AS payment_id, p.status AS payment_status, p.method AS payment_method, p.amount AS payment_amount, p.paid_at
    FROM bookings b
    JOIN rooms r ON r.id = b.room_id
    JOIN availability a ON a.id = b.availability_id
    LEFT JOIN payments p ON p.booking_id = b.id
    WHERE b.id = $1`, bookingID)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }

  inv := model.Invoice{
    BookingID:     rrow.BookingID,
    UserID:        rrow.UserID,
    BookingStatus: rrow.BookingStatus,
    GuestName:     rrow.GuestName,
    GuestPhone:    rrow.GuestPhone,
    Notes:         rrow.Notes,
    RoomName:      rrow.RoomName,
    RoomNo:        rrow.RoomNo,
    RoomType:      rrow.RoomType,
    RoomCapacity:  rrow.RoomCapacity,
    PricePerSlot:  rrow.PricePerSlot,
    Date:          rrow.Date,
    TimeStart:     rrow.TimeStart,
    TimeEnd:       rrow.TimeEnd,
  }
  if rrow.PaymentID.Valid {
    inv.PaymentID = &rrow.PaymentID.Int64
  }
  if rrow.PaymentStatus.Valid {
    inv.PaymentStatus = &rrow.PaymentStatus.String
  }
  if rrow.PaymentMethod.Valid {
    inv.PaymentMethod = &rrow.PaymentMethod.String
  }
  if rrow.PaymentAmount.Valid {
    inv.PaymentAmount = &rrow.PaymentAmount.Int64
  }
  if rrow.PaidAt.Valid {
    inv.PaidAt = &rrow.PaidAt.Time
  }
  return &inv, nil
}
