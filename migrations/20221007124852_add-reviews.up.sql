CREATE TABLE IF NOT EXISTS review (
    uuid TEXT PRIMARY KEY NOT NULL,
    movie_uuid TEXT UNIQUE NOT NULL,
    movie_title TEXT UNIQUE NOT NULL,
    review TEXT NOT NULL,
    liked INTEGER NOT NULL DEFAULT 0,
    created_datetime INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    FOREIGN KEY (movie_uuid) REFERENCES movie(uuid)
);
CREATE INDEX IF NOT EXISTS idx_review_movie_uuid ON review(movie_uuid);
CREATE INDEX IF NOT EXISTS idx_review_movie_title ON review(movie_title);
