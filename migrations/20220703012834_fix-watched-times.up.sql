ALTER TABLE movie_watch ADD COLUMN watched_str TEXT;
UPDATE movie_watch 
    SET watched_str = DATE(watched + 5*60*60, 'unixepoch', 'localtime');
DROP INDEX IF EXISTS idx_movie_watch_imdb_id_watched;
ALTER TABLE movie_watch DROP COLUMN watched;
ALTER TABLE movie_watch RENAME COLUMN watched_str TO watched;
CREATE INDEX IF NOT EXISTS idx_movie_watch_imdb_id_watched ON movie_watch(imdb_id, watched);