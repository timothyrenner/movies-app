/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
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

	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Encountered error opening database %v: %v", DB, err)
	}
	dbClient := DBClient{DB: db}

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

	movieWatchUpdates := make(
		[]GristMovieWatchRecord, len(movieWatchRecords.Records),
	)
	for ii := range movieWatchRecords.Records {
		// Get movie details with Grist ID.
		movieWithGristId, err := dbClient.FindMovieWithGristId(
			movieWatchRecords.Records[ii].Fields.ImdbId,
		)
		if err != nil {
			log.Panicf("Encountered error finding movie with grist id: %v", err)
		}

		if movieWithGristId.GristId.Valid {
			movieWatchUpdates[ii] = GristMovieWatchRecord{
				GristRecord: GristRecord{Id: movieWatchUpdates[ii].Id},
				Fields: GristMovieWatchFields{
					Movie: movieWithGristId.GristId.Int64,
				},
			}
		} else {
			// If there's not a Grist ID associated with the movie, we need
			// to construct the record and push to grist to create one.
			// Note we don't need to explicitly handle nulls here because the
			// default value for a null will result in an omitted field when
			// serializing to json anyway.
			gristMovieRecords := &GristMovieRecords{
				Records: []GristMovieRecord{
					{
						Fields: GristMovieFields{
							Title:       movieWithGristId.Title,
							ImdbLink:    movieWithGristId.ImdbLink,
							Year:        movieWithGristId.Year,
							Rated:       movieWithGristId.Rated.String,
							Released:    movieWithGristId.Released.String,
							Runtime:     movieWithGristId.RuntimeMinutes,
							Plot:        movieWithGristId.Plot.String,
							Country:     movieWithGristId.Country.String,
							Language:    movieWithGristId.Language.String,
							BoxOffice:   movieWithGristId.BoxOffice.String,
							Production:  movieWithGristId.Production.String,
							CallFelissa: movieWithGristId.CallFelissa,
							Slasher:     movieWithGristId.Slasher,
							Zombies:     movieWithGristId.Zombies,
							Beast:       movieWithGristId.Beast,
							Godzilla:    movieWithGristId.Godzilla,
						},
					},
				},
			}
			// Get movie genres.
			movieGenres, err := dbClient.GetGenreNamesForMovie(
				movieWithGristId.Uuid,
			)
			if err != nil {
				log.Panicf("Error getting genres for movie: %v", err)
			}
			gristMovieRecords.Records[0].Fields.Genre = movieGenres
			// Get movie actors.
			movieActors, err := dbClient.GetActorNamesForMovie(
				movieWithGristId.Uuid,
			)
			if err != nil {
				log.Panicf("Error getting directors for movie: %v", err)
			}
			gristMovieRecords.Records[0].Fields.Actor = movieActors
			// Get movie directors.
			movieDirectors, err := dbClient.GetDirectorNamesForMovie(
				movieWithGristId.Uuid,
			)
			if err != nil {
				log.Panicf("Error getting directors for movie: %v", err)
			}
			gristMovieRecords.Records[0].Fields.Director = movieDirectors
			// Get movie writers.
			movieWriters, err := dbClient.GetWriterNamesForMovie(
				movieWithGristId.Uuid,
			)
			if err != nil {
				log.Panicf("Error getting writers for movie: %v", err)
			}
			gristMovieRecords.Records[0].Fields.Writer = movieWriters
			// Get movie ratings.
			movieRatingRows, err := dbClient.GetRatingsForMovie(
				movieWithGristId.Uuid,
			)
			if err != nil {
				log.Panicf("Error getting ratings for movie: %v", err)
			}

			// Push movie ratings to Grist.
			movieRatingRecords := GristMovieRatingRecords{}
			movieRatingRecords.Records = make(
				[]GristMovieRatingRecord, len(movieRatingRows),
			)
			for ii := range movieRatingRows {
				movieRatingRecords.Records[ii] = GristMovieRatingRecord{
					Fields: GristMovieRatingFields{
						Source: movieRatingRows[ii].Source,
						Value:  movieRatingRows[ii].Value,
					},
				}
			}

			movieRatingGristIds, err := gristClient.CreateMovieRatingRecords(
				GRIST_DOCUMENT_ID,
				"Movie_Rating",
				&movieRatingRecords,
			)
			if err != nil {
				log.Panicf("Error creating Grist records for ratings: %v", err)
			}
			gristMovieRecords.Records[0].Fields.Rating = make(
				[]any, len(movieRatingGristIds.Records)+1,
			)
			gristMovieRecords.Records[0].Fields.Rating[0] = "r"
			for ii := range movieRatingGristIds.Records {
				gristMovieRecords.Records[0].Fields.Rating[ii+1] =
					movieRatingGristIds.Records[ii].Id
			}

			// Push movie details to Grist.
			movieGristId, err := gristClient.CreateMovieRecords(
				GRIST_DOCUMENT_ID,
				"Movies",
				gristMovieRecords,
			)
			if err != nil {
				log.Panicf("Error creating Grist records for movie: %v", err)
			}
			if len(movieGristId.Records) != 1 {
				log.Panicf(
					"Expected 1 record for creating 1 movie, got %v",
					len(movieGristId.Records),
				)
			}
			movieWatchUpdates[ii].Fields.Movie = int64(movieGristId.Records[0].Id)
			// Save uuid <> grist ID mappings to database.
			uuidGristIds := make([]UuidGristRow, len(movieRatingGristIds.Records)+1)
			for ii := range movieRatingGristIds.Records {
				uuidGristIds[ii] = UuidGristRow{
					Uuid:    movieRatingRows[ii].Uuid,
					GristId: movieRatingGristIds.Records[ii].Id,
				}
			}
			uuidGristIds[len(movieRatingGristIds.Records)] = UuidGristRow{
				Uuid:    movieWithGristId.Uuid,
				GristId: movieGristId.Records[0].Id,
			}

			if err = dbClient.InsertUuidGristIds(uuidGristIds); err != nil {
				log.Panicf("Error inserting uuid <> Grist IDs: %v", err)
			}
		}
	}
	// Send movie watch updates to Grist.
	if err = gristClient.UpdateMovieWatchRecords(
		GRIST_DOCUMENT_ID,
		"Movie_watches",
		&GristMovieWatchRecords{Records: movieWatchUpdates},
	); err != nil {
		log.Panicf("Error updating movie watch records: %v", err)
	}
	log.Println("All done!")
}

func init() {
	rootCmd.AddCommand(updateMovieDetailsCmd)
	updateMovieDetailsCmd.Flags().IntP(
		"limit", "l", 10, "The number of movie details to update.",
	)
}
