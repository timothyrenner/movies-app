package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestWatchPage() *os.File {
	// TODO: Implement.
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
		Title:       "Uncle Sam (tt0118025)",
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
