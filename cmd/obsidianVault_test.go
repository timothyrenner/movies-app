package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestWatchPage() *os.File {
	file, err := os.CreateTemp(".", "test_watch")
	if err != nil {
		log.Panicf("Error creating temp file: %v", err)
	}

	_, err = file.WriteString(
		`
# Uncle Sam: 2022-07-01

## Data
name:: [[Uncle Sam (tt0118025)]]
watched:: [[2022-07-01]]
imdb_link:: https://www.imdb.com/title/tt0118025/
imdb_id:: tt0118025
service:: Shudder
first_time:: true
joe_bob:: true
slasher:: false
call_felissa:: false
beast:: true
zombies:: false
godzilla:: false
wallpaper_fu:: true

## Tags
#movie-watch

## Notes
"Don't be afraid, it's only friendly fire"
"I must be batting 750 with the bereaved" - army dude who notifies widows
"Must be awful lonely being dead"
Prevert uncle sam on stilts
Literally the worst national anthem rendition in existence.
A fuckin sack race half marathon obstacle course
		`,
	)
	if err != nil {
		os.Remove(file.Name())
		log.Panicf("error writing contents: %v", err)
	}
	return file
}

func createTestMoviePage() *os.File {
	file, err := os.CreateTemp(".", "test_movie")
	if err != nil {
		log.Panicf("Error creating test movie page: %v", err)
	}
	_, err = file.WriteString(`
# XTRO
## Data
title:: XTRO
imdb_link:: https://www.imdb.com/title/tt0086610/

genre:: [[Horror]], [[Sci-Fi]]
director:: [[Harry Bromley Davenport]]
actor:: [[Philip Sayer]], [[Bernice Stegers]], [[Danny Brainin]]
writer:: [[Harry Bromley Davenport]], [[Iain Cassie]], [[Michel Parry]]
year:: 1982
rated:: R
released:: 1983-01-07
runtime_minutes:: 84
plot:: An alien creature impregnates a woman who gives birth to a man that was abducted by aliens three years ago. The man reconnects with his wife and son for a sinister purpose.
country:: United Kingdom
language:: English
box_office:: 
production:: 
call_felissa:: true
slasher:: false
zombies:: false
beast:: true
godzilla:: false
wallpaper_fu:: false

## Tags
#movie
#Horror
#Sci-Fi
`)
	if err != nil {
		os.Remove(file.Name())
		log.Panicf("error writing review_contents: %v", err)
	}
	return file
}

func createTestMovieReviewPage() *os.File {
	file, err := os.CreateTemp(".", "test_review")
	if err != nil {
		log.Panicf("Error creating test review page: %v", err)
	}

	_, err = file.WriteString(`
# Review: Uncle Sam
movie::[[Uncle Sam (tt0118025)]]
liked::true

## Review
Look there's a zombie soldier dressed as Uncle Sam who blows people up with fireworks.
Do you really want anything more in a movie?
Oh you want a prevert peeping Tom not-zombie Uncle Sam on stilts too?
We got you.
		`)

	if err != nil {
		os.Remove(file.Name())
		log.Panicf("error writing review contents: %v", err)
	}
	return file
}

func TestParseWatchPage(t *testing.T) {
	file := createTestWatchPage()
	defer file.Close()
	defer os.Remove(file.Name())

	parser, err := CreateMovieWatchParser()
	if err != nil {
		t.Errorf("Error creating parser: %v", err)
	}

	answer, err := parser.ParsePage(file.Name())
	if err != nil {
		t.Errorf("Error parsing page: %v", err)
	}

	truth := MovieWatchPage{
		Title:       "Uncle Sam",
		FileTitle:   "Uncle Sam",
		Watched:     "2022-07-01",
		ImdbLink:    "https://www.imdb.com/title/tt0118025/",
		ImdbId:      "tt0118025",
		Service:     "Shudder",
		FirstTime:   true,
		JoeBob:      true,
		CallFelissa: false,
		Beast:       true,
		Zombies:     false,
		Godzilla:    false,
		WallpaperFu: true,
		Notes: `
"Don't be afraid, it's only friendly fire"
"I must be batting 750 with the bereaved" - army dude who notifies widows
"Must be awful lonely being dead"
Prevert uncle sam on stilts
Literally the worst national anthem rendition in existence.
A fuckin sack race half marathon obstacle course
		`,
	}
	if !cmp.Equal(truth, *answer) {
		t.Errorf("expected \n%v, got \n%v", truth, *answer)
	}
}

