package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var runtimeRegex = regexp.MustCompile("([0-9]+) min")

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
	RuntimeMinutes int
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
}

func CreateMovieRow(
	movieRecord *OmdbMovieResponse, movieWatch *GristMovieWatchRecord,
) (*MovieRow, error) {

	year, err := strconv.Atoi(movieRecord.Year)
	if err != nil {
		return nil, fmt.Errorf(
			"error converting year %v to string: %v", movieRecord.Year, err,
		)
	}

	var releasedDate string
	if movieRecord.Released == "N/A" {
		releasedDate = movieRecord.Released
	} else {
		released, err := time.Parse("2 Jan 2006", movieRecord.Released)
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing date %v: %v", movieRecord.Released, err,
			)
		}
		releasedDate = released.Format("2006-01-02")
	}

	runtimeMatch := runtimeRegex.FindStringSubmatch(movieRecord.Runtime)
	if runtimeMatch == nil {
		return nil, fmt.Errorf("couldn't parse runtime %v", movieRecord.Runtime)
	}
	if len(runtimeMatch) <= 1 {
		return nil, fmt.Errorf("error parsing runtime %v", movieRecord.Runtime)
	}
	runtimeStr := runtimeMatch[1]
	runtimeInt, err := strconv.Atoi(runtimeStr)
	if err != nil {
		return nil, fmt.Errorf(
			"error converting runtime %v to string: %v", runtimeStr, err,
		)
	}

	return &MovieRow{
		Uuid:           uuid.New().String(),
		Title:          movieWatch.Fields.Name,
		ImdbLink:       movieWatch.Fields.ImdbLink,
		ImdbId:         movieWatch.Fields.ImdbId,
		Year:           year,
		Rated:          textToNullString(movieRecord.Rated),
		Released:       textToNullString(releasedDate),
		RuntimeMinutes: runtimeInt,
		Plot:           textToNullString(movieRecord.Plot),
		Country:        textToNullString(movieRecord.Country),
		Language:       textToNullString(movieRecord.Language),
		BoxOffice:      textToNullString(movieRecord.BoxOffice),
		Production:     textToNullString(movieRecord.Production),
		CallFelissa:    movieWatch.Fields.CallFelissa,
		Slasher:        movieWatch.Fields.Slasher,
		Beast:          movieWatch.Fields.Beast,
		Godzilla:       movieWatch.Fields.Godzilla,
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
	genres := strings.Split(movieRecord.Genre, ",")
	rows := make([]MovieGenreRow, len(genres))
	for ii := range genres {
		rows[ii] = MovieGenreRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      strings.TrimSpace(genres[ii]),
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
	actors := strings.Split(movieRecord.Actors, ",")
	rows := make([]MovieActorRow, len(actors))
	for ii := range actors {
		rows[ii] = MovieActorRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      strings.TrimSpace(actors[ii]),
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
	directors := strings.Split(movieRecord.Director, ",")
	rows := make([]MovieDirectorRow, len(directors))
	for ii := range directors {
		rows[ii] = MovieDirectorRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      strings.TrimSpace(directors[ii]),
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
	writers := strings.Split(movieRecord.Writer, ",")
	rows := make([]MovieWriterRow, len(writers))
	for ii := range writers {
		rows[ii] = MovieWriterRow{
			Uuid:      uuid.New().String(),
			MovieUuid: movieUuid,
			Name:      strings.TrimSpace(writers[ii]),
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
	Watched    int
	Service    string
	FirstTime  bool
	JoeBob     bool
}

func CreateMovieWatchRow(
	movieWatchRecord *GristMovieWatchRecord,
	movieUuid string,
) (*MovieWatchRow, error) {
	if len(movieWatchRecord.Fields.Service) != 2 {
		return nil, fmt.Errorf(
			"expected Service to have length 2, got %v",
			len(movieWatchRecord.Fields.Service),
		)
	}

	return &MovieWatchRow{
		Uuid:       uuid.New().String(),
		MovieUuid:  movieUuid,
		MovieTitle: movieWatchRecord.Fields.Name,
		ImdbId:     movieWatchRecord.Fields.ImdbId,
		Watched:    movieWatchRecord.Fields.Watched,
		Service:    movieWatchRecord.Fields.Service[1],
		FirstTime:  movieWatchRecord.Fields.FirstTime,
		JoeBob:     movieWatchRecord.Fields.JoeBob,
	}, nil
}

type UuidGristRow struct {
	Uuid    string
	GristId int
}

type MovieDetailUuids struct {
	Movie    string
	Genre    []string
	Actor    []string
	Director []string
	Writer   []string
	Rating   []string
}

func (c *DBClient) FindMovieWatch(imdbId string, watched int) (string, error) {
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

func (c *DBClient) GetAllMovieWatches() ([]MovieWatchRow, error) {
	dbRows, err := c.DB.Query(
		`SELECT
			uuid,
			movie_uuid,
			movie_title,
			imdb_id,
			watched,
			service,
			first_time,
			joe_bob
		FROM movie_watch`)
	if err != nil {
		return nil, fmt.Errorf("error retrieving movie watches: %v", err)
	}
	defer dbRows.Close()
	movieWatchRows := make([]MovieWatchRow, 0)
	for dbRows.Next() {
		movieWatchRow := MovieWatchRow{}
		if err := dbRows.Scan(
			&movieWatchRow.Uuid,
			&movieWatchRow.MovieUuid,
			&movieWatchRow.MovieTitle,
			&movieWatchRow.ImdbId,
			&movieWatchRow.Watched,
			&movieWatchRow.Service,
			&movieWatchRow.FirstTime,
			&movieWatchRow.JoeBob,
		); err != nil {
			return nil, fmt.Errorf("error scanning movie watch row: %v", err)
		}
		movieWatchRows = append(movieWatchRows, movieWatchRow)
	}
	return movieWatchRows, nil
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

type MovieRowWithGristId struct {
	GristId sql.NullInt64
	MovieRow
}

func (c *DBClient) FindMovieWithGristId(imdbId string) (*MovieRowWithGristId, error) {
	query := `
	SELECT
		movie.uuid,
		uuid_grist.grist_id,
		movie.title,
		movie.imdb_link,
		movie.imdb_id,
		movie.year,
		movie.rated,
		movie.released,
		movie.runtime_minutes,
		movie.plot,
		movie.country,
		movie.language,
		movie.box_office,
		movie.production,
		movie.call_felissa,
		movie.slasher,
		movie.zombies,
		movie.beast,
		movie.godzilla
	FROM movie
	LEFT JOIN uuid_grist ON
		movie.uuid = uuid_grist.uuid
	WHERE imdb_id = ?
	`

	movieRow := MovieRowWithGristId{}
	row := c.DB.QueryRow(query, imdbId)
	if err := row.Scan(
		&movieRow.Uuid,
		&movieRow.GristId,
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
		return nil, fmt.Errorf(
			"encountered error scanning row: %v", err,
		)
	}

	return &movieRow, nil
}

func (c *DBClient) InsertMovieDetails(
	movie *OmdbMovieResponse,
	movieWatch *GristMovieWatchRecord,
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
			godzilla
		) VALUES(
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
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

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &movieUuids, nil
}

func (c *DBClient) InsertMovieWatch(
	movieWatch *GristMovieWatchRecord, movieUuid string,
) (string, error) {
	movieWatchRow, err := CreateMovieWatchRow(movieWatch, movieUuid)
	if err != nil {
		return "", fmt.Errorf(
			"encountered error creating movie watch row: %v", err,
		)
	}
	_, err = c.DB.Exec(
		`INSERT INTO movie_watch (
			uuid,
			movie_uuid,
			movie_title,
			watched,
			service,
			first_time,
			joe_bob
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		movieWatchRow.Uuid,
		movieWatchRow.MovieUuid,
		movieWatchRow.MovieTitle,
		movieWatchRow.Watched,
		movieWatchRow.Service,
		movieWatchRow.FirstTime,
		movieWatchRow.JoeBob,
	)
	if err != nil {
		return "", fmt.Errorf(
			"encountered error inserting movie watch: %v", err,
		)
	}

	return movieWatchRow.Uuid, nil
}

func (c *DBClient) InsertUuidGrist(movieWatchUuid string, gristId int) error {
	_, err := c.DB.Exec(
		`INSERT INTO uuid_grist (uuid, grist_id) VALUES (?, ?)`,
		movieWatchUuid, gristId,
	)
	if err != nil {
		return fmt.Errorf(
			"encountered error inserting uuid <> grist id row: %v", err,
		)
	}
	return nil
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

func (c *DBClient) InsertUuidGristIds(ids []UuidGristRow) error {
	paramStrings := make([]string, len(ids))
	paramValues := make([]any, len(ids)*2)
	for ii := range ids {
		paramStrings[ii] = "(?,?)"
		paramValues[2*ii] = ids[ii].Uuid
		paramValues[2*ii+1] = ids[ii].GristId
	}
	_, err := c.DB.Exec(
		fmt.Sprintf(
			`INSERT INTO uuid_grist (uuid, grist_id) VALUES %v`,
			strings.Join(paramStrings, ","),
		),
		paramValues...,
	)
	if err != nil {
		return fmt.Errorf("error inserting rows: %v", err)
	}
	return nil
}
