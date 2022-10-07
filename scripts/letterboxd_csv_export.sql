.headers on
.mode csv
.once ./data/letterboxd_export.csv
SELECT
    movie.imdb_id AS imdbID,
    watch.watched AS WatchedDate,
    NOT watch.first_time AS Rewatch
FROM movie_watch AS watch
INNER JOIN movie ON
watch.movie_uuid = movie.uuid
ORDER BY watch.watched DESC;