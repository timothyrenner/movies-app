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
	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/database"
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

func sampleMovieWatchPage() *MovieWatchPage {
	return &MovieWatchPage{
		Title:       "Tenebrae",
		FileTitle:   "Tenebrae",
		Watched:     "2022-05-27",
		ImdbLink:    "https://www.imdb.com/title/tt0084777/",
		ImdbId:      "tt0084777",
		FirstTime:   false,
		JoeBob:      true,
		CallFelissa: false,
		Beast:       false,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     true,
		WallpaperFu: false,
		Service:     "Shudder",
		Notes:       "",
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

func sampleReviewPage() *MovieReviewPage {
	return &MovieReviewPage{
		MovieTitle: "Things",
		Review:     "YOU HAVE JUST EXPERIENCED ... THINGS",
		Liked:      false,
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

func TestCreateInsertMovieParams(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieWatch := *sampleMovieWatchPage()

	insertMovieParams, err := CreateInsertMovieParams(movieRecord, &movieWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	truth := database.InsertMovieParams{
		Uuid:           insertMovieParams.Uuid,
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		ImdbID:         "tt0084777",
		Year:           1982,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes: sql.NullInt64{Int64: 101, Valid: true},
		Plot:           sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:        sql.NullString{String: "Italy", Valid: true},
		Language:       sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    0,
		Beast:          0,
		Slasher:        1,
		Godzilla:       0,
		WallpaperFu:    sql.NullBool{Bool: false, Valid: true},
	}

	if !cmp.Equal(truth, *insertMovieParams) {
		t.Errorf("Expected %v \n got %v", truth, *insertMovieParams)
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
	preyWatch := MovieWatchPage{
		Title:       "Prey",
		FileTitle:   "Prey",
		Watched:     "2022-08-06",
		ImdbLink:    "https://www.imdb.com/title/tt11866324/",
		ImdbId:      "tt11866324",
		FirstTime:   true,
		JoeBob:      false,
		CallFelissa: false,
		Beast:       true,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     false,
		WallpaperFu: false,
		Service:     "Hulu",
		Notes:       "",
	}

	preyRow, err := CreateInsertMovieParams(&prey, &preyWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	preyTruth := database.InsertMovieParams{
		Uuid:           preyRow.Uuid,
		Title:          "Prey",
		ImdbLink:       "https://www.imdb.com/title/tt11866324/",
		ImdbID:         "tt11866324",
		Year:           2022,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "2022-08-05", Valid: true},
		RuntimeMinutes: sql.NullInt64{Int64: 0, Valid: false},
		Plot:           sql.NullString{String: "The origin story of the Predator in the world of the Comanche Nation 300 years ago. Naru, a skilled female warrior, fights to protect her tribe against one of the first highly-evolved Predators to land on Earth.", Valid: true},
		Country:        sql.NullString{String: "United States", Valid: true},
		Language:       sql.NullString{String: "English", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    0,
		Beast:          1,
		Slasher:        0,
		Godzilla:       0,
		WallpaperFu:    sql.NullBool{Bool: false, Valid: true},
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
	barbarianWatch := MovieWatchPage{
		Title:       "Barbarian",
		FileTitle:   "Barbarian",
		Watched:     "2022-09-12",
		ImdbLink:    "https://www.imdb.com/title/tt15791034/",
		ImdbId:      "tt15791034",
		FirstTime:   true,
		JoeBob:      false,
		CallFelissa: false,
		Beast:       false,
		Zombies:     false,
		Godzilla:    false,
		Slasher:     false,
		WallpaperFu: true,
		Service:     "Theater",
		Notes:       "",
	}
	barbarianRow, err := CreateInsertMovieParams(&barbarian, &barbarianWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	barbarianTruth := database.InsertMovieParams{
		Uuid:           barbarianRow.Uuid,
		Title:          "Barbarian",
		ImdbLink:       "https://www.imdb.com/title/tt15791034/",
		ImdbID:         "tt15791034",
		Year:           2022,
		Rated:          sql.NullString{String: "R", Valid: true},
		Released:       sql.NullString{String: "2022-09-09", Valid: true},
		RuntimeMinutes: sql.NullInt64{Int64: 0, Valid: false},
		Plot:           sql.NullString{String: "A woman staying at an Airbnb discovers that the house she has rented is not what it seems.", Valid: true},
		Country:        sql.NullString{String: "United States", Valid: true},
		Language:       sql.NullString{String: "English", Valid: true},
		BoxOffice:      sql.NullString{String: "", Valid: false},
		Production:     sql.NullString{String: "", Valid: false},
		CallFelissa:    0,
		Beast:          0,
		Slasher:        0,
		Godzilla:       0,
		WallpaperFu:    sql.NullBool{Bool: true, Valid: true},
	}
	if !cmp.Equal(barbarianTruth, *barbarianRow) {
		t.Errorf("Expected \n%v, got \n%v", barbarianTruth, *barbarianRow)
	}
}

func TestCreateInsertMovieWatchParams(t *testing.T) {
	movieWatchPage := sampleMovieWatchPage()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieWatchParams(movieWatchPage, movieUuid)
	// Test with null notes.
	truth := database.InsertMovieWatchParams{
		Uuid:       answer.Uuid,
		MovieUuid:  sql.NullString{String: movieUuid, Valid: true},
		MovieTitle: sql.NullString{String: "Tenebrae", Valid: true},
		ImdbID:     "tt0084777",
		Watched:    sql.NullString{String: "2022-05-27", Valid: true},
		Service:    "Shudder",
		FirstTime:  0,
		JoeBob:     1,
		Notes:      sql.NullString{String: "", Valid: false},
	}

	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}

	movieWatchPageWithNotes := sampleMovieWatchPage()
	movieWatchPageWithNotes.Notes = "Great flick"
	answerNotes := CreateInsertMovieWatchParams(movieWatchPageWithNotes, movieUuid)
	truthNotes := database.InsertMovieWatchParams{
		Uuid:       answerNotes.Uuid,
		MovieUuid:  sql.NullString{String: movieUuid, Valid: true},
		MovieTitle: sql.NullString{String: "Tenebrae", Valid: true},
		ImdbID:     "tt0084777",
		Watched:    sql.NullString{String: "2022-05-27", Valid: true},
		Service:    "Shudder",
		FirstTime:  0,
		JoeBob:     1,
		Notes:      sql.NullString{String: "Great flick", Valid: true},
	}
	if !cmp.Equal(truthNotes, *answerNotes) {
		t.Errorf("Expected \n%v, got \n%v", truthNotes, *answerNotes)
	}
}

func TestCreateInsertMovieGenreParams(t *testing.T) {
	movieRecord := omdbSampleMovie()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieGenreParams(movieRecord, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []database.InsertMovieGenreParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Horror",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Mystery",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
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

func TestCreateInsertMovieActorParams(t *testing.T) {
	movieRecord := omdbSampleMovie()

	movieUuid := uuid.New().String()
	answer := CreateInsertMovieActorParams(movieRecord, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}
	truth := []database.InsertMovieActorParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Anthony Franciosa",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Giuliano Gemma",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "John Saxon",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieDirectorParams(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieDirectorParams(movieRecord, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []database.InsertMovieDirectorParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieWriterParams(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieWriterParams(movieRecord, movieUuid)
	if len(answer) != 1 {
		t.Errorf("Expected 1 row, got %v", len(answer))
	}
	truth := []database.InsertMovieWriterParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Name:      "Dario Argento",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieRatingParams(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieRatingParams(movieRecord, movieUuid)
	if len(answer) != 3 {
		t.Errorf("Expected 3 rows, got %v", len(answer))
	}

	truth := []database.InsertMovieRatingParams{
		{
			Uuid:      answer[0].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Internet Movie Database",
			Value:     "7.0/10",
		}, {
			Uuid:      answer[1].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Rotten Tomatoes",
			Value:     "77%",
		}, {
			Uuid:      answer[2].Uuid,
			MovieUuid: sql.NullString{String: movieUuid, Valid: true},
			Source:    "Metacritic",
			Value:     "83/100",
		},
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestCreateInsertMovieReviewParams(t *testing.T) {
	movieReview := sampleReviewPage()
	movieUuid := uuid.New().String()

	answer := CreateInsertMovieReviewParams(movieReview, movieUuid)
	truth := &database.InsertReviewParams{
		Uuid:       answer.Uuid,
		MovieUuid:  movieUuid,
		MovieTitle: "Things",
		Review:     "YOU HAVE JUST EXPERIENCED ... THINGS",
		Liked:      0,
	}

	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, answer)
	}
}

func TestInsertMovieDetails(t *testing.T) {
	c, m := setupDatabase()
	defer teardownDatabase(c, m)

	movie := omdbSampleMovie()
	movieWatch := sampleMovieWatchPage()

	queries := database.New(c.DB)
	ctx := context.Background()
	answer, err := InsertMovieDetails(
		c.DB,
		ctx,
		queries,
		movie,
		movieWatch,
	)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	movieRowAnswer, err := queries.GetMovie(ctx, answer.Movie)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	movieRowTruth := database.Movie{
		Uuid:            answer.Movie,
		CreatedDatetime: movieRowAnswer.CreatedDatetime,
		Title:           "Tenebrae",
		ImdbLink:        "https://www.imdb.com/title/tt0084777/",
		ImdbID:          "tt0084777",
		Year:            1982,
		Rated:           sql.NullString{String: "R", Valid: true},
		Released:        sql.NullString{String: "1984-02-17", Valid: true},
		RuntimeMinutes:  sql.NullInt64{Int64: 101, Valid: true},
		Plot:            sql.NullString{String: "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.", Valid: true},
		Country:         sql.NullString{String: "Italy", Valid: true},
		Language:        sql.NullString{String: "Italian, Spanish", Valid: true},
		BoxOffice:       sql.NullString{String: "", Valid: false},
		Production:      sql.NullString{String: "", Valid: false},
		CallFelissa:     0,
		Slasher:         1,
		Zombies:         0,
		Beast:           0,
		Godzilla:        0,
		WallpaperFu:     sql.NullBool{Bool: false, Valid: true},
	}

	if !cmp.Equal(movieRowTruth, movieRowAnswer) {
		t.Errorf("Expected \n%v, got \n%v", movieRowTruth, movieRowAnswer)
	}

	// Movie genre rows.
	if len(answer.Genre) != 3 {
		t.Errorf("Expected 3 genres, got %v", len(answer.Genre))
	}

	movieGenreTruth := []string{"Horror", "Mystery", "Thriller"}

	movieGenreAnswer, err := queries.GetGenreNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)

	if err != nil {
		t.Errorf("Encountered error querying for movie genre: %v", err)
	}

	if !cmp.Equal(movieGenreTruth, movieGenreAnswer) {
		t.Errorf("Expected %v, got %v", movieGenreTruth, movieGenreAnswer)
	}

	// Movie actor rows.
	if len(answer.Actor) != 3 {
		t.Errorf("Expected 3 actors, got %v", len(answer.Actor))
	}
	movieActorTruth := []string{"Anthony Franciosa", "Giuliano Gemma", "John Saxon"}
	movieActorAnswer, err := queries.GetActorNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying movie_actor: %v", err)
	}
	if !cmp.Equal(movieActorTruth, movieActorAnswer) {
		t.Errorf("Expected %v, got %v", movieActorTruth, movieActorAnswer)
	}

	// Movie director rows.
	if len(answer.Director) != 1 {
		t.Errorf("Expected 1 director uuid, got %v", len(answer.Director))
	}
	movieDirectorTruth := []string{"Dario Argento"}
	movieDirectorAnswer, err := queries.GetDirectorNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie director: %v", err)
	}
	if !cmp.Equal(movieDirectorTruth, movieDirectorAnswer) {
		t.Errorf("Expected %v, got %v", movieDirectorTruth, movieDirectorAnswer)
	}

	// Movie writer rows.
	if len(answer.Writer) != 1 {
		t.Errorf("Expected 1 writer uuid, got %v", len(answer.Writer))
	}
	movieWriterTruth := []string{"Dario Argento"}
	movieWriterAnswer, err := queries.GetWriterNamesForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie writer: %v", err)
	}
	if !cmp.Equal(movieWriterTruth, movieWriterAnswer) {
		t.Errorf("Expected %v, got %v", movieWriterTruth, movieWriterAnswer)
	}
	// Movie rating rows.
	if len(answer.Rating) != 3 {
		t.Errorf("Expected 3 rating uuids, got %v", len(answer.Rating))
	}
	movieRatingAnswer, err := queries.GetRatingsForMovie(
		ctx, sql.NullString{String: answer.Movie, Valid: true},
	)
	if err != nil {
		t.Errorf("Encountered error querying for movie ratings: %v", err)
	}
	if len(movieRatingAnswer) != 3 {
		t.Errorf("Expected 3 rating answers, got %v", len(movieRatingAnswer))
	}
	movieRatingTruth := []database.MovieRating{
		{
			Uuid:            answer.Rating[0],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Internet Movie Database",
			Value:           "7.0/10",
			CreatedDatetime: movieRatingAnswer[0].CreatedDatetime,
		}, {
			Uuid:            answer.Rating[1],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Rotten Tomatoes",
			Value:           "77%",
			CreatedDatetime: movieRatingAnswer[1].CreatedDatetime,
		}, {
			Uuid:            answer.Rating[2],
			MovieUuid:       sql.NullString{String: answer.Movie, Valid: true},
			Source:          "Metacritic",
			Value:           "83/100",
			CreatedDatetime: movieRatingAnswer[2].CreatedDatetime,
		},
	}

	if !cmp.Equal(movieRatingTruth, movieRatingAnswer) {
		t.Errorf("Expected \n%v, got \n%v", movieRatingTruth, movieRatingAnswer)
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
	preyWatch := MovieWatchPage{
		Title:       "Prey",
		ImdbId:      "tt11866324",
		Watched:     "2022-08-06",
		Service:     "Hulu",
		FirstTime:   true,
		JoeBob:      false,
		CallFelissa: false,
		Beast:       true,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     false,
		WallpaperFu: false,
	}
	_, err = InsertMovieDetails(c.DB, ctx, queries, &prey, &preyWatch)
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

	// Check that upsert works properly.
	reviewUpdate := MovieReviewRow{
		Uuid:       "review3", // application generates uuids, so this won't be the same.
		MovieTitle: "Tenebrae",
		MovieUuid:  "abc-123",
		Review:     "Massive wallpaper fu",
		Liked:      false,
	}

	if err := c.InsertReview(&reviewUpdate); err != nil {
		t.Errorf("Encountered error %v", err)
	}
	var updateAnswer MovieReviewRow
	updateRow := c.DB.QueryRow(
		`SELECT uuid, movie_uuid, movie_title, review, liked
		FROM review
		WHERE movie_title='Tenebrae'
		`,
	)
	if err := updateRow.Scan(
		&updateAnswer.Uuid,
		&updateAnswer.MovieUuid,
		&updateAnswer.MovieTitle,
		&updateAnswer.Review,
		&updateAnswer.Liked,
	); err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	updateTruth := MovieReviewRow{
		Uuid:       "review2", // does not get updated.
		MovieTitle: "Tenebrae",
		MovieUuid:  "abc-123",
		Review:     "Massive wallpaper fu",
		Liked:      false,
	}
	if !cmp.Equal(updateTruth, updateAnswer) {
		t.Errorf("Expected \n%v, got \n%v", updateTruth, updateAnswer)
	}
}
