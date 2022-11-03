package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetMovieTitleFromWatchFile(t *testing.T) {
	fileContents := []byte(
		`
# Tenebrae: 2022-05-26

## Data
name:: [[Tenebrae]]
watched:: [[2022-05-27]]
imdb_id:: tt0084777
first_time:: false
joe_bob:: true
service:: Shudder

## Tags
#movie-watch
		`,
	)
	truth := "2022-05-27"
	answer, err := GetMovieTitleFromWatchFile(fileContents)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}

func TestGetImdbIdFromWatchFile(t *testing.T) {
	fileContents := []byte(
		`
# Tenebrae: 2022-05-26

## Data
name:: [[Tenebrae]]
watched:: [[2022-05-27]]
imdb_id:: tt0084777
first_time:: false
joe_bob:: true
service:: Shudder

## Tags
#movie-watch
		`,
	)
	truth := "tt0084777"
	answer, err := GetMovieImdbIdFromWatchFile(fileContents)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected %v, got %v", truth, answer)
	}
}
