package model

import "time"

type Booking struct {
  ID          int64     `db:"id" json:"id"`
  UserID      int64     `db:"user_id" json:"user_id"`
  RoomID      int64     `db:"room_id" json:"room_id"`
  GuestName   string    `db:"guest_name" json:"guest_name"`
  GuestPhone  string    `db:"guest_phone" json:"guest_phone"`
  Notes       string    `db:"notes" json:"notes"`
  Status      string    `db:"status" json:"status"`
  StartDate   string    `db:"start_date" json:"start_date"`
  EndDate     string    `db:"end_date" json:"end_date"`
  TotalDays   int       `db:"total_days" json:"total_days"`
  TotalAmount int64     `db:"total_amount" json:"total_amount"`
  CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type BookingView struct {
  ID          int64     `db:"id" json:"id"`
  UserID      int64     `db:"user_id" json:"user_id"`
  RoomID      int64     `db:"room_id" json:"room_id"`
  GuestName   string    `db:"guest_name" json:"guest_name"`
  GuestPhone  string    `db:"guest_phone" json:"guest_phone"`
  Notes       string    `db:"notes" json:"notes"`
  Status      string    `db:"status" json:"status"`
  StartDate   string    `db:"start_date" json:"start_date"`
  EndDate     string    `db:"end_date" json:"end_date"`
  TotalDays   int       `db:"total_days" json:"total_days"`
  TotalAmount int64     `db:"total_amount" json:"total_amount"`
  CreatedAt   time.Time `db:"created_at" json:"created_at"`
  RoomName    string    `db:"room_name" json:"room_name"`
  RoomNo      string    `db:"room_no" json:"room_no"`
  RoomType    string    `db:"room_type" json:"room_type"`
  RoomCapacity int      `db:"room_capacity" json:"room_capacity"`
  RoomPrice   int64     `db:"room_price" json:"room_price"`
  RoomImage   string    `db:"room_image" json:"room_image"`
}
