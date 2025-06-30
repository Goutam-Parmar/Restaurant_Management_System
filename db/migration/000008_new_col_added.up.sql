BEGIN;

-- Add `created_by` column if not exists
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'created_by'
  ) THEN
ALTER TABLE users
    ADD COLUMN created_by BIGINT REFERENCES users(id) ON DELETE SET NULL;
END IF;
END$$;

-- Add `created_at` column if not exists
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'created_at'
  ) THEN
ALTER TABLE users
    ADD COLUMN created_at TIMESTAMP DEFAULT NOW();
END IF;
END$$;

COMMIT;

