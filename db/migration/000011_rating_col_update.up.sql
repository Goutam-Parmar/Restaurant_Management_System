BEGIN;

-- drop old constraint if exists
ALTER TABLE restaurants
DROP CONSTRAINT IF EXISTS restaurants_rating_check;

-- 2alter the column type from SMALLINT to DOUBLE PRECISION
ALTER TABLE restaurants
ALTER COLUMN rating TYPE DOUBLE PRECISION;

-- 3️add new CHECK constraint for 0.0 to 5.0
ALTER TABLE restaurants
    ADD CONSTRAINT restaurants_rating_check CHECK (rating >= 0.0 AND rating <= 5.0);

COMMIT;
