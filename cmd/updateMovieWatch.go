/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/timothyrenner/movies-app/database"
)

// updateMovieWatchCmd represents the updateMovieWatch command
var updateMovieWatchCmd = &cobra.Command{
	Use:   "update-movie-watch",
	Short: "Updates a movie watch from a movie watch page.",
	Long: `Pulls the movie from OMDB and creates a movie page if one does not
	exist.`,
	Run: updateMovieWatch,
}

func init() {
	rootCmd.AddCommand(updateMovieWatchCmd)
}

func updateMovieWatch(cmd *cobra.Command, args []string) {
	movieWatchPageFile := args[0]

	if OMDB_KEY == "" {
		log.Panic("OMDB_KEY must be present for this script to run.")
	}

	ctx := context.Background()
	log.Println("Opening database.")
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	defer db.Close()

	queries := database.New(db)

	log.Println("Parsing the movie watch page.")
	parser, err := CreateMovieWatchParser()
	if err != nil {
		log.Panicf("Error creating parser: %v", err)
	}
	page, err := parser.ParsePage(movieWatchPageFile)
	if err != nil {
		log.Panicf("Error parsing %v: %v", movieWatchPageFile, err)
	}

	movieUuid, err := queries.FindMovie(ctx, page.ImdbId)
	if err == sql.ErrNoRows {
		log.Printf("No movie found for %v. Creating one.", page.ImdbId)
	} else if err != nil {
		log.Panicf("Error getting movie uuid for %v: %v", page.ImdbId, err)
	}
	if movieUuid == "" {
		log.Printf("Fetching %v from OMDB.", page.Title)
		omdbClient := NewOmdbClient(OMDB_KEY)
		omdbResponse, err := omdbClient.GetMovie(page.ImdbId)
		if err != nil {
			log.Panicf("Error fetching movie from OMDB: %v", err)
		}

		moviePage, err := CreateMoviePage(omdbResponse, page)
		if err != nil {
			log.Panicf("Error creating movie page: %v", err)
		}

		movieDetailUuids, err := InsertMovieDetails(
			db, ctx, queries, moviePage, omdbResponse.Ratings,
		)
		if err != nil {
			log.Panicf(
				"Error inserting movie details into database: %v", err,
			)
		}
		movieUuid = movieDetailUuids.Movie

		moviePageFileName := fmt.Sprintf(
			"%v (%v).md", page.FileTitle, page.ImdbId,
		)
		// The vault dir is two levels up from the watch page.
		// First call to dir removes the file, second call moves up into the
		// root of the vault.
		vaultDir := path.Dir(path.Dir(movieWatchPageFile))
		moviePageFilePath := path.Join(vaultDir, "Movies", moviePageFileName)
		moviePageFile, skipMovie, err := createOrOpenFile(
			false, moviePageFilePath,
		)
		if err != nil {
			log.Panicf("Error opening file %v: %v", moviePageFilePath, err)
		}
		defer moviePageFile.Close()
		if !skipMovie {
			log.Printf("Creating page %v", moviePageFileName)
			movieTemplate, err := template.New("movie").Parse(MOVIE_TEMPLATE)
			if err != nil {
				log.Panicf("Unable to parse movie template: %v", err)
			}
			if err := movieTemplate.Execute(
				moviePageFile, moviePage,
			); err != nil {
				log.Panicf(
					"Error writing movie page %v: %v",
					moviePageFilePath, err,
				)
			}
		} else {
			log.Printf("Page %v already exists, skipping.", moviePageFileName)
		}
	}

	insertMovieWatchParams := CreateInsertMovieWatchParams(page, movieUuid)
	// Get the watch uuid from the database.
	log.Println("Getting movie watch uuid out of the database.")
	findMovieWatchParams := database.FindMovieWatchParams{
		ImdbID:  page.ImdbId,
		Watched: page.Watched,
	}
	movieWatchUuid, err := queries.FindMovieWatch(ctx, findMovieWatchParams)
	if err == sql.ErrNoRows {
		log.Printf("No movie watch found for %v", findMovieWatchParams)
		log.Println("Inserting one.")
	} else if err != nil {
		log.Panicf("Error finding movie watch %v: %v", findMovieWatchParams, err)
	}

	if movieWatchUuid != "" {
		log.Println("Found a movie watch.")
		insertMovieWatchParams.Uuid = movieWatchUuid
	}

	log.Println("Upserting movie watch into the database.")
	if err := queries.InsertMovieWatch(ctx, *insertMovieWatchParams); err != nil {
		log.Panicf("Error inserting movie watch into database: %v", err)
	}
	log.Println("All done.")
}
