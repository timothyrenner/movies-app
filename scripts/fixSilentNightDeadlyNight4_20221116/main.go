package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/timothyrenner/movies-app/cmd"
	"github.com/timothyrenner/movies-app/database"
)

var MOVIE_WATCH_UUIDS_TO_UPDATE = []string{
	"c2f89e38-34e5-4c36-91e5-1cec0634b331",
}

var MOVIE_UUIDS_TO_REMOVE []string = []string{
	"95dad98f-2af3-4c00-828e-4f1ab4598f83",
}
var DB = cmd.DB
var SILENT_NIGHT_DEADLY_NIGHT_4_CORRECT_UUID = "7fc13177-0ee8-441a-8ccd-22791f2beb65"

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

	for ii := range MOVIE_WATCH_UUIDS_TO_UPDATE {

		updateMovieUuidParams := database.UpdateMovieUuidForWatchParams{
			MovieUuid: SILENT_NIGHT_DEADLY_NIGHT_4_CORRECT_UUID,
			Uuid: MOVIE_WATCH_UUIDS_TO_UPDATE[ii],
		}
		if err := qtx.UpdateMovieUuidForWatch(ctx, updateMovieUuidParams); err != nil {
			log.Panicf(
				"Error deleting movie watch %v: %v",
				MOVIE_WATCH_UUIDS_TO_UPDATE[ii],
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
