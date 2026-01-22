ALTER TABLE bookings
  DROP CONSTRAINT IF EXISTS booking_unique_availability;

ALTER TABLE bookings
  DROP CONSTRAINT IF EXISTS bookings_availability_id_fkey;

ALTER TABLE bookings
  DROP COLUMN IF EXISTS availability_id;

ALTER TABLE bookings
  ADD COLUMN IF NOT EXISTS start_date DATE,
  ADD COLUMN IF NOT EXISTS end_date DATE,
  ADD COLUMN IF NOT EXISTS total_days INTEGER,
  ADD COLUMN IF NOT EXISTS total_amount BIGINT;
