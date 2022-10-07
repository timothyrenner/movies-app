/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"log"

	"github.com/spf13/cobra"
)

// updateReviewCmd represents the updateReview command
var updateReviewCmd = &cobra.Command{
	Use:   "update-review",
	Short: "Updates a review in the database from an Obsidian page",
	Run:   updateReview,
	Args:  cobra.RangeArgs(1, 1),
}

func init() {
	rootCmd.AddCommand(updateReviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateReviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateReviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func updateReview(cmd *cobra.Command, args []string) {

	reviewFile := args[0]

	dbc, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	db := DBClient{DB: dbc}
	defer db.Close()

	parser, err := CreateMovieReviewParser()
	if err != nil {
		log.Panicf("Error creating review parser: %v", err)
	}

	page, err := parser.ParseMovieReviewPage(reviewFile)
	if err != nil {
		log.Panicf("Error parsing review page: %v", err)
	}

	// Get the movie uuid from the db for that title.
	movieUuid, err := db.FindMovie(page.ImdbId)
	if err != nil {
		log.Panicf(
			"Error obtaining movie %v (%v): %v",
			page.ImdbId, page.MovieTitle, err,
		)
	}

	movieReviewRow := page.CreateRow(movieUuid)

	if err := db.InsertReview(movieReviewRow); err != nil {
		log.Panicf(
			"Error inserting review for %v: %v", movieReviewRow.MovieTitle, err,
		)
	}
}
