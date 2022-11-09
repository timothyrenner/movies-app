-- name: GetReviewForMovie :one
SELECT uuid,
    movie_uuid,
    movie_title,
    review,
    liked
FROM review
WHERE movie_title = ?;
-- name: InsertReview :exec
INSERT INTO review (uuid, movie_uuid, movie_title, review, liked)
VALUES (?, ?, ?, ?, ?) ON CONFLICT (movie_uuid) DO
UPDATE
SET review = excluded.review,
    liked = excluded.liked;
-- name: UpdateMovieUuidForReview :exec
UPDATE review SET movie_uuid = ? WHERE uuid = ?;