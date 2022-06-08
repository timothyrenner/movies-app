/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// updateMovieDetails.goCmd represents the updateMovieDetails.go command
var updateMovieDetailsCmd = &cobra.Command{
	Use:   "update-movie-details",
	Short: "Pushes movie details to Grist and sets the appropriate relations.",
	Run:   updateMovieDetails,
}

func updateMovieDetails(cmd *cobra.Command, args []string) {
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
	if limit < 0 {
		log.Panicf("Limit must be greater than zero, got %v", limit)
	}

	// Get the movie watches that don't have relations.
	gristClient := NewGristClient(GRIST_KEY)
	movieWatchRecords, err := gristClient.GetMovieWatchRecords(
		GRIST_DOCUMENT_ID,
		"Movie_watches",
		&map[string]any{
			"Movie": []int{0},
		},
		"-Watched",
		limit,
	)
	if err != nil {
		log.Panicf("Encountered error getting movie watches: %v", err)
	}
	log.Printf("Got %v movie watch records", len(movieWatchRecords.Records))
	// TODO: Get movie details with Grist IDs.
	// TODO: Get movie genres.
	// TODO: Get movie actors.
	// TODO: Get movie directors.
	// TODO: Get movie Writers.
	// TODO: Get movie ratings with Grist IDs.
	// TODO: Push movie details to Grist.
	// TODO: Update movie watches with movie IDs.
	// TODO: Save uuid <> grist ID mapping to database.
	// TODO: Push movie ratings to Grist.
	// TODO: Save uuid <> grist ID mapping to database.
	// TODO: Update movie details with rating IDs.

}

func init() {
	rootCmd.AddCommand(updateMovieDetailsCmd)
	updateMovieDetailsCmd.Flags().IntP(
		"limit", "l", 10, "The number of movie details to update.",
	)
}
