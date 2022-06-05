/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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

	movieWatchAddChan := make(chan *GristMovieWatchRecord)
	go func() {
		records := make([]GristMovieWatchRecord, 0)
		for record := range movieWatchAddChan {
			records = append(records, *record)

			movieUuid, err := FindMovie(record)
			if err != nil {
				log.Panicf("Error finding movie: %v", err)
			}
			if movieUuid == "" {
				omdbResponse, err := omdbClient.GetMovie(record.ImdbId())
				if err != nil {
					log.Panicf("Error fetching movie from OMDB: %v", err)
				}
				_, err = InsertMovieDetails(omdbResponse, record)
				if err != nil {
					log.Panicf(
						"Error inserting movie details into database: %v", err,
					)
				}
				continue
			}
		}
		// TODO: Add the records to the database.
		// Update the records in the database with the grist IDs.
	}()

	for ii := range records.Records {
		record := &records.Records[ii]
		// Determine if it's already in the database.
		movieWatchUuid, err := FindMovieWatch(record)
		if err != nil {
			log.Panicf("Encountered error obtaining movie watch: %v", err)
		}

		// If there's a uuid in the Grist record, we know it's in the DB
		// already. Skip.
		if movieWatchUuid != "" {
			log.Printf(
				"Already found %v in database, skipping.", record.Fields.Name,
			)
			continue
		}

		movieWatchAddChan <- record
	}
}

func init() {
	rootCmd.AddCommand(updateMoviesCmd)
	updateMoviesCmd.Flags().IntP(
		"limit", "l", 25, "The number of movies to pull from Grist.",
	)
}
