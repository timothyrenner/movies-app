ALTER TABLE movie_watch ADD COLUMN watched_int INT;
UPDATE movie_watch
    SET watched_int = UNIXEPOCH(watched);
DROP INDEX IF EXISTS idx_movie_watch_imdb_id_watched;
ALTER TABLE movie_watch DROP COLUMN watched;
ALTER TABLE movie_watch RENAME COLUMN watched_int TO watched;
CREATE INDEX IF NOT EXISTS idx_movie_watch_imdb_id_watched ON movie_watch(imdb_id, watched);