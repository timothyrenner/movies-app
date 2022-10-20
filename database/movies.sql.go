// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: movies.sql

package database

import (
	"context"
	"database/sql"
)

const findMovie = `-- name: FindMovie :one
SELECT uuid
FROM movie
WHERE imdb_id = ?
`

func (q *Queries) FindMovie(ctx context.Context, imdbID string) (string, error) {
	row := q.db.QueryRowContext(ctx, findMovie, imdbID)
	var uuid string
	err := row.Scan(&uuid)
	return uuid, err
}

const findMovieWatch = `-- name: FindMovieWatch :one
SELECT uuid
FROM movie_watch
WHERE imdb_id = ?
    AND watched = ?
`

type FindMovieWatchParams struct {
	ImdbID  string
	Watched sql.NullString
}

func (q *Queries) FindMovieWatch(ctx context.Context, arg FindMovieWatchParams) (string, error) {
	row := q.db.QueryRowContext(ctx, findMovieWatch, arg.ImdbID, arg.Watched)
	var uuid string
	err := row.Scan(&uuid)
	return uuid, err
}

const getActorNamesForMovie = `-- name: GetActorNamesForMovie :many
SELECT name
FROM movie_actor
WHERE movie_uuid = ?
`

func (q *Queries) GetActorNamesForMovie(ctx context.Context, movieUuid sql.NullString) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getActorNamesForMovie, movieUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllMovieWatches = `-- name: GetAllMovieWatches :many
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
    INNER JOIN movie AS m ON m.uuid = w.movie_uuid
`

type GetAllMovieWatchesRow struct {
	Uuid        string
	MovieUuid   sql.NullString
	MovieTitle  sql.NullString
	ImdbID      string
	Watched     sql.NullString
	Service     string
	FirstTime   int64
	JoeBob      int64
	Notes       sql.NullString
	ImdbLink    string
	Slasher     int64
	CallFelissa int64
	Beast       int64
	Godzilla    int64
	Zombies     int64
}

