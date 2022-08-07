ALTER TABLE movie
ADD COLUMN runtime_minutes_not_nullable INTEGER NOT NULL DEFAULT 0;
UPDATE movie
SET runtime_minutes_not_nullable = CASE
        WHEN runtime_minutes IS NULL THEN 0
        ELSE runtime_minutes
    END;
ALTER TABLE movie DROP COLUMN runtime_minutes;
ALTER TABLE movie
    RENAME COLUMN runtime_minutes_not_nullable TO runtime_minutes;