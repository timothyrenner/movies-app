/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/timothyrenner/movies-app/database"
)

var deleteMovieWatchCmd = &cobra.Command{
	Use:   "delete-movie-watch",
	Short: "Deletes a movie watch.",
	Run: deleteMovieWatch,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(deleteMovieWatchCmd)
}

func deleteMovieWatch(cmd *cobra.Command, args []string) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening db: %v", err)
	}
	defer db.Close()

	queries := database.New(db)
	ctx := context.Background()

	movieWatchUuid := args[0]
	validate := func(input string) error {
		decasedInput := strings.ToLower(input)
		if (decasedInput != "y") && (decasedInput != "n") {
			return fmt.Errorf("choose 'y', or 'n', not %v", input)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label: fmt.Sprintf("Really delete movie watch %v ?", movieWatchUuid),
		Validate: validate,
	}
	if _,err := prompt.Run(); err != nil {
		log.Panicf("Aborting delete: %v", err)
	}

	if err := queries.DeleteMovieWatch(ctx, movieWatchUuid); err != nil {
		log.Panicf("Error deleting movie watch: %v", err)
	}
	
}
