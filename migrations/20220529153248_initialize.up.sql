PRAGMA foreign_keys = ON;
-- Main table.
CREATE TABLE IF NOT EXISTS movie (
    uuid TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    imdb_link TEXT NOT NULL,
    year INTEGER NOT NULL,
    rated TEXT,
    released TEXT,
    runtime_minutes INTEGER NOT NULL,
    plot TEXT,
    country TEXT,
    language TEXT,
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
    movie_title TEXT,
    watched INTEGER NOT NULL,
    service TEXT NOT NULL,
    first_time INTEGER NOT NULL,
    joe_bob INTEGER NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_movie_watch_movie_uuid ON movie_watch(movie_uuid);
CREATE INDEX IF NOT EXISTS idx_movie_watch_title_watched ON movie_watch(movie_title, watched);
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
-- Table mapping uuids to grist IDs.
CREATE TABLE IF NOT EXISTS uuid_grist (
    uuid TEXT PRIMARY KEY,
    grist_id INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_uuid_grist_grist_id ON uuid_grist(grist_id);