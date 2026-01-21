CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  phone TEXT,
  role TEXT NOT NULL CHECK (role IN ('ADMIN','USER')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rooms (
  id BIGSERIAL PRIMARY KEY,
  room_no TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  capacity INTEGER NOT NULL CHECK (capacity > 0),
  price_per_slot BIGINT NOT NULL CHECK (price_per_slot >= 0),
  status TEXT NOT NULL CHECK (status IN ('ACTIVE','INACTIVE')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS availability (
  id BIGSERIAL PRIMARY KEY,
  room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  time_start TIME NOT NULL,
  time_end TIME NOT NULL,
  is_open BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT availability_unique UNIQUE (room_id, date, time_start, time_end),
  CONSTRAINT availability_time_check CHECK (time_start < time_end)
);
CREATE INDEX IF NOT EXISTS idx_availability_room_date ON availability(room_id, date);

CREATE TABLE IF NOT EXISTS bookings (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  availability_id BIGINT NOT NULL REFERENCES availability(id) ON DELETE CASCADE,
  guest_name TEXT NOT NULL,
  guest_phone TEXT NOT NULL,
  notes TEXT,
  status TEXT NOT NULL CHECK (status IN ('PENDING','APPROVED','REJECTED','CANCELLED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT booking_unique_availability UNIQUE (availability_id)
);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);

CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  booking_id BIGINT NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
  amount BIGINT NOT NULL CHECK (amount >= 0),
  method TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('UNPAID','PAID','FAILED')),
  paid_at TIMESTAMPTZ,
  reference TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

