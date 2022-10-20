-- name: FindMovieWatch :one
SELECT uuid
FROM movie_watch
WHERE imdb_id = ?
    AND watched = ?;
-- name: GetAllMovieWatches :many
SELECT w.uuid,
    w.movie_uuid,
    w.movie_title,
    w.imdb_id,
    w.watched,
    w.service,
    w.first_time,
    w.joe_bob,
    w.notes,
    m.imdb_link,
    m.slasher,
    m.call_felissa,
    m.beast,
    m.godzilla,
    m.zombies
FROM movie_watch AS w
    INNER JOIN movie AS m ON m.uuid = w.movie_uuid;
-- name: FindMovie :one
SELECT uuid
FROM movie
WHERE imdb_id = ?;
-- name: GetMovie :one
SELECT uuid,
    title,
    imdb_link,
    imdb_id,
    year,
    rated,
    released,
    runtime_minutes,
    plot,
    country,
    language,
    box_office,
    production,
    call_felissa,
    slasher,
    zombies,
    beast,
    godzilla
FROM movie
WHERE uuid = ?;
-- name: InsertMovie :exec
INSERT INTO movie (
        uuid,
        title,
        imdb_link,
        imdb_id,
        year,
        rated,
        released,
        runtime_minutes,
        plot,
        country,
        language,
        box_office,
        production,
        call_felissa,
        slasher,
        zombies,
        beast,
        godzilla,
        wallpaper_fu
    )
VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    );
-- name: InsertMovieGenre :exec
INSERT INTO movie_genre (uuid, movie_uuid, name)
VALUES (?, ?, ?);
-- name: InsertMovieActor :exec
INSERT INTO movie_actor (uuid, movie_uuid, name)
VALUES (?, ?, ?);
-- name: InsertMovieDirector :exec
INSERT INTO movie_director (uuid, movie_uuid, name)
VALUES (?, ?, ?);
-- name: InsertMovieWriter :exec
INSERT INTO movie_writer (uuid, movie_uuid, name)
VALUES (?, ?, ?);
-- name: InsertMovieRating :exec
INSERT INTO movie_rating (
        uuid,
        movie_uuid,
        source,
        value
    )
VALUES (?, ?, ?, ?);
-- name: InsertMovieWatch :exec
INSERT INTO movie_watch (
        uuid,
        movie_uuid,
        movie_title,
        imdb_id,
        watched,
        service,
        first_time,
        joe_bob,
        notes
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
-- name: GetGenreNamesForMovie :many
SELECT name
FROM movie_genre
WHERE movie_uuid = ?;
-- name: GetActorNamesForMovie :many
SELECT name
FROM movie_actor
WHERE movie_uuid = ?;
-- name: GetDirectorNamesForMovie :many
SELECT name
FROM movie_director
WHERE movie_uuid = ?;
-- name: GetWriterNamesForMovie :many
SELECT name
FROM movie_writer
WHERE movie_uuid = ?;
-- name: GetRatingsForMovie :many
SELECT uuid,
    movie_uuid,
    source,
    value
FROM movie_rating
WHERE movie_uuid = ?;