ALTER TABLE movie
ADD COLUMN runtime_minutes_nullable INTEGER;
UPDATE movie
SET runtime_minutes_nullable = CASE
        WHEN runtime_minutes = 0 THEN NULL
        ELSE runtime_minutes
    END;
ALTER TABLE movie DROP COLUMN runtime_minutes;
ALTER TABLE movie
    RENAME COLUMN runtime_minutes_nullable TO runtime_minutes;