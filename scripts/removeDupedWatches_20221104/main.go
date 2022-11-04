package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/timothyrenner/movies-app/cmd"
	"github.com/timothyrenner/movies-app/database"
)

var DB = cmd.DB

var MOVIE_WATCH_UUIDS_TO_REMOVE []string = []string{
	"df3ef888-d400-4f05-92c6-dd8261ec6939",
	"650e41db-8370-42d3-a520-1dc2411b4b16",
	"9962ea1c-7212-4574-950b-068a99267f83",
}
var MOVIE_UUIDS_TO_REMOVE []string = []string{
	"9b021a09-5818-4894-8404-41fb0e6478d9",
}

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

	for ii := range MOVIE_WATCH_UUIDS_TO_REMOVE {

		if err := qtx.DeleteMovieWatch(ctx, MOVIE_WATCH_UUIDS_TO_REMOVE[ii]); err != nil {
			log.Panicf(
				"Error deleting movie watch %v: %v",
				MOVIE_WATCH_UUIDS_TO_REMOVE[ii],
				err,
			)
		}
	}

	for ii := range MOVIE_UUIDS_TO_REMOVE {
		movieUuid := sql.NullString{String: MOVIE_UUIDS_TO_REMOVE[ii], Valid: true}

		// delete actors.
		if err := qtx.DeleteActorsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting actors for movie %v: %v",
				movieUuid.String, err,
			)
		}

		// delete directors
		if err := qtx.DeleteDirectorsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting directors for movie %v: %v",
				movieUuid.String, err,
			)
		}

		// delete writers
		if err := qtx.DeleteWritersForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting writers for movie %v: %v",
				movieUuid.String, err,
			)
		}

		// delete genres
		if err := qtx.DeleteGenresForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting genres for movie %v: %v",
				movieUuid.String, err,
			)
		}

		// delete ratings
		if err := qtx.DeleteRatingsForMovie(ctx, movieUuid); err != nil {
			log.Panicf(
				"Error deleting ratings for movie %v: %v",
				movieUuid.String, err,
			)
		}

		// delete movie
		if err := qtx.DeleteMovie(ctx, movieUuid.String); err != nil {
			log.Panicf(
				"Error deleting movie %v: %v",
				movieUuid.String, err,
			)
		}
	}
	log.Println("Committing updates to database.")
	if err := tx.Commit(); err != nil {
		log.Panicf("Error committing transaction: %v", err)
	}
	log.Println("All done!")
}
