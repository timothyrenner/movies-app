package cmd

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setupDatabase() *migrate.Migrate {
	ctx := context.Background()
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

	tx, err := db.BeginTx(ctx, nil)
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
			godzilla,
			grist_id
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
			FALSE,
			1
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie table: %v", err)
	}

	_, err = tx.Exec(
		`INSERT INTO movie_watch (
			uuid,
			movie_uuid,
			movie_title,
			watched,
			service,
			first_time,
			joe_bob,
			grist_id
		) VALUES (
			'def-123',
			'abc-123',
			'Tenebrae',
			1653609600,
			'Shudder',
			FALSE,
			TRUE,
			1
		)`,
	)
	if err != nil {
		log.Panicf("Encountered error loading movie_watch: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Panicf("Encountered error committing transaction: %v", err)
	}

	// Point the database client var at the new database client.
	DB = db
	return m
}

func teardownDatabase(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Panicf("Encountered error tearing down database: %v", err)
	}
}

func TestFindMovieWatch(t *testing.T) {
	m := setupDatabase()
	defer teardownDatabase(m)

	truth := "def-123"
	record := GristMovieWatchRecord{
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
			Uuid:        "",
		},
	}

	uuid, err := FindMovieWatch(&record)
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
