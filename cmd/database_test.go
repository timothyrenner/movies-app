package cmd

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
)

func setupDatabase() (*DBClient, *migrate.Migrate) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("Encountered error opening in-memory database: %v", err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Panicf(
			"Encountered error creating driver for in-memory database: %v", err,
		)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations", "sqlite3", driver,
	)
	if err != nil {
		log.Panicf("Encountered error creating migration: %v", err)
	}
	if err = m.Up(); err != nil {
		log.Panicf("Encountered error running migration: %v", err)
	}

	// Point the database client var at the new database client.
	c := DBClient{DB: db}
	return &c, m
}

func (c *DBClient) loadMovie() {
	ctx := context.Background()
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Panicf("Encountered error beginning transaction: %v", err)
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
		) VALUES 
		(
			'abc-123',
			'Tenebrae',
			'https://www.imdb.com/title/tt0084777/',
			'tt0084777',
			1982,
			'R',
			'1984-02-17',
			101,
			'An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.',
			'Italy',
			'Italian, Spanish',
			NULL,
			NULL,
			FALSE,
			TRUE,
			FALSE,
			FALSE,
			FALSE,
			TRUE
		), (
			'abc-456',
			'Slaughterhouse',
			'https://www.imdb.com/title/tt0093990/',
			'tt0093990',
			1987,
			'R',
			'1987-08-28',
			85,
			'The owner of a slaughterhouse facing foreclosure instructs his obese and mentally disabled son to go on a killing spree against the people who want to buy his property.',
			'United States',
			'English',
			NULL,
			NULL,
			FALSE,
			TRUE,
			FALSE,
			FALSE,
			FALSE,
			FALSE
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_genre (uuid, movie_uuid, name)
		VALUES
		('genre1', 'abc-123', 'Horror'),
		('genre2', 'abc-123', 'Mystery'),
		('genre3', 'abc-123', 'Thriller'),
		('genre4', 'abc-456', 'Comedy'),
		('genre5', 'abc-456', 'Horror')
		`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie genre table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_actor (uuid, movie_uuid, name)
		VALUES
		('actor1', 'abc-123', 'Anthony Franciosa'),
		('actor2', 'abc-123', 'Giuliano Gemma'),
		('actor3', 'abc-123', 'John Saxon'),
		('actor4', 'abc-456', 'Joe B. Barton'),
		('actor5', 'abc-456', 'Don Barrett'),
		('actor6', 'abc-456', 'Sherry Leigh')
		`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie actor table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_director (uuid, movie_uuid, name)
		VALUES
		('director1', 'abc-123', 'Dario Argento'),
		('director2', 'abc-456', 'Rick Roesller')
		`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie director table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_writer (uuid, movie_uuid, name)
		VALUES
		('writer1', 'abc-123', 'Dario Argento'),
		('writer2', 'abc-456', 'Rick Roessler')
		`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie writer table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_rating (uuid, movie_uuid, source, value)
		VALUES
		('rating1', 'abc-123', 'Internet Movie Database', '7.0/10'),
		('rating2', 'abc-123', 'Rotten Tomatoes', '77%'),
		('rating3', 'abc-123', 'Metacritic', '83/100'),
		('rating4', 'abc-456', 'Internet Movie Database', '5.3/10')`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie rating table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO review (uuid, movie_uuid, movie_title, review, liked)
		VALUES
			('review1', 'abc-456', 'Slaughterhouse', 'Here piggy piggy', TRUE)
		`,
	)
	if err != nil {
		log.Panicf("Encountered error loading review table: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Panicf("Encountered error committing transaction: %v", err)
	}
}

func (c *DBClient) loadMovieWatch() {
	ctx := context.Background()
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Panicf("Encountered error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(
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
		) VALUES (
			'def-123',
			'abc-123',
			'Tenebrae',
			'tt0084777',
			'2022-05-27',
			'Shudder',
			FALSE,
			TRUE,
			'Some notes'
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie_watch: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Panicf("Encountered error committing transaction: %v", err)
	}
}

func teardownDatabase(c *DBClient, m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Panicf("Encountered error tearing down database: %v", err)
	}
	if err := c.Close(); err != nil {
		log.Panicf("Encountered error closing database: %v", err)
	}
}

