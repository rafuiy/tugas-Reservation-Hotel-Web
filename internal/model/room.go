package model

import "time"

type Room struct {
  ID           int64     `db:"id" json:"id"`
  RoomNo       string    `db:"room_no" json:"room_no"`
  Name         string    `db:"name" json:"name"`
  Type         string    `db:"type" json:"type"`
  Capacity     int       `db:"capacity" json:"capacity"`
  PricePerSlot int64     `db:"price_per_slot" json:"price_per_slot"`
  Status       string    `db:"status" json:"status"`
  CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