func (q *Queries) GetAllMovieWatches(ctx context.Context) ([]GetAllMovieWatchesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMovieWatches)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMovieWatchesRow
	for rows.Next() {
		var i GetAllMovieWatchesRow
		if err := rows.Scan(
			&i.Uuid,
			&i.MovieUuid,
			&i.MovieTitle,
			&i.ImdbID,
			&i.Watched,
			&i.Service,
			&i.FirstTime,
			&i.JoeBob,
			&i.Notes,
			&i.ImdbLink,
			&i.Slasher,
			&i.CallFelissa,
			&i.Beast,
			&i.Godzilla,
			&i.Zombies,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDirectorNamesForMovie = `-- name: GetDirectorNamesForMovie :many
SELECT name
FROM movie_director
WHERE movie_uuid = ?
`

func (q *Queries) GetDirectorNamesForMovie(ctx context.Context, movieUuid sql.NullString) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getDirectorNamesForMovie, movieUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGenreNamesForMovie = `-- name: GetGenreNamesForMovie :many
SELECT name
FROM movie_genre
WHERE movie_uuid = ?
`

func (q *Queries) GetGenreNamesForMovie(ctx context.Context, movieUuid sql.NullString) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getGenreNamesForMovie, movieUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMovie = `-- name: GetMovie :one
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
WHERE uuid = ?
`

type GetMovieRow struct {
	Uuid           string
	Title          string
	ImdbLink       string
	ImdbID         string
	Year           int64
	Rated          sql.NullString
	Released       sql.NullString
	RuntimeMinutes sql.NullInt64
	Plot           sql.NullString
	Country        sql.NullString
	Language       sql.NullString
	BoxOffice      sql.NullString
	Production     sql.NullString
	CallFelissa    int64
	Slasher        int64
	Zombies        int64
	Beast          int64
	Godzilla       int64
}

func (q *Queries) GetMovie(ctx context.Context, uuid string) (GetMovieRow, error) {
	row := q.db.QueryRowContext(ctx, getMovie, uuid)
	var i GetMovieRow
	err := row.Scan(
		&i.Uuid,
		&i.Title,
		&i.ImdbLink,
		&i.ImdbID,
		&i.Year,
		&i.Rated,
		&i.Released,
		&i.RuntimeMinutes,
		&i.Plot,
		&i.Country,
		&i.Language,
		&i.BoxOffice,
		&i.Production,
		&i.CallFelissa,
		&i.Slasher,
		&i.Zombies,
		&i.Beast,
		&i.Godzilla,
	)
	return i, err
}

const getRatingsForMovie = `-- name: GetRatingsForMovie :many
SELECT uuid,
    movie_uuid,
    source,
    value
FROM movie_rating
WHERE movie_uuid = ?
`

type GetRatingsForMovieRow struct {
	Uuid      string
	MovieUuid sql.NullString
	Source    string
	Value     string
}

func (q *Queries) GetRatingsForMovie(ctx context.Context, movieUuid sql.NullString) ([]GetRatingsForMovieRow, error) {
	rows, err := q.db.QueryContext(ctx, getRatingsForMovie, movieUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRatingsForMovieRow
	for rows.Next() {
		var i GetRatingsForMovieRow
		if err := rows.Scan(
			&i.Uuid,
			&i.MovieUuid,
			&i.Source,
			&i.Value,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWriterNamesForMovie = `-- name: GetWriterNamesForMovie :many
SELECT name
FROM movie_writer
WHERE movie_uuid = ?
`

func (q *Queries) GetWriterNamesForMovie(ctx context.Context, movieUuid sql.NullString) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getWriterNamesForMovie, movieUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertMovie = `-- name: InsertMovie :exec
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
    )
`

type InsertMovieParams struct {
	Uuid           string
	Title          string
	ImdbLink       string
	ImdbID         string
	Year           int64
	Rated          sql.NullString
	Released       sql.NullString
	RuntimeMinutes sql.NullInt64
	Plot           sql.NullString
	Country        sql.NullString
	Language       sql.NullString
	BoxOffice      sql.NullString
	Production     sql.NullString
	CallFelissa    int64
	Slasher        int64
	Zombies        int64
	Beast          int64
	Godzilla       int64
	WallpaperFu    sql.NullBool
}

func (q *Queries) InsertMovie(ctx context.Context, arg InsertMovieParams) error {
	_, err := q.db.ExecContext(ctx, insertMovie,
		arg.Uuid,
		arg.Title,
		arg.ImdbLink,
		arg.ImdbID,
		arg.Year,
		arg.Rated,
		arg.Released,
		arg.RuntimeMinutes,
		arg.Plot,
		arg.Country,
		arg.Language,
		arg.BoxOffice,
		arg.Production,
		arg.CallFelissa,
		arg.Slasher,
		arg.Zombies,
		arg.Beast,
		arg.Godzilla,
		arg.WallpaperFu,
	)
	return err
}

const insertMovieActor = `-- name: InsertMovieActor :exec
INSERT INTO movie_actor (uuid, movie_uuid, name)
VALUES (?, ?, ?)
`

type InsertMovieActorParams struct {
	Uuid      string
	MovieUuid sql.NullString
	Name      string
}

func (q *Queries) InsertMovieActor(ctx context.Context, arg InsertMovieActorParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieActor, arg.Uuid, arg.MovieUuid, arg.Name)
	return err
}

const insertMovieDirector = `-- name: InsertMovieDirector :exec
INSERT INTO movie_director (uuid, movie_uuid, name)
VALUES (?, ?, ?)
`

type InsertMovieDirectorParams struct {
	Uuid      string
	MovieUuid sql.NullString
	Name      string
}

func (q *Queries) InsertMovieDirector(ctx context.Context, arg InsertMovieDirectorParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieDirector, arg.Uuid, arg.MovieUuid, arg.Name)
	return err
}

const insertMovieGenre = `-- name: InsertMovieGenre :exec
INSERT INTO movie_genre (uuid, movie_uuid, name)
VALUES (?, ?, ?)
`

type InsertMovieGenreParams struct {
	Uuid      string
	MovieUuid sql.NullString
	Name      string
}

func (q *Queries) InsertMovieGenre(ctx context.Context, arg InsertMovieGenreParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieGenre, arg.Uuid, arg.MovieUuid, arg.Name)
	return err
}

const insertMovieRating = `-- name: InsertMovieRating :exec
INSERT INTO movie_rating (
        uuid,
        movie_uuid,
        source,
        value
    )
VALUES (?, ?, ?, ?)
`

type InsertMovieRatingParams struct {
	Uuid      string
	MovieUuid sql.NullString
	Source    string
	Value     string
}

func (q *Queries) InsertMovieRating(ctx context.Context, arg InsertMovieRatingParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieRating,
		arg.Uuid,
		arg.MovieUuid,
		arg.Source,
		arg.Value,
	)
	return err
}

const insertMovieWatch = `-- name: InsertMovieWatch :exec
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertMovieWatchParams struct {
	Uuid       string
	MovieUuid  sql.NullString
	MovieTitle sql.NullString
	ImdbID     string
	Watched    sql.NullString
	Service    string
	FirstTime  int64
	JoeBob     int64
	Notes      sql.NullString
}

func (q *Queries) InsertMovieWatch(ctx context.Context, arg InsertMovieWatchParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieWatch,
		arg.Uuid,
		arg.MovieUuid,
		arg.MovieTitle,
		arg.ImdbID,
		arg.Watched,
		arg.Service,
		arg.FirstTime,
		arg.JoeBob,
		arg.Notes,
	)
	return err
}

const insertMovieWriter = `-- name: InsertMovieWriter :exec
INSERT INTO movie_writer (uuid, movie_uuid, name)
VALUES (?, ?, ?)
`

type InsertMovieWriterParams struct {
	Uuid      string
	MovieUuid sql.NullString
	Name      string
}

func (q *Queries) InsertMovieWriter(ctx context.Context, arg InsertMovieWriterParams) error {
	_, err := q.db.ExecContext(ctx, insertMovieWriter, arg.Uuid, arg.MovieUuid, arg.Name)
	return err
}