func omdbSampleMovie() *OmdbMovieResponse {
	return &OmdbMovieResponse{
		Title:    "Tenebrae",
		Year:     "1982",
		Rated:    "R",
		Released: "17 Feb 1984",
		Runtime:  "101 min",
		Genre:    "Horror, Mystery, Thriller",
		Director: "Dario Argento",
		Writer:   "Dario Argento",
		Actors:   "Anthony Franciosa, Giuliano Gemma, John Saxon",
		Plot:     "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.",
		Language: "Italian, Spanish",
		Country:  "Italy",
		Awards:   "N/A",
		Poster:   "https://m.media-amazon.com/images/M/MV5BOTRmNGQ5NTAtNGEzYS00Mjk5LThiZDQtOTk4YTEzNTE1MGZkXkEyXkFqcGdeQXVyNjc1NTYyMjg@._V1_SX300.jpg",
		Ratings: []Rating{
			{Source: "Internet Movie Database", Value: "7.0/10"},
			{Source: "Rotten Tomatoes", Value: "77%"},
			{Source: "Metacritic", Value: "83/100"},
		},
		Metascore:  "83",
		ImdbRating: "7.0",
		ImdbVotes:  "23,156",
		ImdbID:     "tt0084777",
		Type:       "movie",
		DVD:        "20 Sep 2016",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
}

func sampleMovieWatchRow() *MovieWatchRow {
	return &MovieWatchRow{
		MovieTitle: "Tenebrae",
		ImdbId:     "tt0084777",
		FirstTime:  false,
		Watched:    "2022-05-27",
		JoeBob:     true,
		Service:    "Shudder",
		Notes:      textToNullString("Hi there"),
	}
}

func sampleEnrichedMovieWatchRow() *EnrichedMovieWatchRow {
	return &EnrichedMovieWatchRow{
		MovieWatchRow: *sampleMovieWatchRow(),
		CallFelissa:   false,
		Beast:         false,
		Godzilla:      false,
		Zombies:       false,
		Slasher:       true,
		WallpaperFu:   false,
	}
}

func TestFindMovieWatch(t *testing.T) {
	c, m := setupDatabase()
	c.loadMovie()
	c.loadMovieWatch()
	defer teardownDatabase(c, m)

	truth := "def-123"
	record := sampleMovieWatchRow()

	uuid, err := c.FindMovieWatch(record.ImdbId, record.Watched)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if uuid != truth {
		t.Errorf("Expected %v, got %v", truth, uuid)
	}

	uuid2, err := c.FindMovieWatch("tt0093990", "2022-05-27")
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if uuid2 != "" {
		t.Errorf("Expected empty string, got %v", uuid2)
	}
}

func TestGetAllMovieWatches(t *testing.T) {
	c, m := setupDatabase()
	c.loadMovie()
	c.loadMovieWatch()
	defer teardownDatabase(c, m)

	truth := []EnrichedMovieWatchRow{{
		MovieWatchRow: MovieWatchRow{
			Uuid:       "def-123",
			MovieUuid:  "abc-123",
			MovieTitle: "Tenebrae",
			ImdbId:     "tt0084777",
			Watched:    "2022-05-27",
			Service:    "Shudder",
			FirstTime:  false,
			JoeBob:     true,
			Notes:      textToNullString("Some notes"),
		},
		Slasher:     true,
		CallFelissa: false,
		Beast:       false,
		WallpaperFu: false,
		Zombies:     false,
		Godzilla:    false,
		ImdbLink:    "https://www.imdb.com/title/tt0084777/",
	}}

	answer, err := c.GetAllEnrichedMovieWatches()
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, answer)
	}
}

func TestGetLatestMovieWatchDate(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()
	c.loadMovieWatch()

	truth := "2022-05-27"
	answer, err := c.GetLatestMovieWatchDate()
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetMovieDB(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-123"
	truth := MovieRow{
		Uuid:           "abc-123",
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		ImdbId:         "tt0084777",
		Year:           1982,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes: sql.NullInt32{Int32: 101, Valid: true},
		Plot:           sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:        sql.NullString{String: "Italy", Valid: true},
		Language:       sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    false,
		Slasher:        true,
		Zombies:        false,
		Beast:          false,
		Godzilla:       false,
	}
	answer, err := c.GetMovie(movieUuid)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}
}

func TestFindMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()
	movieWatch := sampleMovieWatchRow()
	truth := "abc-123"
	answer, err := c.FindMovie(movieWatch.ImdbId)
	if err != nil {
		t.Errorf("Error getting movie: %v", err)
	}

	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}

}

func TestCreateMovieRow(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieWatch := EnrichedMovieWatchRow{
		MovieWatchRow: MovieWatchRow{
			MovieTitle: "Tenebrae",
			ImdbId:     "tt0084777",
			Watched:    "2022-05-27",
			Service:    "Shudder",
			FirstTime:  false,
			JoeBob:     true,
		},
		CallFelissa: false,
		Beast:       false,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     true,
		WallpaperFu: false,
	}

	movieRow, err := CreateMovieRow(movieRecord, &movieWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	truth := MovieRow{
		Uuid:           movieRow.Uuid,
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		ImdbId:         "tt0084777",
		Year:           1982,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes: sql.NullInt32{Int32: 101, Valid: true},
		Plot:           sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:        sql.NullString{String: "Italy", Valid: true},
		Language:       sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    false,
		Beast:          false,
		Slasher:        true,
		Godzilla:       false,
		WallpaperFu:    false,
	}

	if !cmp.Equal(truth, *movieRow) {
		t.Errorf("Expected %v \n got %v", truth, *movieRow)
	}

	// Now test when there's a null runtime.
	prey := OmdbMovieResponse{
		Title:      "Prey",
		Year:       "2022",
		Rated:      "R",
		Released:   "05 Aug 2022",
		Runtime:    "N/A",
		Genre:      "Action, Drama, Horror",
		Director:   "Dan Trachtenberg",
		Writer:     "Patrick Aison",
		Actors:     "Amber Midthunder, Dane DiLiegro, Harlan Blayne Kytwayhat",
		Plot:       "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.",
		Language:   "English",
		Country:    "United States",
		Awards:     "N/A",
		Poster:     "https://m.media-amazon.com/images/M/MV5BMWE2YjY4MGQtNjRkYy00ZTQxLTkyNTUtODI1Y2I3M2M3ODE2XkEyXkFqcGdeQXVyMTEyMjM2NDc2._V1_SX300.jpg",
		Ratings:    []Rating{},
		Metascore:  "N/A",
		ImdbRating: "N/A",
		ImdbVotes:  "N/A",
		ImdbID:     "tt11866324",
		Type:       "movie",
		DVD:        "05 Aug 2022",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
	preyWatch := EnrichedMovieWatchRow{
		MovieWatchRow: MovieWatchRow{
			MovieTitle: "Prey",
			ImdbId:     "tt11866324",
			Watched:    "2022-08-06",
			Service:    "Hulu",
			FirstTime:  true,
			JoeBob:     false,
		},
		CallFelissa: false,
		Beast:       true,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     false,
		WallpaperFu: false,
	}

	preyRow, err := CreateMovieRow(&prey, &preyWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	preyTruth := MovieRow{
		Uuid:           preyRow.Uuid,
		Title:          "Prey",
		ImdbLink:       "https://www.imdb.com/title/tt11866324/",
		ImdbId:         "tt11866324",
		Year:           2022,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "2022-08-05", Valid: true},
		RuntimeMinutes: sql.NullInt32{Int32: 0, Valid: false},
		Plot:           sql.NullString{String: "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.", Valid: true},
		Country:        sql.NullString{String: "United States", Valid: true},
		Language:       sql.NullString{String: "English", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    false,
		Beast:          true,
		Slasher:        false,
		Godzilla:       false,
		WallpaperFu:    false,
	}

	if !cmp.Equal(preyTruth, *preyRow) {
		t.Errorf("Expected \n%v, got \n%v", preyTruth, *preyRow)
	}

	// Now test when the row is invalid format.
	barbarian := OmdbMovieResponse{
		Title:      "Barbarian",
		Year:       "2022",
		Rated:      "R",
		Released:   "09 Sep 2022",
		Runtime:    "31S min",
		Genre:      "Horror, Thriller",
		Director:   "Zach Cregger",
		Writer:     "Zach Cregger",
		Actors:     "Georgina Campbell, Bill Skarsgård, Justin Long",
		Plot:       "A woman staying at an Airbnb discovers that the house she has rented is not what it seems.",
		Language:   "English",
		Country:    "United States",
		Awards:     "N/A", // <- travesty cause this movie R U L E S
		Poster:     "https://m.media-amazon.com/images/M/MV5BN2M3Y2NhMGYtYjUxOS00M2UwLTlmMGUtYzY4MzFlNjZkYzY2XkEyXkFqcGdeQXVyODc0OTEyNDU@._V1_SX300.jpg",
		Ratings:    []Rating{},
		Metascore:  "N/A",
		ImdbRating: "N/A",
		ImdbVotes:  "N/A",
		ImdbID:     "tt15791034",
		Type:       "movie",
		DVD:        "N/A",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
	barbarianWatch := EnrichedMovieWatchRow{
		MovieWatchRow: MovieWatchRow{
			MovieTitle: "Barbarian",
			ImdbId:     "tt15791034",
			Watched:    "2022-09-12",
			Service:    "Theater",
			FirstTime:  true,
			JoeBob:     false,
		},
		CallFelissa: false,
		Beast:       false,
		Zombies:     false,
		Godzilla:    false,
		Slasher:     false,
		WallpaperFu: true,
	}
	barbarianRow, err := CreateMovieRow(&barbarian, &barbarianWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	barbarianTruth := MovieRow{
		Uuid:           barbarianRow.Uuid,
		Title:          "Barbarian",
		ImdbLink:       "https://www.imdb.com/title/tt15791034/",
		ImdbId:         "tt15791034",
		Year:           2022,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "2022-09-09", Valid: true},
		RuntimeMinutes: sql.NullInt32{Int32: 0, Valid: false},
		Plot:           sql.NullString{String: "A woman staying at an Airbnb discovers that the house she has rented is not what it seems.", Valid: true},
		Country:        sql.NullString{String: "United States", Valid: true},
		Language:       sql.NullString{String: "English", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    false,
		Beast:          false,
		Slasher:        false,
		Godzilla:       false,
		WallpaperFu:    false,
	}
	if !cmp.Equal(barbarianTruth, *barbarianRow) {
		t.Errorf("Expected \n%v, got \n%v", barbarianTruth, *barbarianRow)
	}
}

func TestCreateMovieGenreRow(t *testing.T) {
	movieRecord := omdbSampleMovie()

	answer := CreateMovieGenreRows(movieRecord, "abc-123")
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []MovieGenreRow{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: "abc-123",
			Name:      "Horror",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: "abc-123",
			Name:      "Mystery",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: "abc-123",
			Name:      "Thriller",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestTextToNullString(t *testing.T) {
	na := "N/A"
	naTruth := sql.NullString{}
	naAnswer := textToNullString(na)
	if !cmp.Equal(naTruth, naAnswer) {
		t.Errorf("Expected %v, got %v", naTruth, naAnswer)
	}

	empty := ""
	emptyTruth := sql.NullString{}
	emptyAnswer := textToNullString(empty)
	if !cmp.Equal(emptyTruth, emptyAnswer) {
		t.Errorf("Expected %v, got %v", emptyTruth, emptyAnswer)
	}

	notEmpty := "R"
	notEmptyTruth := sql.NullString{String: "R", Valid: true}
	notEmptyAnswer := textToNullString(notEmpty)
	if !cmp.Equal(notEmptyTruth, notEmptyAnswer) {
		t.Errorf("Expected %v, got %v", notEmptyTruth, notEmptyAnswer)
	}
}

func TestCreateMovieActorRows(t *testing.T) {
	movieRecord := omdbSampleMovie()

	answer := CreateMovieActorRows(movieRecord, "abc-123")
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []MovieActorRow{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: "abc-123",
			Name:      "Anthony Franciosa",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: "abc-123",
			Name:      "Giuliano Gemma",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: "abc-123",
			Name:      "John Saxon",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateMovieDirectorRows(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := "abc-123"

	answer := CreateMovieDirectorRows(movieRecord, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []MovieDirectorRow{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: movieUuid,
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateMovieWriterRows(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := "abc-123"

	answer := CreateMovieWriterRows(movieRecord, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []MovieWriterRow{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: movieUuid,
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateMovieRatingRows(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := "abc-123"

	answer := CreateMovieRatingRows(movieRecord, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}

	truth := []MovieRatingRow{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: "abc-123",
			Source:    "Internet Movie Database",
			Value:     "7.0/10",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: "abc-123",
			Source:    "Rotten Tomatoes",
			Value:     "77%",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: "abc-123",
			Source:    "Metacritic",
			Value:     "83/100",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestInsertMovieDetails(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)

	movie := omdbSampleMovie()
	movieWatch := sampleEnrichedMovieWatchRow()

	answer, err := c.InsertMovieDetails(movie, movieWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	movieRows, err := c.DB.Query(
		`SELECT
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
		FROM movie
		WHERE uuid=?`,
		answer.Movie,
	)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	defer func() {
		if err = movieRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieRowsTruth := []MovieRow{
		{
			Uuid:           answer.Movie,
			Title:          "Tenebrae",
			ImdbLink:       "https://www.imdb.com/title/tt0084777/",
			ImdbId:         "tt0084777",
			Year:           1982,
			Rated:          sql.NullString{String: "R", Valid: true},
			Released:       sql.NullString{String: "1984-02-17", Valid: true},
			RuntimeMinutes: sql.NullInt32{Int32: 101, Valid: true},
			Plot:           sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
			Country:        sql.NullString{String: "Italy", Valid: true},
			Language:       sql.NullString{String: "Italian, Spanish", Valid: true},
			BoxOffice:      sql.NullString{String: "", Valid: false},
			Production:     sql.NullString{String: "", Valid: false},
			CallFelissa:    false,
			Slasher:        true,
			Zombies:        false,
			Beast:          false,
			Godzilla:       false,
			WallpaperFu:    false,
		},
	}
	movieRowsAnswer := make([]MovieRow, 0)
	for movieRows.Next() {
		movieRowAnswer := MovieRow{}
		if err = movieRows.Scan(
			&movieRowAnswer.Uuid,
			&movieRowAnswer.Title,
			&movieRowAnswer.ImdbLink,
			&movieRowAnswer.ImdbId,
			&movieRowAnswer.Year,
			&movieRowAnswer.Rated,
			&movieRowAnswer.Released,
			&movieRowAnswer.RuntimeMinutes,
			&movieRowAnswer.Plot,
			&movieRowAnswer.Country,
			&movieRowAnswer.Language,
			&movieRowAnswer.BoxOffice,
			&movieRowAnswer.Production,
			&movieRowAnswer.CallFelissa,
			&movieRowAnswer.Slasher,
			&movieRowAnswer.Zombies,
			&movieRowAnswer.Beast,
			&movieRowAnswer.Godzilla,
			&movieRowAnswer.WallpaperFu,
		); err != nil {
			t.Errorf("Encountered error scanning movie row: %v", err)
		}
		movieRowsAnswer = append(movieRowsAnswer, movieRowAnswer)
	}
	if !cmp.Equal(movieRowsTruth, movieRowsAnswer) {
		t.Errorf("Expected %v, got %v", movieRowsTruth, movieRowsAnswer)
	}

	// Movie genre rows.

	if len(answer.Genre) != 3 {
		t.Errorf("Expected 3 genre uuids, got %v", len(answer.Genre))
	}
	movieGenreTruth := []MovieGenreRow{
		{
			Uuid:      answer.Genre[0],
			MovieUuid: answer.Movie,
			Name:      "Horror",
		}, {
			Uuid:      answer.Genre[1],
			MovieUuid: answer.Movie,
			Name:      "Mystery",
		}, {
			Uuid:      answer.Genre[2],
			MovieUuid: answer.Movie,
			Name:      "Thriller",
		},
	}

	genreRows, err := c.DB.Query(
		`SELECT 
			uuid, movie_uuid, name 
		FROM movie_genre 
		WHERE movie_uuid = ? 
		`,
		answer.Movie,
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie genre: %v", err)
	}
	defer func() {
		if err = genreRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieGenreAnswer := make([]MovieGenreRow, 0)
	for genreRows.Next() {
		movieGenreRowsAnswer := MovieGenreRow{}
		if err = genreRows.Scan(
			&movieGenreRowsAnswer.Uuid,
			&movieGenreRowsAnswer.MovieUuid,
			&movieGenreRowsAnswer.Name,
		); err != nil {
			t.Errorf("Encountered error scanning genre row: %v", err)
		}
		movieGenreAnswer = append(movieGenreAnswer, movieGenreRowsAnswer)
	}
	if !cmp.Equal(movieGenreTruth, movieGenreAnswer) {
		t.Errorf("Expected %v, got %v", movieGenreTruth, movieGenreAnswer)
	}

	// Movie actor rows.
	if len(answer.Actor) != 3 {
		t.Errorf("Expected 3 actors, got %v", len(answer.Actor))
	}
	movieActorTruth := []MovieActorRow{
		{
			Uuid:      answer.Actor[0],
			MovieUuid: answer.Movie,
			Name:      "Anthony Franciosa",
		}, {
			Uuid:      answer.Actor[1],
			MovieUuid: answer.Movie,
			Name:      "Giuliano Gemma",
		}, {
			Uuid:      answer.Actor[2],
			MovieUuid: answer.Movie,
			Name:      "John Saxon",
		},
	}
	actorRows, err := c.DB.Query(
		`SELECT uuid, movie_uuid, name
		FROM movie_actor WHERE movie_uuid = ?`,
		answer.Movie,
	)
	if err != nil {
		t.Errorf("Encountered error querying movie_actor: %v", err)
	}
	defer func() {
		if err = actorRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieActorAnswer := make([]MovieActorRow, 0)
	for actorRows.Next() {
		movieActorRowAnswer := MovieActorRow{}
		if err = actorRows.Scan(
			&movieActorRowAnswer.Uuid,
			&movieActorRowAnswer.MovieUuid,
			&movieActorRowAnswer.Name,
		); err != nil {
			t.Errorf("Encountered error scanning actor row: %v", err)
		}
		movieActorAnswer = append(movieActorAnswer, movieActorRowAnswer)
	}
	if !cmp.Equal(movieActorTruth, movieActorAnswer) {
		t.Errorf("Expected %v, got %v", movieActorTruth, movieActorAnswer)
	}

	// Movie director rows.
	if len(answer.Director) != 1 {
		t.Errorf("Expected 1 director uuid, got %v", len(answer.Director))
	}
	movieDirectorTruth := []MovieDirectorRow{
		{
			Uuid:      answer.Director[0],
			MovieUuid: answer.Movie,
			Name:      "Dario Argento",
		},
	}
	directorRows, err := c.DB.Query(
		` SELECT uuid, movie_uuid, name
		FROM movie_director WHERE movie_uuid = ?`,
		answer.Movie,
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie director: %v", err)
	}
	defer func() {
		if err = directorRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieDirectorAnswer := make([]MovieDirectorRow, 0)
	for directorRows.Next() {
		movieDirectorRowAnswer := MovieDirectorRow{}
		if err = directorRows.Scan(
			&movieDirectorRowAnswer.Uuid,
			&movieDirectorRowAnswer.MovieUuid,
			&movieDirectorRowAnswer.Name,
		); err != nil {
			t.Errorf("Encountered error scanning director row: %v", err)
		}
		movieDirectorAnswer = append(movieDirectorAnswer, movieDirectorRowAnswer)
	}
	if !cmp.Equal(movieDirectorTruth, movieDirectorAnswer) {
		t.Errorf("Expected %v, got %v", movieDirectorTruth, movieDirectorAnswer)
	}

	// Movie writer rows.
	if len(answer.Writer) != 1 {
		t.Errorf("Expected 1 writer uuid, got %v", len(answer.Writer))
	}
	movieWriterTruth := []MovieWriterRow{
		{
			Uuid:      answer.Writer[0],
			MovieUuid: answer.Movie,
			Name:      "Dario Argento",
		},
	}
	writerRows, err := c.DB.Query(
		`SELECT uuid, movie_uuid, name
		FROM movie_writer WHERE movie_uuid = ?`,
		answer.Movie,
	)

	if err != nil {
		t.Errorf("Encountered error querying for movie writer: %v", err)
	}
	defer func() {
		if err = writerRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieWriterAnswer := make([]MovieWriterRow, 0)
	for writerRows.Next() {
		movieWriterRowAnswer := MovieWriterRow{}
		if err = writerRows.Scan(
			&movieWriterRowAnswer.Uuid,
			&movieWriterRowAnswer.MovieUuid,
			&movieWriterRowAnswer.Name,
		); err != nil {
			t.Errorf("Encountered error scanning writer row: %v", err)
		}
		movieWriterAnswer = append(movieWriterAnswer, movieWriterRowAnswer)
	}
	if !cmp.Equal(movieWriterTruth, movieWriterAnswer) {
		t.Errorf("Expected %v, got %v", movieWriterTruth, movieWriterAnswer)
	}
	// Movie rating rows.
	if len(answer.Rating) != 3 {
		t.Errorf("Expected 1 rating uuid, got %v", len(answer.Rating))
	}
	movieRatingTruth := []MovieRatingRow{
		{
			Uuid:      answer.Rating[0],
			MovieUuid: answer.Movie,
			Source:    "Internet Movie Database",
			Value:     "7.0/10",
		}, {
			Uuid:      answer.Rating[1],
			MovieUuid: answer.Movie,
			Source:    "Rotten Tomatoes",
			Value:     "77%",
		}, {
			Uuid:      answer.Rating[2],
			MovieUuid: answer.Movie,
			Source:    "Metacritic",
			Value:     "83/100",
		},
	}
	ratingRows, err := c.DB.Query(
		`SELECT uuid, movie_uuid, source, value
		FROM movie_rating WHERE movie_uuid = ?`,
		answer.Movie,
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie ratings: %v", err)
	}
	defer func() {
		if err = ratingRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	movieRatingAnswer := make([]MovieRatingRow, 0)
	for ratingRows.Next() {
		movieRatingRowAnswer := MovieRatingRow{}
		if err = ratingRows.Scan(
			&movieRatingRowAnswer.Uuid,
			&movieRatingRowAnswer.MovieUuid,
			&movieRatingRowAnswer.Source,
			&movieRatingRowAnswer.Value,
		); err != nil {
			t.Errorf("Encountered error scanning movie rating row: %v", err)
		}
		movieRatingAnswer = append(movieRatingAnswer, movieRatingRowAnswer)
	}
	if !cmp.Equal(movieRatingTruth, movieRatingAnswer) {
		t.Errorf("Expected %v, got %v", movieRatingTruth, movieRatingAnswer)
	}

	// Now do one where the runtime minutes is null. We just need to make
	// sure this is going to error.
	prey := OmdbMovieResponse{
		Title:      "Prey",
		Year:       "2022",
		Rated:      "R",
		Released:   "05 Aug 2022",
		Runtime:    "N/A",
		Genre:      "Action, Drama, Horror",
		Director:   "Dan Trachtenberg",
		Writer:     "Patrick Aison",
		Actors:     "Amber Midthunder, Dane DiLiegro, Harlan Blayne Kytwayhat",
		Plot:       "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.",
		Language:   "English",
		Country:    "United States",
		Awards:     "N/A",
		Poster:     "https://m.media-amazon.com/images/M/MV5BMWE2YjY4MGQtNjRkYy00ZTQxLTkyNTUtODI1Y2I3M2M3ODE2XkEyXkFqcGdeQXVyMTEyMjM2NDc2._V1_SX300.jpg",
		Ratings:    []Rating{},
		Metascore:  "N/A",
		ImdbRating: "N/A",
		ImdbVotes:  "N/A",
		ImdbID:     "tt11866324",
		Type:       "movie",
		DVD:        "05 Aug 2022",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
	preyWatch := EnrichedMovieWatchRow{
		MovieWatchRow: MovieWatchRow{
			MovieTitle: "Prey",
			ImdbId:     "tt11866324",
			Watched:    "2022-08-06",
			Service:    "Hulu",
			FirstTime:  true,
			JoeBob:     false,
		},
		CallFelissa: false,
		Beast:       true,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     false,
		WallpaperFu: false,
	}
	_, err = c.InsertMovieDetails(&prey, &preyWatch)
	if err != nil {
		t.Errorf("Error inserting movie with null runtime: %v", err)
	}
}

func TestInsertMovieWatch(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieWatchRecord := sampleMovieWatchRow()
	movieUuid := "abc-123"
	movieWatchRecord.MovieUuid = movieUuid

	uuid, err := c.InsertMovieWatch(movieWatchRecord)
	if err != nil {
		t.Errorf("Encountered error inserting movie watch: %v", err)
	}

	answerRows, err := c.DB.Query(
		`SELECT
			uuid,
			movie_uuid,
			movie_title,
			imdb_id,
			watched,
			service,
			first_time,
			joe_bob,
			notes
		FROM movie_watch
		WHERE movie_uuid = ?`, movieUuid,
	)
	if err != nil {
		t.Errorf("Encountered error querying movie row.")
	}
	defer func() {
		if err = answerRows.Close(); err != nil {
			t.Errorf("Encountered error: %v", err)
		}
	}()
	answer := make([]MovieWatchRow, 0)
	for answerRows.Next() {
		answerMovieWatchRow := MovieWatchRow{}
		if err = answerRows.Scan(
			&answerMovieWatchRow.Uuid,
			&answerMovieWatchRow.MovieUuid,
			&answerMovieWatchRow.MovieTitle,
			&answerMovieWatchRow.ImdbId,
			&answerMovieWatchRow.Watched,
			&answerMovieWatchRow.Service,
			&answerMovieWatchRow.FirstTime,
			&answerMovieWatchRow.JoeBob,
			&answerMovieWatchRow.Notes,
		); err != nil {
			t.Errorf("Encountered error scanning movie watch row: %v", err)
		}
		answer = append(answer, answerMovieWatchRow)
	}
	truth := []MovieWatchRow{
		{
			Uuid:       uuid,
			MovieUuid:  "abc-123",
			MovieTitle: "Tenebrae",
			ImdbId:     "tt0084777",
			Watched:    "2022-05-27",
			Service:    "Shudder",
			FirstTime:  false,
			JoeBob:     true,
			Notes:      textToNullString("Hi there"),
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}

}

func TestGetGenreNamesForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-123"
	truth := []string{"Horror", "Mystery", "Thriller"}
	answer, err := c.GetGenreNamesForMovie(movieUuid)
	if err != nil {
		t.Errorf("Error getting genres for movie: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetActorNamesForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-456"
	truth := []string{"Joe B. Barton", "Don Barrett", "Sherry Leigh"}
	answer, err := c.GetActorNamesForMovie(movieUuid)
	if err != nil {
		t.Errorf("Error getting actors for movie: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetDirectorNamesForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-123"
	truth := []string{"Dario Argento"}
	answer, err := c.GetDirectorNamesForMovie(movieUuid)
	if err != nil {
		t.Errorf("Error getting directors for movie: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetWriterNamesForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-123"
	truth := []string{"Dario Argento"}
	answer, err := c.GetWriterNamesForMovie(movieUuid)
	if err != nil {
		t.Errorf("Error getting writers for movie: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetRatingsForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	movieUuid := "abc-123"
	truth := []MovieRatingRow{
		{
			Uuid:      "rating1",
			MovieUuid: "abc-123",
			Source:    "Internet Movie Database",
			Value:     "7.0/10",
		}, {
			Uuid:      "rating2",
			MovieUuid: "abc-123",
			Source:    "Rotten Tomatoes",
			Value:     "77%",
		}, {
			Uuid:      "rating3",
			MovieUuid: "abc-123",
			Source:    "Metacritic",
			Value:     "83/100",
		},
	}
	answer, err := c.GetRatingsForMovie(movieUuid)
	if err != nil {
		t.Errorf("Error getting ratings for movie: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetReviewForMovie(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	truth := MovieReviewRow{
		Uuid:       "review1",
		MovieUuid:  "abc-456",
		MovieTitle: "Slaughterhouse",
		Liked:      true,
		Review:     "Here piggy piggy",
	}

	answer, err := c.GetReviewForMovie("Slaughterhouse")
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}
}

func TestInsertReview(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)
	c.loadMovie()

	review := MovieReviewRow{
		Uuid:       "review2",
		MovieTitle: "Tenebrae",
		MovieUuid:  "abc-123",
		Review:     "Great twist",
		Liked:      true,
	}

	if err := c.InsertReview(&review); err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	var answer MovieReviewRow
	row := c.DB.QueryRow(
		`SELECT uuid, movie_uuid, movie_title, review, liked
		FROM review
		WHERE movie_title='Tenebrae'
		`,
	)
	if err := row.Scan(
		&answer.Uuid,
		&answer.MovieUuid,
		&answer.MovieTitle,
		&answer.Review,
		&answer.Liked,
	); err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	if !cmp.Equal(review, answer) {
		t.Errorf("Expected \n%v, got \n%v", review, answer)
	}
}
