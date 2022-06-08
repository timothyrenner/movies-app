ALTER TABLE movie DROP COLUMN imdb_id;
DROP INDEX IF EXISTS idx_movie_imdb_id;

CREATE INDEX IF NOT EXISTS idx_movie_title ON movie(title);
CREATE INDEX IF NOT EXISTS idx_movie_year ON movie(year);