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

func setupDatabase() *migrate.Migrate {
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
	DB = db
	return m
}

func loadMovie() {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Panicf("Encountered error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO movie (
			uuid,
			title,
			imdb_link,
			year,
			rated,
			released,
			runtime_minutes,
			plot,
			country,
			box_office,
			production,
			call_felissa,
			slasher,
			zombies,
			beast,
			godzilla
		) VALUES (
			'abc-123',
			'Tenebrae',
			'https://www.imdb.com/title/tt0084777/',
			1982,
			'R',
			'1984-02-17',
			101,
			'An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.',
			'Italy',
			NULL,
			NULL,
			FALSE,
			TRUE,
			FALSE,
			FALSE,
			FALSE
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie table: %v", err)
	}
	if err = tx.Commit(); err != nil {
		log.Panicf("Encountered error committing transaction: %v", err)
	}
}

func loadMovieWatch() {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Panicf("Encountered error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(
		`INSERT INTO movie_watch (
			uuid,
			movie_uuid,
			movie_title,
			watched,
			service,
			first_time,
			joe_bob
		) VALUES (
			'def-123',
			'abc-123',
			'Tenebrae',
			1653609600,
			'Shudder',
			FALSE,
			TRUE
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie_watch: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Panicf("Encountered error committing transaction: %v", err)
	}
}

func teardownDatabase(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Panicf("Encountered error tearing down database: %v", err)
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

func gristSampleMovieWatch() *GristMovieWatchRecord {
	return &GristMovieWatchRecord{
		GristRecord: GristRecord{Id: 1},
		Fields: GristMovieWatchFields{
			Name:        "Tenebrae",
			ImdbLink:    "https://www.imdb.com/title/tt0084777/",
			FirstTime:   false,
			Watched:     1653609600,
			JoeBob:      true,
			CallFelissa: false,
			Beast:       false,
			Godzilla:    false,
			Zombies:     false,
			Slasher:     true,
			Service:     []string{"L", "Shudder"},
		},
	}
}

func TestFindMovieWatch(t *testing.T) {
	m := setupDatabase()
	loadMovie()
	loadMovieWatch()
	defer teardownDatabase(m)

	truth := "def-123"
	record := gristSampleMovieWatch()

	uuid, err := FindMovieWatch(record)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if uuid != truth {
		t.Errorf("Expected %v, got %v", truth, uuid)
	}

	record2 := GristMovieWatchRecord{
		GristRecord: GristRecord{Id: 2},
		Fields: GristMovieWatchFields{
			Name:        "Slaughterhouse",
			ImdbLink:    "",
			FirstTime:   false,
			Watched:     1653609600,
			JoeBob:      true,
			CallFelissa: false,
			Beast:       false,
			Godzilla:    false,
			Zombies:     false,
			Slasher:     true,
			Service:     []string{"L", "Shudder"},
		},
	}

	uuid2, err := FindMovieWatch(&record2)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if uuid2 != "" {
		t.Errorf("Expected empty string, got %v", uuid2)
	}
}

func TestCreateMovieRow(t *testing.T) {
	movieRecord := omdbSampleMovie()
	movieWatch := &GristMovieWatchRecord{
		GristRecord: GristRecord{Id: 1},
		Fields: GristMovieWatchFields{
			Name:        "Tenebrae",
			ImdbLink:    "https://www.imdb.com/title/tt0084777/",
			FirstTime:   false,
			Watched:     1653609600,
			JoeBob:      true,
			CallFelissa: false,
			Beast:       false,
			Godzilla:    false,
			Zombies:     false,
			Slasher:     true,
			Service:     []string{"L", "Shudder"},
		},
	}

	movieRow, err := CreateMovieRow(movieRecord, movieWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	truth := MovieRow{
		Uuid:           movieRow.Uuid,
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		Year:           1982,
		Rated:          "R",
		Released:       "1984-02-17",
		RuntimeMinutes: 101,
		Plot:           "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.",
		Country:        "Italy",
		Language:       "Italian, Spanish",
		BoxOffice:      "N/A",
		Production:     "N/A",
		CallFelissa:    false,
		Beast:          false,
		Slasher:        true,
		Godzilla:       false,
	}

	if !cmp.Equal(truth, *movieRow) {
		t.Errorf("Expected %v \n got %v", truth, *movieRow)
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
