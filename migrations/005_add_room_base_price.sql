ALTER TABLE rooms
  ADD COLUMN IF NOT EXISTS base_price BIGINT;

UPDATE rooms
SET base_price = price_per_slot
WHERE base_price IS NULL;
