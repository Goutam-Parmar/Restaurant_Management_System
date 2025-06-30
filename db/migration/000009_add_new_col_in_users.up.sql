BEGIN;

-- ✅ Add `phone` column if not exists
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'phone'
  ) THEN
ALTER TABLE users
    ADD COLUMN phone VARCHAR(10) UNIQUE CHECK (phone ~ '^[0-9]{10}$');
END IF;
END$$;

-- ✅ Add `updated_at` column if not exists
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'updated_at'
  ) THEN
ALTER TABLE users
    ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();
END IF;
END$$;

-- ✅ Add `deleted_at` column if not exists
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'deleted_at'
  ) THEN
ALTER TABLE users
    ADD COLUMN deleted_at TIMESTAMP;
END IF;
END$$;

COMMIT;
