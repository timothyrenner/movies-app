package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestWatchPage() *os.File {
	file, err := ioutil.TempFile(".", "test_watch")
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
		log.Panicf("error writing contents: %v", err)
		os.Remove(file.Name())
	}
	return file
}

func TestParsePage(t *testing.T) {
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

func TestCreateRow(t *testing.T) {
	page := MovieWatchPage{
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

	truth := EnrichedMovieWatchRow{
		MovieWatchRow: MovieWatchRow{
			MovieTitle: "Uncle Sam",
			Watched:    "2022-07-01",
			ImdbId:     "tt0118025",
			Service:    "Shudder",
			FirstTime:  true,
			JoeBob:     true,
			Notes: `
"Don't be afraid, it's only friendly fire"
"I must be batting 750 with the bereaved" - army dude who notifies widows
"Must be awful lonely being dead"
Prevert uncle sam on stilts
Literally the worst national anthem rendition in existence.
A fuckin sack race half marathon obstacle course
		`,
		},
		CallFelissa: false,
		Beast:       true,
		Zombies:     false,
		Godzilla:    false,
		WallpaperFu: true,
	}

	answer := page.CreateRow()
	if !cmp.Equal(truth, *answer) {
		t.Errorf("Expected \n%v, got \n%v", truth, *answer)
	}
}

func TestCreateMoviePage(t *testing.T) {
	omdbResponse := omdbSampleMovie()
	movieWatch := &EnrichedMovieWatchRow{
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
