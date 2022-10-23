package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/database"
)

type DBClient struct {
	DB *sql.DB
}

func (c *DBClient) Close() error {
	if err := c.DB.Close(); err != nil {
		return err
	}
	return nil
}

func CreateInsertMovieParams(
	movieRecord *OmdbMovieResponse, movieWatch *MovieWatchPage,
) (*database.InsertMovieParams, error) {

	year, err := strconv.Atoi(movieRecord.Year)
	if err != nil {
		return nil, fmt.Errorf(
			"error converting year %v to string: %v", movieRecord.Year, err,
		)
	}

	releasedDate, err := ParseReleased(movieRecord.Released)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing date %v: %v", movieRecord.Released, err,
		)
	}

	var runtime sql.NullInt64
	runtimeInt, err := ParseRuntime(movieRecord.Runtime)
	if err != nil {
		log.Printf(
			"Unable to parse %v (%v). setting to null.",
			movieRecord.Runtime, err,
		)
		runtime = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	} else {
		runtime = sql.NullInt64{
			Int64: int64(runtimeInt),
			Valid: true,
		}
	}

	var callFelissa int64
	if movieWatch.CallFelissa {
		callFelissa = 1
	}

	var slasher int64
	if movieWatch.Slasher {
		slasher = 1
	}

	var beast int64
	if movieWatch.Beast {
		beast = 1
	}

	var godzilla int64
	if movieWatch.Godzilla {
		godzilla = 1
	}

	return &database.InsertMovieParams{
		Uuid:           uuid.New().String(),
		Title:          movieWatch.Title,
		ImdbLink:       fmt.Sprintf("https://www.imdb.com/title/%v/", movieWatch.ImdbId),
		ImdbID:         movieWatch.ImdbId,
		Year:           int64(year),
		Rated:          textToNullString(movieRecord.Rated),
		Released:       textToNullString(releasedDate),
		RuntimeMinutes: runtime,
		Plot:           textToNullString(movieRecord.Plot),
		Country:        textToNullString(movieRecord.Country),
		Language:       textToNullString(movieRecord.Language),
		BoxOffice:      textToNullString(movieRecord.BoxOffice),
		Production:     textToNullString(movieRecord.Production),
		CallFelissa:    callFelissa,
		Slasher:        slasher,
		Beast:          beast,
		Godzilla:       godzilla,
		WallpaperFu:    sql.NullBool{Bool: movieWatch.WallpaperFu, Valid: true},
	}, nil
}

func CreateInsertMovieWatchParams(movieWatch *MovieWatchPage, movieUuid string) *database.InsertMovieWatchParams {

	var firstTime int64
	if movieWatch.FirstTime {
		firstTime = 1
	}
	var joeBob int64
	if movieWatch.JoeBob {
		joeBob = 1
	}
	var movieNotes sql.NullString
	if movieWatch.Notes != "" {
		movieNotes.String = movieWatch.Notes
		movieNotes.Valid = true
	}
	return &database.InsertMovieWatchParams{
		Uuid:       uuid.New().String(),
		MovieUuid:  sql.NullString{String: movieUuid, Valid: true},
		MovieTitle: sql.NullString{String: movieWatch.Title, Valid: true},
		ImdbID:     movieWatch.ImdbId,
		Watched:    sql.NullString{String: movieWatch.Watched, Valid: true},
		Service:    movieWatch.Service,
		FirstTime:  firstTime,
		JoeBob:     joeBob,
		Notes:      movieNotes,
	}
}

