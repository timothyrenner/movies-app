package cmd

import (
	"database/sql"
	"fmt"
)

func FindMovieWatch(movieWatchRecord *GristMovieWatchRecord) (string, error) {
	query := `
	SELECT
		uuid
	FROM
		movie_watch
	WHERE
		movie_title = $1 AND
		watched = $2
	`

	dbRow := DB.QueryRow(
		query, movieWatchRecord.Fields.Name, movieWatchRecord.Fields.Watched,
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

func FindMovie(movieWatchRecord *GristMovieWatchRecord) (string, error) {
	query := `
	SELECT
		uuid
	FROM
		movie
	WHERE
		title = $1
	`
	dbRow := DB.QueryRow(query, movieWatchRecord.Fields.Name)
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
