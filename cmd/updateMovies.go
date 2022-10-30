/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/timothyrenner/movies-app/database"
)

// updateMoviesCmd represents the updateMovies command
var updateMoviesCmd = &cobra.Command{
	Use:   "update-movies",
	Short: "Runs the data pipeline for pulling movies.",
	Long: `Pulls new movie watches from the vault and updates the local database.
	Hydrates the movies with additional info from OMDB if required.
	`,
	Run:  updateMovies,
	Args: cobra.RangeArgs(1, 1),
}

var WATCHED_DATE_EXTRACTOR *regexp.Regexp = regexp.MustCompile(
	`watched::\s*\[\[(\d{4}-\d{2}-\d{2})\]\]`,
)
var IMDB_ID_EXTRACTOR *regexp.Regexp = regexp.MustCompile(
	`imdb_id::\s+(tt\d{7})`,
)

func GetMovieTitleFromWatchFile(fileContents []byte) (string, error) {
	matches := WATCHED_DATE_EXTRACTOR.FindSubmatch(fileContents)
	if len(matches) != 2 {
		return "", errors.New("unable to match watched date")
	}
	return string(matches[1]), nil
}

func GetMovieImdbIdFromWatchFile(fileContents []byte) (string, error) {
	matches := IMDB_ID_EXTRACTOR.FindSubmatch(fileContents)
	if len(matches) != 2 {
		return "", errors.New("unable to match watched date")
	}
	return string(matches[1]), nil
}

func updateMovies(cmd *cobra.Command, args []string) {
	vaultDir := args[0]

	checkAll, err := cmd.Flags().GetBool("check-all")
	if err != nil {
		log.Panicf("Error getting value of check-all: %v", err)
	}

	if OMDB_KEY == "" {
		log.Panic("OMDB_KEY must be present for this script to run.")
	}

	omdbClient := NewOmdbClient(OMDB_KEY)

	ctx := context.Background()
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	defer db.Close()
	queries := database.New(db)

	latestMovieWatchRow, err := queries.GetLatestMovieWatchDate(ctx)
	if err != nil {
		log.Panicf("Error getting latest movie watch date: %v", err)
	}
	// Panics if the conversion doesn't work.
	latestMovieWatch := latestMovieWatchRow.(string)

	newMovieWatchFiles := make([]string, 0)
	allWatches, err := os.ReadDir(path.Join(vaultDir, "Watches"))
	if err != nil {
		log.Panicf("Error reading dir %v", path.Join(vaultDir, "Watches"))
	}
	for ii := range allWatches {
		watchDate := strings.Split(allWatches[ii].Name(), " ")[0]
		if !checkAll && (watchDate < latestMovieWatch) {
			continue
		} else {
			newMovieWatchFiles = append(
				newMovieWatchFiles,
				path.Join(vaultDir, "Watches", allWatches[ii].Name()),
			)
		}
	}
	log.Printf("Found %v possibly new watches", len(newMovieWatchFiles))

	// Initialize the parser.
	movieWatchParser, err := CreateMovieWatchParser()
	if err != nil {
		log.Panicf("Error creating movie watch parser: %v", err)
	}

	// Initialize the template for movie pages.
	movieTemplate, err := template.New("movie").Parse(MOVIE_TEMPLATE)
	if err != nil {
		log.Panicf("Unable to parse movie template: %v", err)
	}

	newMovies := 0
	for ii := range newMovieWatchFiles {
		watchFile := newMovieWatchFiles[ii]
		// Parse the watch file.
		log.Printf("Parsing %v", watchFile)
		movieWatchPage, err := movieWatchParser.ParsePage(watchFile)
		if err != nil {
			log.Panicf("Error parsing page %v: %v", watchFile, err)
		}
		// Determine if it's already in the database.

		movieWatchUuid, err := queries.FindMovieWatch(
			ctx, database.FindMovieWatchParams{
				ImdbID: movieWatchPage.ImdbId,
				Watched: sql.NullString{
					String: movieWatchPage.Watched, Valid: true,
				},
			},
		)
		if err != nil {
			log.Panicf("Encountered error obtaining movie watch: %v", err)
		}

		// If there's a uuid for the movie watch in the database, skip it.
		if movieWatchUuid != "" {
			log.Printf(
				"Already found %v - %v in database, skipping.",
				movieWatchPage.Title, movieWatchPage.Watched,
			)
			continue
		}

		// See if the movie and details are already in the database.
		movieUuid, err := queries.FindMovie(ctx, movieWatchPage.ImdbId)
		if err != nil {
			log.Panicf("Error finding movie: %v", err)
		}
		// If the movie's not in the database, we need to insert it and the
		// details.
		if movieUuid == "" {
			log.Printf("Fetching %v from OMDB.", movieWatchPage.Title)
			omdbResponse, err := omdbClient.GetMovie(movieWatchPage.ImdbId)
			if err != nil {
				log.Panicf("Error fetching movie from OMDB: %v", err)
			}

			moviePage, err := CreateMoviePage(omdbResponse, movieWatchPage)
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
				"%v (%v).md", movieWatchPage.FileTitle, movieWatchPage.ImdbId,
			)
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
			moviePageFile.Close()
		}
		// Now that we have a movie uuid for the foreign key we can insert the
		// movie watch itself.
		movieWatchParams := CreateInsertMovieWatchParams(movieWatchPage, movieUuid)
		if err := queries.InsertMovieWatch(ctx, *movieWatchParams); err != nil {
			log.Panicf("Error inserting movie watch into database: %v", err)
		}
		newMovies += 1
	}
	log.Printf("Completed. Inserted %v new movie watches.", newMovies)
}

func init() {
	rootCmd.AddCommand(updateMoviesCmd)
	updateMoviesCmd.Flags().BoolP(
		"check-all", "c", false, "Whether to check all the movie watches or not.",
	)
}
