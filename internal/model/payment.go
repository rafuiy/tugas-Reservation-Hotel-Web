package model

import "time"

type Payment struct {
  ID        int64      `db:"id" json:"id"`
  BookingID int64      `db:"booking_id" json:"booking_id"`
  Amount    int64      `db:"amount" json:"amount"`
  Method    string     `db:"method" json:"method"`
  Status    string     `db:"status" json:"status"`
  PaidAt    *time.Time `db:"paid_at" json:"paid_at"`
  Reference *string    `db:"reference" json:"reference,omitempty"`
  CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type PaymentView struct {
  ID            int64      `db:"id" json:"id"`
  BookingID     int64      `db:"booking_id" json:"booking_id"`
  Amount        int64      `db:"amount" json:"amount"`
  Method        string     `db:"method" json:"method"`
  Status        string     `db:"status" json:"status"`
  PaidAt        *time.Time `db:"paid_at" json:"paid_at"`
  Reference     *string    `db:"reference" json:"reference,omitempty"`
  CreatedAt     time.Time  `db:"created_at" json:"created_at"`
  UserID        int64      `db:"user_id" json:"user_id"`
  BookingStatus string     `db:"booking_status" json:"booking_status"`
  RoomName      string     `db:"room_name" json:"room_name"`
  RoomNo        string     `db:"room_no" json:"room_no"`
  StartDate     string     `db:"start_date" json:"start_date"`
  EndDate       string     `db:"end_date" json:"end_date"`
  TotalDays     int        `db:"total_days" json:"total_days"`
  TotalAmount   int64      `db:"total_amount" json:"total_amount"`
}
