ALTER TABLE movie ADD COLUMN imdb_id TEXT NOT NULL DEFAULT '';
CREATE INDEX idx_movie_imdb_id ON movie(imdb_id);
DROP INDEX IF EXISTS idx_movie_title;
DROP INDEX IF EXISTS idx_movie_year;

ALTER TABLE movie_watch ADD COLUMN imdb_id TEXT NOT NULL DEFAULT '';
CREATE INDEX idx_movie_watch_imdb_id_watched ON movie_watch(imdb_id, watched);
DROP INDEX IF EXISTS idx_movie_watch_title_watched;