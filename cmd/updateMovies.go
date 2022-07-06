/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
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

	dbc, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	db := DBClient{DB: dbc}
	defer db.Close()

	latestMovieWatch, err := db.GetLatestMovieWatchDate()
	if err != nil {
		log.Panicf("Error getting latest movie watch date: %v", err)
	}

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
				path.Join(vaultDir, "Movies", allWatches[ii].Name()),
			)
		}
	}
	log.Printf("Found %v possibly new watches", len(newMovieWatchFiles))
	for ii := range newMovieWatchFiles {
		watchFileContents, err := os.ReadFile(newMovieWatchFiles[ii])
		if err != nil {
			log.Panicf("Error reading %v: %v", watchFileContents, err)
		}
	}

	// Initialize the parser.
	movieWatchParser, err := CreateMovieWatchParser()
	if err != nil {
		log.Panicf("Error creating movie watch parser: %v", err)
	}

	newMovies := 0
	for ii := range newMovieWatchFiles {
		watchFile := newMovieWatchFiles[ii]
		// Parse the watch file.
		movieWatchPage, err := movieWatchParser.ParsePage(watchFile)
		if err != nil {
			log.Panicf("Error parsing page %v: %v", watchFile, err)
		}
		// Determine if it's already in the database.

		movieWatchUuid, err := db.FindMovieWatch(
			movieWatchPage.ImdbId, movieWatchPage.Watched,
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
		movieUuid, err := db.FindMovie(movieWatchPage.ImdbId)
		if err != nil {
			log.Panicf("Error finding movie: %v", err)
		}
		// If the movie's not in the database, we need to insert it and the
		// details.
		movieWatchRow := movieWatchPage.CreateRow()
		if movieUuid == "" {
			log.Printf("Fetching %v from OMDB.", movieWatchPage.Title)
			omdbResponse, err := omdbClient.GetMovie(movieWatchPage.ImdbId)
			if err != nil {
				log.Panicf("Error fetching movie from OMDB: %v", err)
			}
			movieDetailUuids, err := db.InsertMovieDetails(
				omdbResponse, movieWatchRow,
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
		movieWatchRow.MovieUuid = movieUuid
		_, err = db.InsertMovieWatch(movieWatchRow)
		if err != nil {
			log.Panicf(
				"Error inserting movie watch into database: %v", err,
			)
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
