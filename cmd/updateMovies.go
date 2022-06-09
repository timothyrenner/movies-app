/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"
)

// updateMoviesCmd represents the updateMovies command
var updateMoviesCmd = &cobra.Command{
	Use:   "update-movies",
	Short: "Runs the data pipeline for pulling movies.",
	Long: `Pulls the movie watches from Grist and updates the local database.
	Hydrates the movies with additional info from OMDB if required.
	`,
	Run: updateMovies,
}

func updateMovies(cmd *cobra.Command, args []string) {
	if GRIST_KEY == "" {
		log.Panic("GRIST_KEY must be present for this script to run.")
	}

	if GRIST_DOCUMENT_ID == "" {
		log.Panic("GRIST_DOCUMENT_ID must be present for this script to run.")
	}

	if OMDB_KEY == "" {
		log.Panic("OMDB_KEY must be present for this script to run.")
	}

	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		log.Panicf("Error obtaining limit: %v", err)
	}

	gristClient := NewGristClient(GRIST_KEY)
	if limit > 0 {
		log.Printf("Pulling %v records from Grist.", limit)
	} else {
		log.Println("Pulling all records from Grist.")
	}

	omdbClient := NewOmdbClient(OMDB_KEY)

	dbc, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	db := DBClient{DB: dbc}
	defer db.Close()

	records, err := gristClient.GetMovieWatchRecords(
		GRIST_DOCUMENT_ID,
		"Movie_watches",
		nil,
		"-Watched",
		limit,
	)
	if err != nil {
		log.Panicf("Encountered error pulling records from Grist: %v", err)
	}

	log.Printf(
		"Pulled %v documents from Grist, processing.", len(records.Records),
	)

	newMovies := 0

	for ii := range records.Records {
		record := &records.Records[ii]
		// Determine if it's already in the database.
		movieWatchUuid, err := db.FindMovieWatch(record)
		if err != nil {
			log.Panicf("Encountered error obtaining movie watch: %v", err)
		}

		// If there's a uuid for the movie watch in the database, skip it.
		if movieWatchUuid != "" {
			log.Printf(
				"Already found %v in database, skipping.", record.Fields.Name,
			)
			continue
		}

		// See if the movie and details are already in the database.
		movieUuid, err := db.FindMovie(record)
		if err != nil {
			log.Panicf("Error finding movie: %v", err)
		}
		// If the movie's not in the database, we need to insert it and the
		// details.
		if movieUuid == "" {
			log.Printf("Fetching %v from OMDB.", record.Fields.Name)
			omdbResponse, err := omdbClient.GetMovie(record.ImdbId())
			if err != nil {
				log.Panicf("Error fetching movie from OMDB: %v", err)
			}
			movieDetailUuids, err := db.InsertMovieDetails(
				omdbResponse, record,
			)
			if err != nil {
				log.Panicf(
					"Error inserting movie details into database: %v", err,
				)
			}
			movieUuid = movieDetailUuids.Movie
		}
		// Now that we have a movie uuid for the foreign key we can insert the
		// movie watch itself.
		movieWatchUuid, err = db.InsertMovieWatch(record, movieUuid)
		if err != nil {
			log.Panicf(
				"Error inserting movie watch into database: %v", err,
			)
		}
		// Now add the Grist ID <> movie watch ID mapping.
		err = db.InsertUuidGrist(movieWatchUuid, record.Id)
		if err != nil {
			log.Panicf("Error inserting uuid <> grist ID pair into database: %v", err)
		}
		newMovies += 1
	}
	log.Printf("Completed. Inserted %v new movie watches.", newMovies)
}

func init() {
	rootCmd.AddCommand(updateMoviesCmd)
	updateMoviesCmd.Flags().IntP(
		"limit", "l", 25, "The number of movies to pull from Grist.",
	)
}
