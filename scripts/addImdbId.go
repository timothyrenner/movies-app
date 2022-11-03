package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/timothyrenner/movies-app/cmd"
	"github.com/timothyrenner/movies-app/database"
)

var DB = cmd.DB

func main() {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%v?mode=ro", DB))
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Panicf("Error closing DB: %v", err)
		}
	}()

	log.Println("Getting all movies.")
	allMovieRows, err := db.Query(`SELECT uuid, imdb_link FROM movie`)
	if err != nil {
		log.Panicf("Error getting movie rows: %v", err)
	}

	// We have to load into memory here because we can't have an open read
	// connection _and_ write to the database at the same time, it'll lock.
	movieRows := make([]database.Movie, 0)
	for allMovieRows.Next() {
		movieRow := database.Movie{}
		if err = allMovieRows.Scan(
			&movieRow.Uuid,
			&movieRow.ImdbLink,
		); err != nil {
			log.Panicf("Error scanning row: %v", err)
		}

		urlComponents := strings.Split(movieRow.ImdbLink, "/")
		if urlComponents[len(urlComponents)-1] != "" {
			log.Panicf("Encountered error getting ID from %v", movieRow.ImdbLink)
		}
		movieRow.ImdbID = urlComponents[len(urlComponents)-2]
		movieRows = append(movieRows, movieRow)
	}
	if err = allMovieRows.Close(); err != nil {
		log.Panicf("Error closing rows: %v", err)
	}

	log.Printf("Got %v movies for update.", len(movieRows))
	log.Println("Creating prepared statements.")

	updateMovieStatement, err := db.Prepare(
		`UPDATE movie SET imdb_id = ? WHERE uuid = ?`,
	)
	if err != nil {
		log.Panicf(
			"Error creating movie update prepared statement: %v", err,
		)
	}
	defer updateMovieStatement.Close()

	updateMovieWatchStatement, err := db.Prepare(
		`UPDATE movie_watch SET imdb_id = ? WHERE movie_uuid = ?`,
	)
	if err != nil {
		log.Panicf(
			"Error creating movie watch update prepared statement: %v", err,
		)
	}
	defer updateMovieWatchStatement.Close()

	log.Println("Updating movies and movie watches to add IMDB id.")
	for ii := range movieRows {
		tx, err := db.Begin()
		if err != nil {
			log.Panicf("Error beginning transaction: %v", err)
		}
		defer tx.Rollback()
		// Update the row's imdb id in the database.
		_, err = tx.Stmt(updateMovieStatement).Exec(
			movieRows[ii].ImdbID, movieRows[ii].Uuid,
		)
		if err != nil {
			log.Panicf("Encountered error updating movie row: %v", err)
		}
		// Update the corresponding movie watches too.
		_, err = tx.Stmt(updateMovieWatchStatement).Exec(
			movieRows[ii].ImdbID, movieRows[ii].Uuid,
		)
		if err != nil {
			log.Panicf("Encountered error updating movie watch rows: %v", err)
		}
		if err = tx.Commit(); err != nil {
			log.Panicf("Error committing transaction: %v", err)
		}
	}
	log.Println("Done.")
}
