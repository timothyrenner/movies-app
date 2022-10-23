/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/spf13/cobra"
	"github.com/timothyrenner/movies-app/database"
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
}

func updateReview(cmd *cobra.Command, args []string) {

	reviewFile := args[0]

	ctx := context.Background()
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	queries := database.New(db)

	parser, err := CreateMovieReviewParser()
	if err != nil {
		log.Panicf("Error creating review parser: %v", err)
	}

	page, err := parser.ParseMovieReviewPage(reviewFile)
	if err != nil {
		log.Panicf("Error parsing review page: %v", err)
	}

	// Get the movie uuid from the db for that title.
	movieUuid, err := queries.FindMovie(ctx, page.ImdbId)
	if err != nil {
		log.Panicf(
			"Error obtaining movie %v (%v): %v",
			page.ImdbId, page.MovieTitle, err,
		)
	}

	movieReviewParams := CreateInsertMovieReviewParams(page, movieUuid)

	if err := queries.InsertReview(ctx, *movieReviewParams); err != nil {
		log.Panicf(
			"Error inserting review for %v: %v",
			movieReviewParams.MovieTitle, err,
		)
	}

	log.Printf("Review successfully updated for %v", reviewFile)
}
