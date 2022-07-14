package main

import (
	"database/sql"
	"log"

	"github.com/timothyrenner/movies-app/cmd"
)

var DB = cmd.DB

/* QUERY
-- The busted imdb IDs are only in the movie_watch table, not the movie table.
SELECT
	w.uuid, w.movie_title, w.watched, m.imdb_id
FROM movie_watch AS w
JOIN movie AS m
	ON w.movie_uuid = m.uuid
WHERE w.imdb_id='';
*/

var RECORDS_TO_REPAIR map[string]string = map[string]string{
	// Seedpeople.
	"0f789d68-8b6f-4269-a56d-87f8a3a2835d": "tt0105347",
	// The Found Footage Phenomenon.
	"fc905f75-007f-48a6-a432-a38d874d75fd": "tt12882656",
	// The Monster Club.
	"c254c02a-35da-41f2-8226-3502c6bfa918": "tt0081178",
	// Hellbender.
	"41db4010-cc33-4306-a65c-b597987d7a32": "tt14905650",
	// The Baby.
	"063be892-9cd1-46e7-9f4b-830337b66917": "tt0069754",
	// Butcher, Baker, Nightmare Maker.
	"d66d095b-feae-4fb2-b758-2501d75e5a51": "tt0082813",
	// Grizzly 2: Revenge
	"dfb88032-de3a-4316-a067-ed1ef4bed9a6": "tt0093119",
	// The Stepfather
	"9a7bb097-41f2-4e76-9054-1449fcaacf42": "tt0094035",
	// The Freakmaker
	"11ecded0-eb63-44db-ba6b-49a1f5a07730": "tt0070423",
	// Offseason
	"f9956a84-60f3-420b-8308-b38c6e53450d": "tt11454320",
	// The Thing
	"688ccbcf-bb01-4081-88c1-58e8185d5bfc": "tt0084787",
	// Destroyer
	"b2a5da42-4f1c-47f2-ad9c-76ea6cd86cc3": "tt0095005",
	// Doctor Strange in the Multiverse of Madness
	"a3f021b4-71ff-46c0-8137-ac0a7377406c": "tt9419884",
	// Head of the Family
	"0552eedc-81dd-4bd6-bb7c-e6f9a52e3069": "tt0116503",
	// Habit
	"a7e92e8f-aaf0-45d3-a3f6-bfbf67272af3": "tt0113241",
}

func main() {
	dbc, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening db: %v", err)
	}
	defer dbc.Close()

	updateMovieWatchStatement, err := dbc.Prepare(`
		UPDATE movie_watch SET imdb_id = ? WHERE uuid = ?
	`)
	if err != nil {
		log.Panicf(
			"Error creating movie watch update prepared statement: %v", err,
		)
	}

	for watchUuid, movieImdbId := range RECORDS_TO_REPAIR {
		log.Printf("Updating %v with %v", watchUuid, movieImdbId)
		// We probably don't need a transaction here but whatever.
		tx, err := dbc.Begin()
		if err != nil {
			log.Panicf("Error beginning transaction: %v", err)
		}
		defer tx.Rollback()
		_, err = tx.Stmt(updateMovieWatchStatement).Exec(
			movieImdbId, watchUuid,
		)
		if err != nil {
			log.Panicf("Encountered error updating movie watch row: %v", err)
		}
		if err = tx.Commit(); err != nil {
			log.Panicf("Error committing transaction: %v", err)
		}
	}
	log.Println("Done.")
}
