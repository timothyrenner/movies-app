-- Drop primary table.
DROP TABLE IF EXISTS movie;
DROP INDEX IF EXISTS idx_movie_title;
DROP INDEX IF EXISTS idx_movie_year;
-- Drop primary ledger.
DROP TABLE IF EXISTS movie_watch;
DROP INDEX IF EXISTS idx_movie_watch_movie_uuid;
DROP INDEX IF EXISTS idx_movie_watch_watched;
-- Drop auxiliary table for genres.
DROP TABLE IF EXISTS movie_genre;
DROP INDEX IF EXISTS idx_movie_genre_movie_uuid;
-- Drop auxiliary table for actors.
DROP TABLE IF EXISTS movie_actor;
DROP INDEX IF EXISTS idx_movie_actor_movie_uuid;
-- Drop auxiliary table for directors.
DROP TABLE IF EXISTS movie_director;
DROP INDEX IF EXISTS idx_movie_director_movie_uuid;
-- Drop auxiliary table for producers.
DROP TABLE IF EXISTS movie_producer;
DROP INDEX IF EXISTS idx_movie_producer_movie_uuid;
-- Drop auxiliary table for writers.
DROP TABLE IF EXISTS movie_writer;
DROP INDEX IF EXISTS idx_movie_writer_movie_uuid;
-- Drop auxiliary table for ratings.
DROP TABLE IF EXISTS movie_rating;
DROP INDEX IF EXISTS idx_movie_rating_movie_uuid;