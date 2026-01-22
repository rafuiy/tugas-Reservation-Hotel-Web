package model

import "time"

type Availability struct {
  ID        int64     `db:"id" json:"id"`
  RoomID    int64     `db:"room_id" json:"room_id"`
  Date      string    `db:"date" json:"date"`
  TimeStart string    `db:"time_start" json:"time_start"`
  TimeEnd   string    `db:"time_end" json:"time_end"`
  IsOpen    bool      `db:"is_open" json:"is_open"`
  CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AvailabilityView struct {
  ID           int64     `db:"id" json:"id"`
  RoomID       int64     `db:"room_id" json:"room_id"`
  Date         string    `db:"date" json:"date"`
  TimeStart    string    `db:"time_start" json:"time_start"`
  TimeEnd      string    `db:"time_end" json:"time_end"`
  IsOpen       bool      `db:"is_open" json:"is_open"`
  IsBooked     bool      `db:"is_booked" json:"is_booked"`
  BookingStatus string   `db:"booking_status" json:"booking_status"`
  CreatedAt    time.Time `db:"created_at" json:"created_at"`
  RoomName     string    `db:"room_name" json:"room_name"`
  RoomNo       string    `db:"room_no" json:"room_no"`
  RoomType     string    `db:"room_type" json:"room_type"`
  RoomCapacity int       `db:"room_capacity" json:"room_capacity"`
  RoomPrice    int64     `db:"room_price" json:"room_price"`
  RoomStatus   string    `db:"room_status" json:"room_status"`
}
