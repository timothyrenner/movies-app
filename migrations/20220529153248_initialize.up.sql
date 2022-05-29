PRAGMA foreign_keys = ON;
-- Main table.
CREATE TABLE IF NOT EXISTS movie (
    uuid TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    year INTEGER NOT NULL,
    rated TEXT,
    released TEXT,
    runtime_minutes INTEGER,
    plot TEXT,
    country TEXT,
    box_office TEXT,
    production TEXT,
    call_felissa INTEGER NOT NULL,
    slasher INTEGER NOT NULL,
    zombies INTEGER NOT NULL,
    beast INTEGER NOT NULL,
    godzilla INTEGER NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH())
);
CREATE INDEX IF NOT EXISTS idx_movie_title ON movie(title);
CREATE INDEX IF NOT EXISTS idx_movie_year ON movie(year);
-- Watch ledger.
CREATE TABLE IF NOT EXISTS movie_watch (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    watched INTEGER NOT NULL,
    service TEXT NOT NULL,
    first_time INTEGER NOT NULL,
    joe_bob INTEGER NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_watch_movie_uuid ON movie_watch(movie_uuid);
CREATE INDEX IF NOT EXISTS idx_movie_watch_watched ON movie_watch(watched);
-- Auxiliary table for genre.
CREATE TABLE IF NOT EXISTS movie_genre (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_genre_movie_uuid ON movie_genre(movie_uuid);
-- Auxiliary table for actors.
CREATE TABLE IF NOT EXISTS movie_actor (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_actor_movie_uuid ON movie_actor(movie_uuid);
-- Auxiliary table for directors.
CREATE TABLE IF NOT EXISTS movie_director (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_director_movie_uuid ON movie_director(movie_uuid);
-- Auxiliary table for producers.
CREATE TABLE IF NOT EXISTS movie_producer (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    name TEXT NOT NULL,
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_producer_movie_uuid ON movie_producer(movie_uuid);
-- Auxiliary table for writers.
CREATE TABLE IF NOT EXISTS movie_writer (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_writer_movie_uuid ON movie_writer(movie_uuid);
-- Auxiliary table for movie ratings.
CREATE TABLE IF NOT EXISTS movie_rating (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT,
    source TEXT NOT NULL,
    value TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_rating_movie_uuid ON movie_rating(movie_uuid);