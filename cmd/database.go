package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DBClient struct {
	DB  *sql.DB
	ctx context.Context
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

func CreateMovieRow(
	movieRecord *OmdbMovieResponse, movieWatch *EnrichedMovieWatchRow,
) (*MovieRow, error) {

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

	var runtime sql.NullInt32
	runtimeInt, err := ParseRuntime(movieRecord.Runtime)
	if err != nil {
		log.Printf(
			"Unable to parse %v (%v). setting to null.",
			movieRecord.Runtime, err,
		)
		runtime = sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	} else {
		runtime = sql.NullInt32{
			Int32: int32(runtimeInt),
			Valid: true,
		}
	}

	return &MovieRow{
		Uuid:           uuid.New().String(),
		Title:          movieWatch.MovieTitle,
		ImdbLink:       fmt.Sprintf("https://www.imdb.com/title/%v/", movieWatch.ImdbId),
		ImdbId:         movieWatch.ImdbId,
		Year:           year,
		Rated:          textToNullString(movieRecord.Rated),
		Released:       textToNullString(releasedDate),
		RuntimeMinutes: runtime,
		Plot:           textToNullString(movieRecord.Plot),
		Country:        textToNullString(movieRecord.Country),
		Language:       textToNullString(movieRecord.Language),
		BoxOffice:      textToNullString(movieRecord.BoxOffice),
		Production:     textToNullString(movieRecord.Production),
		CallFelissa:    movieWatch.CallFelissa,
		Slasher:        movieWatch.Slasher,
		Beast:          movieWatch.Beast,
		Godzilla:       movieWatch.Godzilla,
	}, nil
}

type MovieGenreRow struct {
	Uuid      string
	MovieUuid string
	Name      string
}

func CreateMovieGenreRows(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []MovieGenreRow {
	genres := SplitOnCommaAndTrim(movieRecord.Genre)
	rows := make([]MovieGenreRow, len(genres))
	for ii := range genres {
		rows[ii] = MovieGenreRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
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

func CreateMovieActorRows(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []MovieActorRow {
	actors := SplitOnCommaAndTrim(movieRecord.Actors)
	rows := make([]MovieActorRow, len(actors))
	for ii := range actors {
		rows[ii] = MovieActorRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
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

func CreateMovieDirectorRows(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []MovieDirectorRow {
	directors := SplitOnCommaAndTrim(movieRecord.Director)
	rows := make([]MovieDirectorRow, len(directors))
	for ii := range directors {
		rows[ii] = MovieDirectorRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
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

func CreateMovieWriterRows(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []MovieWriterRow {
	writers := SplitOnCommaAndTrim(movieRecord.Writer)
	rows := make([]MovieWriterRow, len(writers))
	for ii := range writers {
		rows[ii] = MovieWriterRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
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

func CreateMovieRatingRows(
	movieRecord *OmdbMovieResponse,
	movieUuid string,
) []MovieRatingRow {
	rows := make([]MovieRatingRow, len(movieRecord.Ratings))
	for ii := range movieRecord.Ratings {
		rows[ii] = MovieRatingRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
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
	movieWatch, err := models.MovieWatches(
		qm.Select(models.MovieWatchColumns.UUID),
		qm.Where(models.MovieWatchColumns.ImdbID+"=?", imdbId),
		qm.Where(models.MovieWatchColumns.Watched+"=?", watched),
	).One(c.ctx, c.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", fmt.Errorf("encountered error with query: %v", err)
		}
	}
	return movieWatch.UUID.String, nil
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
	latestMovieWatch, err := models.MovieWatches(
		qm.Select(models.MovieWatchColumns.Watched),
		qm.OrderBy(fmt.Sprintf("%v DESC", models.MovieWatchColumns.Watched)),
	).One(c.ctx, c.DB)
	if err != nil {
		return "", fmt.Errorf("error getting latest watched date: %v", err)
	}
	return latestMovieWatch.Watched.String, nil
}

func (c *DBClient) FindMovie(imdbId string) (string, error) {
	movie, err := models.Movies(
		qm.Select(models.MovieColumns.UUID),
		qm.Where(models.MovieColumns.ImdbID+"=?", imdbId),
	).One(c.ctx, c.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", fmt.Errorf("encountered error with query: %v", err)
		}
	}
	return movie.UUID.String, nil
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

func (c *DBClient) InsertMovieDetails(
	movie *OmdbMovieResponse,
	movieWatch *EnrichedMovieWatchRow,
) (*MovieDetailUuids, error) {

	movieUuids := MovieDetailUuids{}
	movieRow, err := CreateMovieRow(movie, movieWatch)
	if err != nil {
		return nil, fmt.Errorf("error creating movie row: %v", err)
	}
	movieUuids.Movie = movieRow.Uuid

	ctx := context.Background()
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(
		`INSERT INTO movie (
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
			godzilla,
			wallpaper_fu
		) VALUES(
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`,
		movieRow.Uuid,
		movieRow.Title,
		movieRow.ImdbLink,
		movieRow.ImdbId,
		movieRow.Year,
		movieRow.Rated,
		movieRow.Released,
		movieRow.RuntimeMinutes,
		movieRow.Plot,
		movieRow.Country,
		movieRow.Language,
		movieRow.BoxOffice,
		movieRow.Production,
		movieRow.CallFelissa,
		movieRow.Slasher,
		movieRow.Zombies,
		movieRow.Beast,
		movieRow.Godzilla,
		movieRow.WallpaperFu,
	)

	if err != nil {
		return nil, fmt.Errorf("encountered error inserting movie: %v", err)
	}

	movieGenreRows := CreateMovieGenreRows(movie, movieRow.Uuid)
	movieUuids.Genre = make([]string, len(movieGenreRows))
	values := make([]string, len(movieGenreRows))
	args := make([]any, len(movieGenreRows)*3)

	for ii := range movieGenreRows {
		values[ii] = "(?, ?, ?)"
		args[3*ii] = movieGenreRows[ii].Uuid
		args[3*ii+1] = movieGenreRows[ii].MovieUuid
		args[3*ii+2] = movieGenreRows[ii].Name
		movieUuids.Genre[ii] = movieGenreRows[ii].Uuid
	}
	_, err = tx.Exec(
		fmt.Sprintf(`INSERT INTO movie_genre (
			uuid,
			movie_uuid,
			name
		) VALUES %v
		`, strings.Join(values, ",")),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting movie genre: %v", err)
	}

	movieActorRows := CreateMovieActorRows(movie, movieRow.Uuid)
	movieUuids.Actor = make([]string, len(movieActorRows))
	values = make([]string, len(movieActorRows))
	args = make([]any, len(movieActorRows)*3)

	for ii := range movieActorRows {
		values[ii] = "(?, ?, ?)"
		args[3*ii] = movieActorRows[ii].Uuid
		args[3*ii+1] = movieActorRows[ii].MovieUuid
		args[3*ii+2] = movieActorRows[ii].Name
		movieUuids.Actor[ii] = movieActorRows[ii].Uuid
	}
	_, err = tx.Exec(fmt.Sprintf(
		`INSERT INTO movie_actor (
			uuid,
			movie_uuid,
			name
		) VALUES %v
		`, strings.Join(values, ",")), args...,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting movie actor: %v", err)
	}

	movieDirectorRows := CreateMovieDirectorRows(movie, movieRow.Uuid)
	movieUuids.Director = make([]string, len(movieDirectorRows))
	values = make([]string, len(movieDirectorRows))
	args = make([]any, len(movieDirectorRows)*3)
	for ii := range movieDirectorRows {
		values[ii] = "(?, ?, ?)"
		args[3*ii] = movieDirectorRows[ii].Uuid
		args[3*ii+1] = movieDirectorRows[ii].MovieUuid
		args[3*ii+2] = movieDirectorRows[ii].Name
		movieUuids.Director[ii] = movieDirectorRows[ii].Uuid
	}
	_, err = tx.Exec(fmt.Sprintf(
		`INSERT INTO movie_director (
			uuid,
			movie_uuid,
			name
		) VALUES %v
		`, strings.Join(values, ",")), args...,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting movie director: %v", err)
	}

	movieWriterRows := CreateMovieWriterRows(movie, movieRow.Uuid)
	movieUuids.Writer = make([]string, len(movieWriterRows))
	values = make([]string, len(movieWriterRows))
	args = make([]any, len(movieWriterRows)*3)
	for ii := range movieWriterRows {
		values[ii] = "(?, ?, ?)"
		args[3*ii] = movieWriterRows[ii].Uuid
		args[3*ii+1] = movieWriterRows[ii].MovieUuid
		args[3*ii+2] = movieWriterRows[ii].Name
		movieUuids.Writer[ii] = movieWriterRows[ii].Uuid
	}
	_, err = tx.Exec(fmt.Sprintf(
		`INSERT INTO movie_writer (
			uuid,
			movie_uuid,
			name
		) VALUES %v
		`, strings.Join(values, ",")), args...,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting movie writer: %v", err)
	}

	movieRatingRows := CreateMovieRatingRows(movie, movieRow.Uuid)
	if len(movieRatingRows) > 0 {
		movieUuids.Rating = make([]string, len(movieRatingRows))
		values = make([]string, len(movieRatingRows))
		args = make([]any, len(movieRatingRows)*4)
		for ii := range movieRatingRows {
			values[ii] = "(?, ?, ?, ?)"
			args[4*ii] = movieRatingRows[ii].Uuid
			args[4*ii+1] = movieRatingRows[ii].MovieUuid
			args[4*ii+2] = movieRatingRows[ii].Source
			args[4*ii+3] = movieRatingRows[ii].Value
			movieUuids.Rating[ii] = movieRatingRows[ii].Uuid
		}
		_, err = tx.Exec(fmt.Sprintf(
			`INSERT INTO movie_rating (
			uuid,
			movie_uuid,
			source,
			value
		) VALUES %v
		`, strings.Join(values, ",")), args...,
		)
		if err != nil {
			return nil, fmt.Errorf("error inserting movie rating: %v", err)
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
