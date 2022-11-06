PRAGMA foreign_keys = OFF;
-- make movie_uuid in movie_watch not null.
DROP INDEX idx_movie_watch_movie_uuid;
DROP INDEX idx_movie_watch_imdb_id_watched;
ALTER TABLE movie_watch RENAME TO movie_watch_old;
CREATE TABLE movie_watch (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    movie_title TEXT NOT NULL,
    service TEXT NOT NULL,
    first_time INTEGER NOT NULL,
    joe_bob INTEGER NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    imdb_id TEXT NOT NULL DEFAULT '',
    watched TEXT NOT NULL,
    notes TEXT,
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_watch_movie_uuid ON movie_watch(movie_uuid);
CREATE INDEX idx_movie_watch_imdb_id_watched ON movie_watch(imdb_id, watched);
INSERT INTO movie_watch SELECT * FROM movie_watch_old;
DROP TABLE movie_watch_old;

-- make movie_uuid in movie_actor not null.
DROP INDEX idx_movie_actor_movie_uuid;
ALTER TABLE movie_actor RENAME TO movie_actor_old;
CREATE TABLE movie_actor (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_actor_movie_uuid ON movie_actor(movie_uuid);
INSERT INTO movie_actor SELECT * FROM movie_actor_old;
DROP TABLE movie_actor_old;

-- make movie_uuid in movie_director not null.
DROP INDEX idx_movie_director_movie_uuid;
ALTER TABLE movie_director RENAME TO movie_director_old;
CREATE TABLE movie_director (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_director_movie_uuid ON movie_director(movie_uuid);
INSERT INTO movie_director SELECT * FROM movie_director_old;
DROP TABLE movie_director_old;

-- make movie_uuid in movie_genre not null.
DROP INDEX idx_movie_genre_movie_uuid;
ALTER TABLE movie_genre RENAME TO movie_genre_old;
CREATE TABLE movie_genre (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_genre_movie_uuid ON movie_genre(movie_uuid);
INSERT INTO movie_genre SELECT * FROM movie_genre_old;
DROP TABLE movie_genre_old;

-- make movie_uuid in movie_rating not null.
DROP INDEX idx_movie_rating_movie_uuid;
ALTER TABLE movie_rating RENAME TO movie_rating_old;
CREATE TABLE movie_rating (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    source TEXT NOT NULL,
    value TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_rating_movie_uuid ON movie_rating(movie_uuid);
INSERT INTO movie_rating SELECT * FROM movie_rating_old;
DROP TABLE movie_rating_old;

-- make movie_uuid in movie_writer not null.
DROP INDEX idx_movie_writer_movie_uuid;
ALTER TABLE movie_writer RENAME TO movie_writer_old;
CREATE TABLE movie_writer (
    uuid TEXT PRIMARY KEY,
    movie_uuid TEXT NOT NULL,
    name TEXT NOT NULL,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY(movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX idx_movie_writer_movie_uuid ON movie_writer(movie_uuid);
INSERT INTO movie_writer SELECT * FROM movie_writer_old;
DROP TABLE movie_writer_old;
PRAGMA foreign_keys = ON;