package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/timothyrenner/movies-app/cmd"
	"github.com/timothyrenner/movies-app/database"
)

var MOVIE_WATCH_UUIDS_TO_REMOVE = []string{
	"c05c55a2-7c93-4eff-a321-ec1f5957065f",
}

var MOVIE_UUIDS_TO_REMOVE []string = []string{
	"a8337a90-4919-4618-a047-fe52b781c6db",
	"53e0a021-b6c7-41b5-b82c-4c4ad0d35d02",
	"c519ba6d-595a-4b23-92c5-a10e13fbc767",
	"ddfec4be-4ebd-4763-9e66-d95f948be30a",
	"9c4652ef-2bcf-4c82-a223-05320a3aacbc",
	"21626395-d426-4139-96d9-696fe9b14de6",
	"8c877f48-7369-4b46-b010-7b66ca920549",
}
var DB = cmd.DB

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
