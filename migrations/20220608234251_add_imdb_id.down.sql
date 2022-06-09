DROP INDEX IF EXISTS idx_movie_imdb_id;
ALTER TABLE movie DROP COLUMN imdb_id;

CREATE INDEX IF NOT EXISTS idx_movie_title ON movie(title);
CREATE INDEX IF NOT EXISTS idx_movie_year ON movie(year);

DROP INDEX IF EXISTS idx_movie_watch_imdb_id_watched;
ALTER TABLE movie_watch DROP COLUMN imdb_id;

CREATE INDEX IF NOT EXISTS idx_movie_watch_title_watched ON movie_watch(movie_title, watched);