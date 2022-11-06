package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/timothyrenner/movies-app/cmd"
	"github.com/timothyrenner/movies-app/database"
)

var MOVIE_WATCH_UUIDS_TO_DELETE = []string{
	"d7eabdba-e951-43ea-a7c2-e628655de5d0",
}
var MOVIE_REVIEW_UUIDS_TO_UPDATE = []string{
	"672394ab-918d-46e1-8d41-5a94b661ba23",
}
var MOVIE_UUIDS_TO_REMOVE []string = []string{
	"9ed7592d-ba0d-47cc-a979-a7e9996ccaa7",
}
var DB = cmd.DB
var THINGS_CORRECT_MOVIE_UUID = "18167891-6091-49bc-b41f-c2b6a9b940ff"

func main() {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening db: %v", err)
	}
	defer db.Close()

	queries := database.New(db)
	tx, err := db.Begin()
	if err != nil {
		log.Panicf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	ctx := context.Background()

	for ii := range MOVIE_WATCH_UUIDS_TO_DELETE {
		if err := qtx.DeleteMovieWatch(ctx, MOVIE_WATCH_UUIDS_TO_DELETE[ii]); err != nil {
			log.Panicf(
				"Error deleting movie watch %v: %v",
				MOVIE_WATCH_UUIDS_TO_DELETE[ii],
				err,
			)
		}
	}
	for ii := range MOVIE_REVIEW_UUIDS_TO_UPDATE {
		updateMovieUuidParams := database.UpdateMovieUuidForReviewParams{
			MovieUuid: THINGS_CORRECT_MOVIE_UUID,
			Uuid:      MOVIE_REVIEW_UUIDS_TO_UPDATE[ii],
		}

		if err := qtx.UpdateMovieUuidForReview(ctx, updateMovieUuidParams); err != nil {
			log.Panicf(
				"Error updating review %v: %v",
				MOVIE_REVIEW_UUIDS_TO_UPDATE[ii],
				err,
			)
		}
	}
	for ii := range MOVIE_UUIDS_TO_REMOVE {
		movieUuid := MOVIE_UUIDS_TO_REMOVE[ii]

		// delete actors.
		if err := qtx.DeleteActorsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting actors for movie %v: %v",
				movieUuid, err,
			)
		}

		// delete directors
		if err := qtx.DeleteDirectorsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting directors for movie %v: %v",
				movieUuid, err,
			)
		}

		// delete writers
		if err := qtx.DeleteWritersForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting writers for movie %v: %v",
				movieUuid, err,
			)
		}

		// delete genres
		if err := qtx.DeleteGenresForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting genres for movie %v: %v",
				movieUuid, err,
			)
		}

		// delete ratings
		if err := qtx.DeleteRatingsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting ratings for movie %v: %v",
				movieUuid, err,
			)
		}

		// delete movie
		if err := qtx.DeleteMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting movie %v: %v",
				movieUuid, err,
			)
		}
	}
	log.Println("Committing updates to database.")
	if err := tx.Commit(); err != nil {
		log.Panicf("Error committing transaction: %v", err)
	}
	log.Println("All done!")
}
