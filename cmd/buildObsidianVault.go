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
	"strings"
	"sync"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/timothyrenner/movies-app/database"
)

// buildObsidianVaultCmd represents the buildObsidianVault command
var buildObsidianVaultCmd = &cobra.Command{
	Use:   "build-obsidian-vault",
	Short: "Builds an Obsidian vault from the movies database.",
	Run:   buildObsidianVault,
	Args:  cobra.RangeArgs(1, 1),
}

func init() {
	rootCmd.AddCommand(buildObsidianVaultCmd)

	// Here you will define your flags and configuration settings.
	buildObsidianVaultCmd.Flags().BoolP(
		"force", "f", false, "Whether to force rebuild the whole vault or not.",
	)
	buildObsidianVaultCmd.Flags().IntP(
		"limit", "l", 0,
		"The maximum number of records to pull. 0 means pull all of them.",
	)
}

func createOrOpenFile(force bool, path string) (*os.File, bool, error) {
	created := false
	var err error
	var file *os.File
	if force {
		log.Printf("Creating and truncating existing file: %v", path)
		file, err = os.Create(path)
		if err != nil {
			return nil, false, fmt.Errorf("error opening %v: %v", path, err)
		}
	} else {
		log.Printf("Creating file if it doesn't exist: %v", path)
		file, err = os.OpenFile(
			path, os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0666,
		)
		if errors.Is(err, os.ErrExist) {
			created = true
		} else if err != nil {
			return nil, false, fmt.Errorf("error opening %v: %v", path, err)
		}
	}

	return file, created, nil
}

func buildObsidianVault(cmd *cobra.Command, args []string) {

	vaultDir := args[0]

	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		log.Panicf("Error obtaining limit: %v", err)
	}
	if limit < 0 {
		log.Panicf("limit must be >= 0, got %v", limit)
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		log.Panicf("Error obtaining force value: %v", err)
	}
	if force {
		log.Println("Rebuilding entire vault (except notes).")
	}

	ctx := context.Background()
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening database %v: %v", DB, err)
	}
	queries := database.New(db)
	defer db.Close()

	// Set up the directories.
	watchesDir := path.Join(vaultDir, "Watches")
	if err = os.Mkdir(watchesDir, 0755); err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Printf("%v exists", watchesDir)
		} else {
			log.Panicf("Error creating %v", watchesDir)
		}
	}
	moviesDir := path.Join(vaultDir, "Movies")
	if err = os.Mkdir(moviesDir, 0755); err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Printf("%v exists", moviesDir)
		} else {
			log.Panicf("Error creating %v", moviesDir)
		}
	}

	// Step 1: Get all the movie watch records.
	// Note: this is should be like ... paginated or something. Future
	// improvement if for some crazy reason memory becomes an issue.
	log.Println("Getting movie watches.")
	allMovieWatches, err := queries.GetAllMovieWatches(ctx)
	if err != nil {
		log.Panicf("Error getting all movie watches: %v", err)
	}

	var movieWatches []database.GetAllMovieWatchesRow
	if limit > 0 {
		movieWatches = allMovieWatches[:limit]
	} else {
		movieWatches = allMovieWatches
	}
	log.Printf("Building vault info for %v watches.", len(movieWatches))

	movieWatchTemplate, err := template.New("movie_watch").Parse(MOVIE_WATCH_TEMPLATE)
	if err != nil {
		log.Panicf("Unable to parse movie watch template: %v", err)
	}
	movieTemplate, err := template.New("movie").Parse(MOVIE_TEMPLATE)
	if err != nil {
		log.Panicf("Unable to parse movie template: %v", err)
	}
	var wg sync.WaitGroup
	for ii := range movieWatches {
		movieWatchPage := CreateMovieWatchPage(&movieWatches[ii])
		// Step 2: If there's no movie watch page, create one.
		wg.Add(1)
		go func() {
			defer wg.Done()
			filename := fmt.Sprintf(
				"%v %v.md", movieWatchPage.Watched, movieWatchPage.FileTitle,
			)

			filePath := path.Join(watchesDir, filename)
			file, skipWatch, err := createOrOpenFile(force, filePath)
			if err != nil {
				log.Panicf("Error creating or opening %v: %v", filePath, err)
			}
			defer file.Close()
			if !skipWatch {
				if err := movieWatchTemplate.Execute(
					file, movieWatchPage,
				); err != nil {
					log.Panicf("Error writing movie watch page: %v", err)
				}
				file.Close()
			}
		}()
		// Step 3: If there's no movie page, retrieve the movie and create one.
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			moviePageFileName := fmt.Sprintf(
				"%v (%v).md", movieWatchPage.FileTitle, movieWatchPage.ImdbId,
			)
			moviePageFilePath := path.Join(moviesDir, moviePageFileName)
			moviePageFile, skipMovie, err := createOrOpenFile(
				force, moviePageFilePath,
			)
			if err != nil {
				log.Panicf(
					"Error creating or opening %v: %v",
					moviePageFilePath, err,
				)
			}
			defer moviePageFile.Close()
			if !skipMovie {
				// Get directors for movie.
				directors, err := queries.GetDirectorNamesForMovie(
					ctx,
					movieWatches[ii].MovieUuid,
				)
				if err != nil {
					log.Panicf(
						"Error getting directors for %v: %v",
						movieWatches[ii].MovieTitle,
						err,
					)
				}

				// Get writers for movie.
				writers, err := queries.GetWriterNamesForMovie(
					ctx,
					movieWatches[ii].MovieUuid,
				)
				if err != nil {
					log.Panicf(
						"Error getting writers for %v: %v",
						movieWatches[ii].MovieTitle,
						err,
					)
				}

				// Get actors for movie.
				actors, err := queries.GetActorNamesForMovie(
					ctx,
					movieWatches[ii].MovieUuid,
				)
				if err != nil {
					log.Panicf(
						"Error getting actors for %v: %v",
						movieWatches[ii].MovieTitle,
						err,
					)
				}

				// Get genres for movie.
				genres, err := queries.GetGenreNamesForMovie(
					ctx,
					movieWatches[ii].MovieUuid,
				)
				if err != nil {
					log.Panicf(
						"Error getting genres for %v: %v",
						movieWatches[ii].MovieTitle,
						err,
					)
				}

				movieRow, err := queries.GetMovie(
					ctx, movieWatches[ii].MovieUuid,
				)
				if err != nil {
					log.Panicf(
						"Error getting movie %v: %v",
						movieWatches[ii].MovieTitle, err,
					)
				}
				moviePage := CreateMoviePageFromRow(
					&movieRow, genres, directors, writers, actors,
				)
				if err := movieTemplate.Execute(
					moviePageFile, moviePage,
				); err != nil {
					log.Panicf(
						"Error writing movie page %v: %v",
						moviePageFilePath, err,
					)
				}
				moviePageFile.Close()
			}
		}(ii)
	}
	wg.Wait()
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, "#", "")
	title = strings.ReplaceAll(title, "^", "")
	title = strings.ReplaceAll(title, "[", "")
	title = strings.ReplaceAll(title, "]", "")
	title = strings.ReplaceAll(title, "|", "")

	return title
}
