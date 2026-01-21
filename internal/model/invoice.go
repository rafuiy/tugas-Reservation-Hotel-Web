package model

import "time"

type Invoice struct {
  BookingID     int64      `json:"booking_id"`
  UserID        int64      `json:"user_id"`
  BookingStatus string     `json:"booking_status"`
  GuestName     string     `json:"guest_name"`
  GuestPhone    string     `json:"guest_phone"`
  Notes         string     `json:"notes"`
  RoomName      string     `json:"room_name"`
  RoomNo        string     `json:"room_no"`
  RoomType      string     `json:"room_type"`
  RoomCapacity  int        `json:"room_capacity"`
  PricePerSlot  int64      `json:"price_per_slot"`
  Date          string     `json:"date"`
  TimeStart     string     `json:"time_start"`
  TimeEnd       string     `json:"time_end"`
  PaymentID     *int64     `json:"payment_id"`
  PaymentStatus *string    `json:"payment_status"`
  PaymentMethod *string    `json:"payment_method"`
  PaymentAmount *int64     `json:"payment_amount"`
  PaidAt        *time.Time `json:"paid_at"`
}
