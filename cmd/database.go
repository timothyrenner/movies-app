package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/database"
)

func CreateInsertMovieParams(
	moviePage *MoviePage,
) (*database.InsertMovieParams, error) {

	// I can't say I like compiling this regex each time in the function but
	// whatevs.
	imdbIdExtractor := regexp.MustCompile(`\s*https://www\.imdb\.com/title/(tt\d{7,8})/`)
	imdbIdMatch := imdbIdExtractor.FindSubmatch([]byte(moviePage.ImdbLink))
	if len(imdbIdMatch) != 2 {
		return nil, fmt.Errorf("expected 2 matches for IMDB ID, got %v", len(imdbIdMatch))
	}

	var runtime sql.NullInt64
	if moviePage.RuntimeMinutes == 0 {
		runtime = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	} else {
		runtime = sql.NullInt64{
			Int64: int64(moviePage.RuntimeMinutes),
			Valid: true,
		}
	}

	var callFelissa int64
	if moviePage.CallFelissa {
		callFelissa = 1
	}

	var slasher int64
	if moviePage.Slasher {
		slasher = 1
	}

	var beast int64
	if moviePage.Beast {
		beast = 1
	}

	var godzilla int64
	if moviePage.Godzilla {
		godzilla = 1
	}

	var wallpaperFu int64
	if moviePage.WallpaperFu{
		wallpaperFu = 1
	}

	return &database.InsertMovieParams{
		Uuid:           uuid.New().String(),
		Title:          moviePage.Title,
		ImdbLink:       moviePage.ImdbLink,
		ImdbID:         string(imdbIdMatch[1]),
		Year:           int64(moviePage.Year),
		Rated:          textToNullString(moviePage.Rating),
		Released:       textToNullString(moviePage.Released),
		RuntimeMinutes: runtime,
		Plot:           textToNullString(moviePage.Plot),
		Country:        textToNullString(moviePage.Country),
		Language:       textToNullString(moviePage.Language),
		BoxOffice:      textToNullString(moviePage.BoxOffice),
		Production:     textToNullString(moviePage.Production),
		CallFelissa:    callFelissa,
		Slasher:        slasher,
		Beast:          beast,
		Godzilla:       godzilla,
		WallpaperFu:    wallpaperFu,
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
		MovieUuid:  movieUuid,
		MovieTitle: movieWatch.Title,
		ImdbID:     movieWatch.ImdbId,
		Watched:    movieWatch.Watched,
		Service:    movieWatch.Service,
		FirstTime:  firstTime,
		JoeBob:     joeBob,
		Notes:      movieNotes,
	}
}

func CreateInsertMovieGenreParams(
	moviePage *MoviePage,
	movieUuid string,
) []database.InsertMovieGenreParams {
	rows := make([]database.InsertMovieGenreParams, len(moviePage.Genres))
	for ii := range moviePage.Genres {
		rows[ii] = database.InsertMovieGenreParams{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      moviePage.Genres[ii],
		}
	}
	return rows
}

func CreateInsertMovieActorParams(
	moviePage *MoviePage,
	movieUuid string,
) []database.InsertMovieActorParams {
	rows := make([]database.InsertMovieActorParams, len(moviePage.Actors))
	for ii := range moviePage.Actors {
		rows[ii] = database.InsertMovieActorParams{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      moviePage.Actors[ii],
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
	moviePage *MoviePage,
	movieUuid string,
) []database.InsertMovieDirectorParams {
	rows := make([]database.InsertMovieDirectorParams, len(moviePage.Directors))
	for ii := range moviePage.Directors {
		rows[ii] = database.InsertMovieDirectorParams{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      moviePage.Directors[ii],
		}
	}
	return rows
}

func CreateInsertMovieWriterParams(
	moviePage *MoviePage,
	movieUuid string,
) []database.InsertMovieWriterParams {
	rows := make([]database.InsertMovieWriterParams, len(moviePage.Writers))
	for ii := range moviePage.Writers {
		rows[ii] = database.InsertMovieWriterParams{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      moviePage.Writers[ii],
		}
	}
	return rows
}

func CreateInsertMovieRatingParams(
	ratings []Rating,
	movieUuid string,
) []database.InsertMovieRatingParams {
	rows := make([]database.InsertMovieRatingParams, len(ratings))
	for ii := range ratings {
		rows[ii] = database.InsertMovieRatingParams{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Source:    ratings[ii].Source,
			Value:     ratings[ii].Value,
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
	movie *MoviePage,
	ratings []Rating,
) (*MovieDetailUuids, error) {

	movieUuids := MovieDetailUuids{}
	movieParams, err := CreateInsertMovieParams(movie)
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

	movieRatingParams := CreateInsertMovieRatingParams(ratings, movieParams.Uuid)
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
