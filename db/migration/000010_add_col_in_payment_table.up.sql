BEGIN;

-- Add `transaction_id` if it does not exist
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'payments' AND column_name = 'transaction_id'
  ) THEN
ALTER TABLE payments
    ADD COLUMN transaction_id TEXT UNIQUE;
END IF;
END$$;

COMMIT;