func CreateInsertMovieGenreParams(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []database.InsertMovieGenreParams {
	genres := SplitOnCommaAndTrim(movieRecord.Genre)
	rows := make([]database.InsertMovieGenreParams, len(genres))
	for ii := range genres {
		rows[ii] = database.InsertMovieGenreParams{
			Uuid:      uuid.New().String(),
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      genres[ii],
		}
	}
	return rows
}

func CreateInsertMovieActorParams(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []database.InsertMovieActorParams {
	actors := SplitOnCommaAndTrim(movieRecord.Actors)
	rows := make([]database.InsertMovieActorParams, len(actors))
	for ii := range actors {
		rows[ii] = database.InsertMovieActorParams{
			Uuid:      uuid.New().String(),
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      actors[ii],
		}
	}
	return rows
}

func textToNullString(text string) sql.NullString {
	if text == "N/A" || len(text) == 0 {
		return sql.NullString{}
	} else {
		return sql.NullString{
			String: text,
			Valid:  true,
		}
	}
}

func CreateInsertMovieDirectorParams(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []database.InsertMovieDirectorParams {
	directors := SplitOnCommaAndTrim(movieRecord.Director)
	rows := make([]database.InsertMovieDirectorParams, len(directors))
	for ii := range directors {
		rows[ii] = database.InsertMovieDirectorParams{
			Uuid:      uuid.New().String(),
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      directors[ii],
		}
	}
	return rows
}

func CreateInsertMovieWriterParams(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []database.InsertMovieWriterParams {
	writers := SplitOnCommaAndTrim(movieRecord.Writer)
	rows := make([]database.InsertMovieWriterParams, len(writers))
	for ii := range writers {
		rows[ii] = database.InsertMovieWriterParams{
			Uuid:      uuid.New().String(),
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      writers[ii],
		}
	}
	return rows
}

func CreateInsertMovieRatingParams(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []database.InsertMovieRatingParams {
	rows := make([]database.InsertMovieRatingParams, len(movieRecord.Ratings))
	for ii := range movieRecord.Ratings {
		rows[ii] = database.InsertMovieRatingParams{
			Uuid:      uuid.New().String(),
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    movieRecord.Ratings[ii].Source,
			Value:     movieRecord.Ratings[ii].Value,
		}
	}
	return rows
}

func CreateInsertMovieReviewParams(
	movieReviewPage *MovieReviewPage,
	movieUuid string,
) *database.InsertReviewParams {
	var liked int64
	if movieReviewPage.Liked {
		liked = 1
	}
	return &database.InsertReviewParams{
		Uuid:       uuid.New().String(),
		MovieUuid:  movieUuid,
		MovieTitle: movieReviewPage.MovieTitle,
		Review:     movieReviewPage.Review,
		Liked:      liked,
	}
}

type MovieDetailUuids struct {
	Movie    string
	Genre    []string
	Actor    []string
	Director []string
	Writer   []string
	Rating   []string
}

func InsertMovieDetails(
	db *sql.DB,
	ctx context.Context,
	queries *database.Queries,
	movie *OmdbMovieResponse,
	movieWatch *MovieWatchPage,
) (*MovieDetailUuids, error) {

	movieUuids := MovieDetailUuids{}
	movieParams, err := CreateInsertMovieParams(movie, movieWatch)
	if err != nil {
		return nil, fmt.Errorf("error creating movie row: %v", err)
	}
	movieUuids.Movie = movieParams.Uuid

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	if err := qtx.InsertMovie(ctx, *movieParams); err != nil {
		return nil, fmt.Errorf("error inserting movie: %v", err)
	}

	movieGenreParams := CreateInsertMovieGenreParams(movie, movieParams.Uuid)
	movieUuids.Genre = make([]string, len(movieGenreParams))

	for ii := range movieGenreParams {
		movieUuids.Genre[ii] = movieGenreParams[ii].Uuid
		if err := qtx.InsertMovieGenre(ctx, movieGenreParams[ii]); err != nil {
			return nil, fmt.Errorf("error inserting movie genre: %v", err)
		}
	}

	movieActorParams := CreateInsertMovieActorParams(movie, movieParams.Uuid)
	movieUuids.Actor = make([]string, len(movieActorParams))

	for ii := range movieActorParams {
		movieUuids.Actor[ii] = movieActorParams[ii].Uuid
		if err := qtx.InsertMovieActor(ctx, movieActorParams[ii]); err != nil {
			return nil, fmt.Errorf("error inserting movie actor: %v", err)
		}
	}

	movieDirectorParams := CreateInsertMovieDirectorParams(movie, movieParams.Uuid)
	movieUuids.Director = make([]string, len(movieDirectorParams))
	for ii := range movieDirectorParams {
		movieUuids.Director[ii] = movieDirectorParams[ii].Uuid
		if err := qtx.InsertMovieDirector(ctx, movieDirectorParams[ii]); err != nil {
			return nil, fmt.Errorf("error inserting movie director: %v", err)
		}
	}

	movieWriterParams := CreateInsertMovieWriterParams(movie, movieParams.Uuid)
	movieUuids.Writer = make([]string, len(movieWriterParams))
	for ii := range movieWriterParams {
		movieUuids.Writer[ii] = movieWriterParams[ii].Uuid
		if err := qtx.InsertMovieWriter(ctx, movieWriterParams[ii]); err != nil {
			return nil, fmt.Errorf("error inserting movie writer: %v", err)
		}
	}

	movieRatingParams := CreateInsertMovieRatingParams(movie, movieParams.Uuid)
	if len(movieRatingParams) > 0 {
		movieUuids.Rating = make([]string, len(movieRatingParams))
		for ii := range movieRatingParams {
			movieUuids.Rating[ii] = movieRatingParams[ii].Uuid
			if err := qtx.InsertMovieRating(ctx, movieRatingParams[ii]); err != nil {
				return nil, fmt.Errorf("error inserting movie rating: %v", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &movieUuids, nil
}