func TestParseMoviePage(t *testing.T) {
	file := createTestMoviePage()
	defer file.Close()
	defer os.Remove(file.Name())

	truth := &MoviePage{
		Title:          "XTRO",
		ImdbLink:       "https://www.imdb.com/title/tt0086610/",
		Genres:         []string{"Horror", "Sci-Fi"},
		Directors:      []string{"Harry Bromley Davenport"},
		Actors:         []string{"Philip Sayer", "Bernice Stegers", "Danny Brainin"},
		Writers:        []string{"Harry Bromley Davenport", "Iain Cassie", "Michel Parry"},
		Year:           1982,
		Rating:         "R",
		Released:       "1983-01-07",
		RuntimeMinutes: 84,
		Plot:           "An alien creature impregnates a woman who gives birth to a man that was abducted by aliens three years ago. The man reconnects with his wife and son for a sinister purpose.",
		Country:        "United Kingdom",
		Language:       "English",
		BoxOffice:      "",
		Production:     "",
		CallFelissa:    true,
		Slasher:        false,
		Zombies:        false,
		Beast:          true,
		Godzilla:       false,
		WallpaperFu:    false,
	}

	parser, err := CreateMovieParser()
	if err != nil {
		t.Errorf("error creating parser: %v", err)
	}
	answer, err := parser.ParsePage(file.Name())
	if err != nil {
		t.Errorf("error parsing page: %v", err)
	}
	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, answer)
	}
}

func TestCreateMoviePage(t *testing.T) {
	omdbResponse := omdbSampleMovie()
	movieWatch := &MovieWatchPage{
		Title:       "Tenebrae",
		ImdbId:      "tt0084777",
		Watched:     "2022-05-27",
		Service:     "Shudder",
		FirstTime:   false,
		JoeBob:      true,
		CallFelissa: false,
		Beast:       false,
		Godzilla:    false,
		Zombies:     false,
		Slasher:     true,
		WallpaperFu: false,
	}

	truth := &MoviePage{
		Title:          "Tenebrae",
		ImdbLink:       "https://www.imdb.com/title/tt0084777/",
		Genres:         []string{"Horror", "Mystery", "Thriller"},
		Directors:      []string{"Dario Argento"},
		Actors:         []string{"Anthony Franciosa", "Giuliano Gemma", "John Saxon"},
		Writers:        []string{"Dario Argento"},
		Year:           1982,
		RuntimeMinutes: 101,
		Rating:         "R",
		Released:       "1984-02-17",
		Plot:           "An American writer in Rome is stalked and harassed by a serial killer who is murdering everyone associated with his work on his latest book.",
		Country:        "Italy",
		Language:       "Italian, Spanish",
		BoxOffice:      "N/A",
		Production:     "N/A",
		CallFelissa:    false,
		Slasher:        true,
		Zombies:        false,
		Beast:          false,
		Godzilla:       false,
		WallpaperFu:    false,
	}

	answer, err := CreateMoviePage(omdbResponse, movieWatch)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}

	if !cmp.Equal(truth, answer) {
		t.Errorf("Expected \n%v, got \n%v", *truth, *answer)
	}
}

func TestParseMovieReviewPage(t *testing.T) {
	file := createTestMovieReviewPage()
	defer file.Close()
	defer os.Remove(file.Name())

	parser, err := CreateMovieReviewParser()
	if err != nil {
		t.Errorf("Error creating parser: %v", err)
	}

	answer, err := parser.ParseMovieReviewPage(file.Name())
	if err != nil {
		t.Errorf("Error parsing page: %v", err)
	}
	truth := MovieReviewPage{
		MovieTitle: "Uncle Sam",
		ImdbId:     "tt0118025",
		Liked:      true,
		Review: `
Look there's a zombie soldier dressed as Uncle Sam who blows people up with fireworks.
Do you really want anything more in a movie?
Oh you want a prevert peeping Tom not-zombie Uncle Sam on stilts too?
We got you.
		`,
	}
	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}
}
