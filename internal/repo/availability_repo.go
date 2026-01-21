package repo

import (
  "context"
  "database/sql"
  "fmt"

  "github.com/jmoiron/sqlx"

  "project_ap3/internal/model"
)

type AvailabilityRepo struct {
  DB *sqlx.DB
}

func NewAvailabilityRepo(db *sqlx.DB) *AvailabilityRepo {
  return &AvailabilityRepo{DB: db}
}

func (r *AvailabilityRepo) Create(ctx context.Context, a *model.Availability) (int64, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var id int64
  err := r.DB.QueryRowContext(ctx,
    `INSERT INTO availability (room_id, date, time_start, time_end, is_open) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
    a.RoomID, a.Date, a.TimeStart, a.TimeEnd, a.IsOpen,
  ).Scan(&id)
  return id, err
}

func (r *AvailabilityRepo) Update(ctx context.Context, a *model.Availability) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  _, err := r.DB.ExecContext(ctx,
    `UPDATE availability SET room_id=$1, date=$2, time_start=$3, time_end=$4, is_open=$5 WHERE id=$6`,
    a.RoomID, a.Date, a.TimeStart, a.TimeEnd, a.IsOpen, a.ID,
  )
  return err
}

func (r *AvailabilityRepo) Delete(ctx context.Context, id int64) error {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  _, err := r.DB.ExecContext(ctx, `DELETE FROM availability WHERE id=$1`, id)
  return err
}

func (r *AvailabilityRepo) GetByID(ctx context.Context, id int64) (*model.Availability, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  var a model.Availability
  err := r.DB.GetContext(ctx, &a, `SELECT id, room_id, date::text AS date, to_char(time_start,'HH24:MI') AS time_start, to_char(time_end,'HH24:MI') AS time_end, is_open, created_at FROM availability WHERE id=$1`, id)
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  return &a, nil
}

func (r *AvailabilityRepo) List(ctx context.Context, date string, roomID *int64, onlyOpen bool) ([]model.AvailabilityView, error) {
  ctx, cancel := withTimeout(ctx)
  defer cancel()

  base := `SELECT a.id, a.room_id, a.date::text AS date, to_char(a.time_start,'HH24:MI') AS time_start, to_char(a.time_end,'HH24:MI') AS time_end,
                  a.is_open, a.created_at,
                  r.name AS room_name, r.room_no, r.type AS room_type, r.capacity AS room_capacity, r.price_per_slot AS room_price, r.status AS room_status
           FROM availability a
           JOIN rooms r ON r.id = a.room_id`

  query := base + " WHERE 1=1"
  args := []interface{}{}
  if onlyOpen {
    query += " AND a.is_open = TRUE"
  }
  if date != "" {
    args = append(args, date)
    query += fmt.Sprintf(" AND a.date = $%d", len(args))
  }
  if roomID != nil {
    args = append(args, *roomID)
    query += fmt.Sprintf(" AND a.room_id = $%d", len(args))
  }
  query += " ORDER BY a.date, a.time_start"

  var rows []model.AvailabilityView
  err := r.DB.SelectContext(ctx, &rows, query, args...)
  return rows, err
}
