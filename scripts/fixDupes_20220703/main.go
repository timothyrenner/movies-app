package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/timothyrenner/movies-app/cmd"
)

var DB = cmd.DB

var RECORDS_TO_REMOVE []string = []string{
	// Butcher, Baker, Nightmare Maker (2 dupes).
	"06548780-58a4-4a00-8db1-27bb2312b8c2",
	"2ccaa31c-fffb-4bc9-a0bd-1fee54c167eb",
	// Grizzly 2: Revenge (will never ever ever ever watch that again. maybe).
	"deb346d1-8e4b-40f4-bb11-1989c03b2932",
	// Hellbender.
	"c479cef6-5003-40ee-9f90-627ff2b98db6",
	// The Baby. So uncomfortable.
	"d2a5000d-bcf7-40be-a4ae-4e9c2d109f2a",
	// The Found Footage Phenomenon.
	"7c4f16b6-27c5-45ea-9b95-f176ad2d7e28",
	// The Freakmaker.
	"6d004a0f-d16f-4d03-b478-573a716da88f",
	// The Monster Club.
	"f2cf26d3-021b-4507-8aa8-808819e321d8",
	// The Stepfather (who am I here?)
	"a020ef35-06a2-44e6-b3c7-7df83de83ea8",
	// The Thing.
	"0dba56e4-22b3-4cb7-867e-d7d78ffad4d8",
	// Offseason
	"272a6d82-d5d4-40a8-83fd-3b898bb9f06a",
}

func main() {
	dbc, err := sql.Open("sqlite3", DB)
	if err != nil {
		log.Panicf("Error opening db: %v", err)
	}
	defer dbc.Close()
	paramSlice := make([]string, len(RECORDS_TO_REMOVE))
	for ii := range paramSlice {
		paramSlice[ii] = "?"
	}
	params := strings.Join(paramSlice, ",")
	query := fmt.Sprintf(`DELETE FROM movie_watch WHERE uuid IN (%v)`, params)
	_, err = dbc.Exec(
		query,
		// for some reason I can't splat the damn string slice, so we're doing
		// it the "hard" way.
		RECORDS_TO_REMOVE[0],
		RECORDS_TO_REMOVE[1],
		RECORDS_TO_REMOVE[2],
		RECORDS_TO_REMOVE[3],
		RECORDS_TO_REMOVE[4],
		RECORDS_TO_REMOVE[5],
		RECORDS_TO_REMOVE[6],
		RECORDS_TO_REMOVE[7],
		RECORDS_TO_REMOVE[8],
		RECORDS_TO_REMOVE[9],
		RECORDS_TO_REMOVE[10],
	)
	if err != nil {
		log.Panicf("Error deleting records: %v", err)
	}
}
