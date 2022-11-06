// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package database

import (
	"database/sql"
)

type Movie struct {
	Uuid            string
	Title           string
	ImdbLink        string
	Year            int64
	Rated           sql.NullString
	Released        sql.NullString
	Plot            sql.NullString
	Country         sql.NullString
	Language        sql.NullString
	BoxOffice       sql.NullString
	Production      sql.NullString
	CallFelissa     int64
	Slasher         int64
	Zombies         int64
	Beast           int64
	Godzilla        int64
	CreatedDatetime int64
	ImdbID          string
	WallpaperFu     sql.NullBool
	RuntimeMinutes  sql.NullInt64
}

type MovieActor struct {
	Uuid            string
	MovieUuid       string
	Name            string
	CreatedDatetime int64
}

type MovieDirector struct {
	Uuid            string
	MovieUuid       string
	Name            string
	CreatedDatetime int64
}

type MovieGenre struct {
	Uuid            string
	MovieUuid       string
	Name            string
	CreatedDatetime int64
}

type MovieRating struct {
	Uuid            string
	MovieUuid       string
	Source          string
	Value           string
	CreatedDatetime int64
}

type MovieWatch struct {
	Uuid            string
	MovieUuid       string
	MovieTitle      string
	Service         string
	FirstTime       int64
	JoeBob          int64
	CreatedDatetime int64
	ImdbID          string
	Watched         string
	Notes           sql.NullString
}

type MovieWriter struct {
	Uuid            string
	MovieUuid       string
	Name            string
	CreatedDatetime int64
}

type Review struct {
	Uuid            string
	MovieUuid       string
	MovieTitle      string
	Review          string
	Liked           int64
	CreatedDatetime int64
}

type UuidGrist struct {
	Uuid    string
	GristID int64
}