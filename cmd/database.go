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

type MovieRow struct {
	Uuid           string
	Title          string
	ImdbLink       string
	ImdbId         string
	Year           int
	Rated          sql.NullString
	Released       sql.NullString
	RuntimeMinutes sql.NullInt32
	Plot           sql.NullString
	Country        sql.NullString
	Language       sql.NullString
	BoxOffice      sql.NullString
	Production     sql.NullString
	CallFelissa    bool
	Slasher        bool
	Zombies        bool
	Beast          bool
	Godzilla       bool
	WallpaperFu    bool
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

type MovieGenreRow struct {
	Uuid      string
	MovieUuid string
	Name      string
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

type MovieActorRow struct {
	Uuid      string
	MovieUuid string
	Name      string
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

type MovieDirectorRow struct {
	Uuid      string
	MovieUuid string
	Name      string
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

type MovieWriterRow struct {
	Uuid      string
	MovieUuid string
	Name      string
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

type MovieRatingRow struct {
	Uuid      string
	MovieUuid string
	Source    string
	Value     string
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

type MovieWatchRow struct {
	Uuid       string
	MovieUuid  string
	MovieTitle string
	ImdbId     string
	Watched    string
	Service    string
	FirstTime  bool
	JoeBob     bool
	Notes      sql.NullString
}

type EnrichedMovieWatchRow struct {
	MovieWatchRow
	// The extra info that gets recorded when the movie watch happens.
	// It's not how we'd store it in the DB because the info is redundant,
	// however these are all the fields that get recorded.
	ImdbLink    string
	Slasher     bool
	CallFelissa bool
	WallpaperFu bool
	Beast       bool
	Godzilla    bool
	Zombies     bool
}

type MovieDetailUuids struct {
	Movie    string
	Genre    []string
	Actor    []string
	Director []string
	Writer   []string
	Rating   []string
}

func (c *DBClient) FindMovieWatch(imdbId string, watched string) (string, error) {
	query := `
	SELECT
		uuid
	FROM
		movie_watch
	WHERE
		imdb_id = ? AND
		watched = ?
	`

	dbRow := c.DB.QueryRow(
		query, imdbId, watched,
	)

	var uuid string

	if err := dbRow.Scan(&uuid); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			return "", fmt.Errorf("encountered error with query: %v", err)
		}
	}
	return uuid, nil
}

func (c *DBClient) GetAllEnrichedMovieWatches() (
	[]EnrichedMovieWatchRow, error,
) {
	dbRows, err := c.DB.Query(
		`SELECT
			w.uuid,
			w.movie_uuid,
			w.movie_title,
			w.imdb_id,
			w.watched,
			w.service,
			w.first_time,
			w.joe_bob,
			w.notes,
			m.imdb_link,
			m.slasher,
			m.call_felissa,
			m.beast,
			m.godzilla,
			m.zombies
		FROM movie_watch AS w
		JOIN movie AS m
			ON m.uuid = w.movie_uuid`)
	if err != nil {
		return nil, fmt.Errorf("error retrieving movie watches: %v", err)
	}
	defer dbRows.Close()
	movieWatchRows := make([]EnrichedMovieWatchRow, 0)
	for dbRows.Next() {
		movieWatchRow := EnrichedMovieWatchRow{}
		if err := dbRows.Scan(
			&movieWatchRow.Uuid,
			&movieWatchRow.MovieUuid,
			&movieWatchRow.MovieTitle,
			&movieWatchRow.ImdbId,
			&movieWatchRow.Watched,
			&movieWatchRow.Service,
			&movieWatchRow.FirstTime,
			&movieWatchRow.JoeBob,
			&movieWatchRow.Notes,
			&movieWatchRow.ImdbLink,
			&movieWatchRow.Slasher,
			&movieWatchRow.CallFelissa,
			&movieWatchRow.Beast,
			&movieWatchRow.Godzilla,
			&movieWatchRow.Zombies,
		); err != nil {
			return nil, fmt.Errorf("error scanning movie watch row: %v", err)
		}
		movieWatchRows = append(movieWatchRows, movieWatchRow)
	}
	return movieWatchRows, nil
}

func (c *DBClient) GetLatestMovieWatchDate() (string, error) {
	dbRow := c.DB.QueryRow(`SELECT MAX(watched) FROM movie_watch`)
	var watched string
	if err := dbRow.Scan(&watched); err != nil {
		return "", fmt.Errorf("error getting latest watched date: %v", err)
	}
	return watched, nil
}

func (c *DBClient) FindMovie(imdbId string) (string, error) {
	query := `
	SELECT
		uuid
	FROM
		movie
	WHERE
		imdb_id = ?
	`
	dbRow := c.DB.QueryRow(query, imdbId)
	var uuid string
	if err := dbRow.Scan(&uuid); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			return "", fmt.Errorf("encountered error with query: %v", err)
		}
	}

	return uuid, nil
}

func (c *DBClient) GetMovie(movieUuid string) (*MovieRow, error) {
	query := `
	SELECT
		uuid,
		title,
		imdb_link,
		imdb_id,
		year,
		rated,
		released,
		runtime_minutes,
		plot,
		country,
		language,
		box_office,
		production,
		call_felissa,
		slasher,
		zombies,
		beast,
		godzilla
	FROM
		movie
	WHERE
		uuid = ?
	`
	dbRow := c.DB.QueryRow(query, movieUuid)
	var movieRow MovieRow
	if err := dbRow.Scan(
		&movieRow.Uuid,
		&movieRow.Title,
		&movieRow.ImdbLink,
		&movieRow.ImdbId,
		&movieRow.Year,
		&movieRow.Rated,
		&movieRow.Released,
		&movieRow.RuntimeMinutes,
		&movieRow.Plot,
		&movieRow.Country,
		&movieRow.Language,
		&movieRow.BoxOffice,
		&movieRow.Production,
		&movieRow.CallFelissa,
		&movieRow.Slasher,
		&movieRow.Zombies,
		&movieRow.Beast,
		&movieRow.Godzilla,
	); err != nil {
		return nil, fmt.Errorf("error getting movie: %v", err)
	}

	return &movieRow, nil

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

func (c *DBClient) InsertMovieWatch(movieWatch *MovieWatchRow) (string, error) {
	if movieWatch.Uuid == "" {
		movieWatch.Uuid = uuid.NewString()
	}
	_, err := c.DB.Exec(
		`INSERT INTO movie_watch (
			uuid,
			movie_uuid,
			movie_title,
			imdb_id,
			watched,
			service,
			first_time,
			joe_bob,
			notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		movieWatch.Uuid,
		movieWatch.MovieUuid,
		movieWatch.MovieTitle,
		movieWatch.ImdbId,
		movieWatch.Watched,
		movieWatch.Service,
		movieWatch.FirstTime,
		movieWatch.JoeBob,
		movieWatch.Notes,
	)
	if err != nil {
		return "", fmt.Errorf(
			"encountered error inserting movie watch: %v", err,
		)
	}

	return movieWatch.Uuid, nil
}

func (c *DBClient) GetGenreNamesForMovie(movieUuid string) (
	[]string, error,
) {
	rows, err := c.DB.Query(
		`SELECT name FROM movie_genre WHERE movie_uuid = ?`,
		movieUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error making query: %v", err)
	}
	defer rows.Close()

	movieGenreNames := make([]string, 0)
	for rows.Next() {
		var movieGenreName string
		if err := rows.Scan(&movieGenreName); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}
		movieGenreNames = append(movieGenreNames, movieGenreName)
	}
	return movieGenreNames, nil
}

func (c *DBClient) GetActorNamesForMovie(movieUuid string) (
	[]string, error,
) {
	rows, err := c.DB.Query(
		`SELECT name FROM movie_actor WHERE movie_uuid = ?`,
		movieUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error making query: %v", err)
	}
	defer rows.Close()

	movieActorNames := make([]string, 0)
	for rows.Next() {
		var movieActorName string
		if err := rows.Scan(&movieActorName); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}
		movieActorNames = append(movieActorNames, movieActorName)
	}

	return movieActorNames, nil
}

func (c *DBClient) GetDirectorNamesForMovie(movieUuid string) (
	[]string, error,
) {
	rows, err := c.DB.Query(
		`SELECT name FROM movie_director WHERE movie_uuid = ?`,
		movieUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error making query: %v", err)
	}
	defer rows.Close()

	movieDirectorNames := make([]string, 0)
	for rows.Next() {
		var movieDirectorName string
		if err := rows.Scan(&movieDirectorName); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}
		movieDirectorNames = append(movieDirectorNames, movieDirectorName)
	}

	return movieDirectorNames, nil
}

func (c *DBClient) GetWriterNamesForMovie(movieUuid string) ([]string, error) {
	rows, err := c.DB.Query(
		`SELECT name FROM movie_writer WHERE movie_uuid = ?`,
		movieUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error making query: %v", err)
	}
	defer rows.Close()

	movieWriterNames := make([]string, 0)
	for rows.Next() {
		var movieWriterName string
		if err := rows.Scan(&movieWriterName); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}
		movieWriterNames = append(movieWriterNames, movieWriterName)
	}

	return movieWriterNames, nil
}

func (c *DBClient) GetRatingsForMovie(movieUuid string) (
	[]MovieRatingRow, error,
) {
	rows, err := c.DB.Query(
		`SELECT uuid, movie_uuid, source, value 
		FROM movie_rating WHERE movie_uuid = ?`,
		movieUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error making query: %v", err)
	}
	defer rows.Close()

	movieRatings := make([]MovieRatingRow, 0)
	for rows.Next() {
		movieRating := MovieRatingRow{}
		if err := rows.Scan(
			&movieRating.Uuid,
			&movieRating.MovieUuid,
			&movieRating.Source,
			&movieRating.Value,
		); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}
		movieRatings = append(movieRatings, movieRating)
	}
	return movieRatings, nil
}

type MovieReviewRow struct {
	Uuid       string
	MovieUuid  string
	MovieTitle string
	Review     string
	Liked      bool
}

func (c *DBClient) GetReviewForMovie(movieTitle string) (*MovieReviewRow, error) {
	row := c.DB.QueryRow(`
	SELECT uuid, movie_uuid, movie_title, review, liked
	FROM review
	WHERE movie_title = ?
	`, movieTitle)
	var movieReview MovieReviewRow
	if err := row.Scan(
		&movieReview.Uuid,
		&movieReview.MovieUuid,
		&movieReview.MovieTitle,
		&movieReview.Review,
		&movieReview.Liked,
	); err != nil {
		return nil, fmt.Errorf("error getting review for %v: %v", movieTitle, err)
	}

	return &movieReview, nil
}

func (c *DBClient) InsertReview(review *MovieReviewRow) error {

	if _, err := c.DB.Exec(
		`INSERT INTO review (uuid, movie_uuid, movie_title, review, liked)
		VALUES
		(?, ?, ?, ?, ?)
		ON CONFLICT (movie_uuid) DO
		UPDATE SET review = excluded.review, liked = excluded.liked
		`,
		review.Uuid,
		review.MovieUuid,
		review.MovieTitle,
		review.Review,
		review.Liked,
	); err != nil {
		return fmt.Errorf("error inserting review %v: %v", review, err)
	}
	return nil
}
