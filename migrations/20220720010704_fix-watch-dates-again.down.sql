UPDATE movie_watch
SET watched = DATE(watched, '-1 day')
WHERE (
    (watched >= '2020-11-02') AND
    (watched <= '2021-03-13')
) OR (
    (watched >= '2021-11-08') AND
    (watched <= '2022-03-12')
);