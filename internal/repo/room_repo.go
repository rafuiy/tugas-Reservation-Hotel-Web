package repo

import (
  "context"
  "database/sql"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type RoomRepo struct {
  DB *sqlx.DB
}

func NewRoomRepo(db *sqlx.DB) *RoomRepo {
  return &RoomRepo{DB: db}
}

func (r *RoomRepo) Create(ctx context.Context, room *model.Room) (int64, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var id int64
  err := r.DB.QueryRowContext(ctx,
    `INSERT INTO rooms (room_no, name, type, capacity, price_per_slot, status) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
    room.RoomNo, room.Name, room.Type, room.Capacity, room.PricePerSlot, room.Status,
  ).Scan(&id)
  return id, err
}

func (r *RoomRepo) Update(ctx context.Context, room *model.Room) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  _, err := r.DB.ExecContext(ctx,
    `UPDATE rooms SET room_no=$1, name=$2, type=$3, capacity=$4, price_per_slot=$5, status=$6 WHERE id=$7`,
    room.RoomNo, room.Name, room.Type, room.Capacity, room.PricePerSlot, room.Status, room.ID,
  )
  return err
}

func (r *RoomRepo) Delete(ctx context.Context, id int64) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  _, err := r.DB.ExecContext(ctx, `DELETE FROM rooms WHERE id=$1`, id)
  return err
}

func (r *RoomRepo) GetByID(ctx context.Context, id int64) (*model.Room, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var room model.Room
  err := r.DB.GetContext(ctx, &room, `SELECT * FROM rooms WHERE id=$1`, id)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &room, nil
}

func (r *RoomRepo) ListAll(ctx context.Context) ([]model.Room, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var rooms []model.Room
  err := r.DB.SelectContext(ctx, &rooms, `SELECT * FROM rooms ORDER BY id DESC`)
  return rooms, err
}

func (r *RoomRepo) ListActive(ctx context.Context) ([]model.Room, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var rooms []model.Room
  err := r.DB.SelectContext(ctx, &rooms, `SELECT * FROM rooms WHERE status='ACTIVE' ORDER BY room_no`)
  return rooms, err
}
