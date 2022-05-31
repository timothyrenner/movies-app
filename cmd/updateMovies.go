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

	movieWatchUpdateChan := make(chan *GristMovieWatchRecord)
	go func() {
		records := make([]GristMovieWatchRecord, 0)
		for record := range movieWatchUpdateChan {
			records = append(records, *record)
		}
		err = gristClient.UpdateMovieWatchRecords(
			GRIST_DOCUMENT_ID,
			"Movie_watches",
			&GristMovieWatchRecords{
				Records: records,
			},
		)
		if err != nil {
			log.Panicf("Error updating movie watch records: %v", err)
		}
	}()

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
				// TODO: Fetch the move info from OMDB and load that into the
				// TODO: database.
				continue
			}
		}
		// TODO: Add the records to the database, and send the UUIDs over
		// TODO: to be updated.
		// Update the records in the database with the grist IDs.
	}()

	for ii := range records.Records {
		record := &records.Records[ii]
		// Determine if it's already in the database.

		// If there's a uuid in the Grist record, we know it's in the DB
		// already. Skip.
		if record.Fields.Uuid != "" {
			continue
		}

		// There might be a case where it _is_ in the database, but has not
		// been synced to Grist yet.
		movieWatchUuid, err := FindMovieWatch(record)
		if err != nil {
			log.Panicf("Encountered error obtaining movie watch: %v", err)
		}

		if movieWatchUuid != "" {
			// TODO: It's worth considering _not_ updating Grist with our
			// TODO: uuids and going straight for the relations instead, as
			// TODO: part of a later step in the pipeline.
			// TODO: We store the Grist ID in our database as its own table:
			// TODO: grist id <> uuid.
			record.Fields.Uuid = movieWatchUuid
			// Update Grist to hold the new Uuid value.
			movieWatchUpdateChan <- record
		} else {
			// Make new record in database.
			movieWatchAddChan <- record
		}
	}
	close(movieWatchUpdateChan)
	close(movieWatchAddChan)
}

func init() {
	rootCmd.AddCommand(updateMoviesCmd)
	updateMoviesCmd.Flags().IntP(
		"limit", "l", 25, "The number of movies to pull from Grist.",
	)
}
