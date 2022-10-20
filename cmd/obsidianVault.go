package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/timothyrenner/movies-app/database"
)

type MovieWatchParser struct {
	TitleExtractor       *regexp.Regexp
	WatchedDateExtractor *regexp.Regexp
	ImdbLinkExtractor    *regexp.Regexp
	ImdbIdExtractor      *regexp.Regexp
	ServiceExtractor     *regexp.Regexp
	FirstTimeExtractor   *regexp.Regexp
	JoeBobExtractor      *regexp.Regexp
	SlasherExtractor     *regexp.Regexp
	CallFelissaExtractor *regexp.Regexp
	BeastExtractor       *regexp.Regexp
	ZombiesExtractor     *regexp.Regexp
	GodzillaExtractor    *regexp.Regexp
	WallpaperFuExtractor *regexp.Regexp
	NotesExtractor       *regexp.Regexp
}

func CreateMovieWatchParser() (*MovieWatchParser, error) {
	parser := MovieWatchParser{}
	// Time for some regex fu.

	titleExtractor, err := regexp.Compile(
		`name::\s*\[\[([a-zA-z0-9:\-/' ]+) \(tt\d{7,8}\)\]\]\s*\n`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for title: %v", err)
	}
	parser.TitleExtractor = titleExtractor

	watchedDateExtractor, err := regexp.Compile(
		`watched::\s*\[\[(\d{4}-\d{2}-\d{2})\]\]`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for watched date: %v", err)
	}
	parser.WatchedDateExtractor = watchedDateExtractor

	imdbLinkExtractor, err := regexp.Compile(
		`imdb_link::\s*(https://www\.imdb\.com/title/tt\d{7,8}/)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for imdb link: %v", err)
	}
	parser.ImdbLinkExtractor = imdbLinkExtractor

	imdbIdExtractor, err := regexp.Compile(
		`imdb_id::\s*(tt\d{7,8})`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for imdb id: %v", err)
	}
	parser.ImdbIdExtractor = imdbIdExtractor

	serviceExtractor, err := regexp.Compile(
		`service::\s*([a-zA-Z+ ]+)\s*\n`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for service: %v", err)
	}
	parser.ServiceExtractor = serviceExtractor

	firstTimeExtractor, err := regexp.Compile(
		`first_time::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for first_time: %v", err)
	}
	parser.FirstTimeExtractor = firstTimeExtractor

	joeBobExtractor, err := regexp.Compile(
		`joe_bob::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for joe_bob: %v", err)
	}
	parser.JoeBobExtractor = joeBobExtractor

	slasherExtractor, err := regexp.Compile(
		`slasher::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for slasher: %v", err)
	}
	parser.SlasherExtractor = slasherExtractor

	callFelissaExtractor, err := regexp.Compile(
		`call_felissa::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for call_felissa: %v", err)
	}
	parser.CallFelissaExtractor = callFelissaExtractor

	beastExtractor, err := regexp.Compile(
		`beast::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for beast: %v", err)
	}
	parser.BeastExtractor = beastExtractor

	zombiesExtractor, err := regexp.Compile(
		`zombies::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for zombies: %v", err)
	}
	parser.ZombiesExtractor = zombiesExtractor

	godzillaExtractor, err := regexp.Compile(
		`godzilla::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for godzilla: %v", err)
	}
	parser.GodzillaExtractor = godzillaExtractor

	wallpaperFuExtractor, err := regexp.Compile(
		`wallpaper_fu::\s*(true|false)`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for wallpaper_fu: %v", err)
	}
	parser.WallpaperFuExtractor = wallpaperFuExtractor

	notesExtractor, err := regexp.Compile(`(?s)## Notes(.*)$`)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for notes: %v", err)
	}
	parser.NotesExtractor = notesExtractor

	return &parser, nil
}

func (p *MovieWatchParser) ParsePage(fileName string) (*MovieWatchPage, error) {

	pageText, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v: %v", fileName, err)
	}

	page := MovieWatchPage{}

	movieNameMatch := p.TitleExtractor.FindSubmatch(pageText)
	if len(movieNameMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for movie name: got %v", len(movieNameMatch),
		)
	}
	page.Title = string(movieNameMatch[1])
	page.FileTitle = cleanTitle(page.Title)

	watchMatch := p.WatchedDateExtractor.FindSubmatch(pageText)
	if len(watchMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for watch: got %v", len(watchMatch),
		)
	}
	page.Watched = string(watchMatch[1])

	imdbLinkMatch := p.ImdbLinkExtractor.FindSubmatch(pageText)
	if len(imdbLinkMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for imdb_link: got %v", len(imdbLinkMatch),
		)
	}
	page.ImdbLink = string(imdbLinkMatch[1])

	imdbIdMatch := p.ImdbIdExtractor.FindSubmatch(pageText)
	if len(imdbIdMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for imdb_id: got %v", len(imdbIdMatch),
		)
	}
	page.ImdbId = string(imdbIdMatch[1])

	serviceMatch := p.ServiceExtractor.FindSubmatch(pageText)
	if len(serviceMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for service: got %v", len(serviceMatch),
		)
	}
	page.Service = string(serviceMatch[1])

	firstTimeMatch := p.FirstTimeExtractor.FindSubmatch(pageText)
	if len(firstTimeMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for first_time: got %v", len(firstTimeMatch),
		)
	}
	firstTime, err := strconv.ParseBool(string(firstTimeMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing first time match %v: %v",
			string(firstTimeMatch[1]),
			err,
		)
	}
	page.FirstTime = firstTime

	joeBobMatch := p.JoeBobExtractor.FindSubmatch(pageText)
	if len(joeBobMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for joe_bob, got %v", len(joeBobMatch),
		)
	}
	joeBob, err := strconv.ParseBool(string(joeBobMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing joe_bob match %v: %v",
			string(joeBobMatch[1]),
			err,
		)
	}
	page.JoeBob = joeBob

	slasherMatch := p.SlasherExtractor.FindSubmatch(pageText)
	if len(slasherMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for slasher, got %v", len(slasherMatch),
		)
	}
	slasher, err := strconv.ParseBool(string(slasherMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing slasher match %v: %v",
			string(slasherMatch[1]),
			err,
		)
	}
	page.Slasher = slasher

	callFelissaMatch := p.CallFelissaExtractor.FindSubmatch(pageText)
	if len(callFelissaMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for call_felissa, got %v", len(callFelissaMatch),
		)
	}
	callFelissa, err := strconv.ParseBool(string(callFelissaMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing call_felissa match %v: %v",
			string(callFelissaMatch[1]),
			err,
		)
	}
	page.CallFelissa = callFelissa

	beastMatch := p.BeastExtractor.FindSubmatch(pageText)
	if len(beastMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for beast, got %v", len(beastMatch),
		)
	}
	beast, err := strconv.ParseBool(string(beastMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing beast match %v: %v", string(beastMatch[1]), err,
		)
	}
	page.Beast = beast

	zombiesMatch := p.ZombiesExtractor.FindSubmatch(pageText)
	if len(zombiesMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for zombies, got %v", len(zombiesMatch),
		)
	}
	zombies, err := strconv.ParseBool(string(zombiesMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing zombies match %v: %v", string(zombiesMatch[1]), err,
		)
	}
	page.Zombies = zombies

	godzillaMatch := p.GodzillaExtractor.FindSubmatch(pageText)
	if len(godzillaMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for godzilla, got %v", len(godzillaMatch),
		)
	}
	godzilla, err := strconv.ParseBool(string(godzillaMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing godzilla match %v: %v", string(godzillaMatch[1]), err,
		)
	}
	page.Godzilla = godzilla

	wallpaperFuMatch := p.WallpaperFuExtractor.FindSubmatch(pageText)
	if len(wallpaperFuMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for wallpaper_fu, got %v", len(wallpaperFuMatch),
		)
	}
	wallpaperFu, err := strconv.ParseBool(string(wallpaperFuMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing wallpaper_fu match %v: %v",
			string(wallpaperFuMatch[1]), err,
		)
	}
	page.WallpaperFu = wallpaperFu

	notesMatch := p.NotesExtractor.FindSubmatch(pageText)
	if len(notesMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for notes, got %v", len(notesMatch),
		)
	}
	page.Notes = string(notesMatch[1])

	return &page, nil
}

var MOVIE_WATCH_TEMPLATE = `
# {{.Title}}: {{.Watched}}

## Data
name:: [[{{.FileTitle}} ({{.ImdbId}})]]
watched:: [[{{.Watched}}]]
imdb_link:: {{.ImdbLink}}
imdb_id:: {{.ImdbId}}
service:: {{.Service}}
first_time:: {{.FirstTime}}
joe_bob:: {{.JoeBob}}
slasher:: {{.Slasher}}
call_felissa:: {{.CallFelissa}}
beast:: {{.Beast}}
zombies:: {{.Zombies}}
godzilla:: {{.Godzilla}}
wallpaper_fu:: {{.WallpaperFu}}

## Tags
#movie-watch

## Notes
{{.Notes}}
`

type MovieWatchPage struct {
	Title       string
	FileTitle   string
	Watched     string
	ImdbLink    string
	ImdbId      string
	FirstTime   bool
	JoeBob      bool
	CallFelissa bool
	Beast       bool
	Godzilla    bool
	Zombies     bool
	Slasher     bool
	WallpaperFu bool
	Service     string
	Notes       string
}

func CreateMovieWatchPage(row *database.GetAllMovieWatchesRow) *MovieWatchPage {
	return &MovieWatchPage{
		Title:       row.MovieTitle.String,
		FileTitle:   cleanTitle(row.MovieTitle.String),
		Watched:     row.Watched.String,
		ImdbLink:    row.ImdbLink,
		ImdbId:      row.ImdbID,
		FirstTime:   row.FirstTime != 0,
		JoeBob:      row.JoeBob != 0,
		CallFelissa: row.CallFelissa != 0,
		Beast:       row.Beast != 0,
		Godzilla:    row.Godzilla != 0,
		Zombies:     row.Zombies != 0,
		Slasher:     row.Slasher != 0,
		WallpaperFu: row.WallpaperFu.Bool,
		Service:     row.Service,
		Notes:       row.Notes.String,
	}
}

func (p *MovieWatchPage) CreateRow() *EnrichedMovieWatchRow {
	return &EnrichedMovieWatchRow{
		// The uuids live only in the database.
		// The will not be parsed from the page.
		MovieWatchRow: MovieWatchRow{
			MovieTitle: p.Title,
			ImdbId:     p.ImdbId,
			Watched:    p.Watched,
			Service:    p.Service,
			FirstTime:  p.FirstTime,
			JoeBob:     p.JoeBob,
			Notes:      textToNullString(p.Notes),
		},
		ImdbLink:    p.ImdbLink,
		Slasher:     p.Slasher,
		CallFelissa: p.CallFelissa,
		WallpaperFu: p.WallpaperFu,
		Beast:       p.Beast,
		Godzilla:    p.Godzilla,
		Zombies:     p.Zombies,
	}
}

var MOVIE_TEMPLATE = `
# {{.Title}}
## Data
title:: {{.Title}}
imdb_link:: {{.ImdbLink}}
{{$sep := ""}}
genre:: {{range $elem := .Genres}}{{$sep}}[[{{$elem}}]]{{$sep = ", "}}{{end}}
director:: {{$sep = ""}}{{range $elem := .Directors}}{{$sep}}[[{{$elem}}]]{{$sep = ", "}}{{end}}
actor:: {{$sep = ""}}{{range $elem := .Actors}}{{$sep}}[[{{$elem}}]]{{$sep = ", "}}{{end}}
writer:: {{$sep = ""}}{{range $elem := .Writers}}{{$sep}}[[{{$elem}}]]{{$sep = ", "}}{{end}}
year:: {{.Year}}
rated:: {{.Rating}}
released:: {{.Released}}
runtime_minutes:: {{if .RuntimeMinutes}} {{.RuntimeMinutes}} {{end}}
plot:: {{.Plot}}
country:: {{.Country}}
language:: {{.Language}}
box_office:: {{.BoxOffice}}
production:: {{.Production}}
call_felissa:: {{.CallFelissa}}
slasher:: {{.Slasher}}
zombies:: {{.Zombies}}
beast:: {{.Beast}}
godzilla:: {{.Godzilla}}
wallpaper_fu:: {{.WallpaperFu}}

## Tags
#movie
{{$sep = ""}}{{range $elem := .Genres}}{{$sep}}#{{$elem}}{{$sep = "\n"}}{{end}}
`

type MoviePage struct {
	Title          string
	ImdbLink       string
	Genres         []string
	Directors      []string
	Actors         []string
	Writers        []string
	Year           int
	RuntimeMinutes int
	Rating         string
	Released       string
	Plot           string
	Country        string
	Language       string
	BoxOffice      string
	Production     string
	CallFelissa    bool
	Slasher        bool
	Zombies        bool
	Beast          bool
	Godzilla       bool
	WallpaperFu    bool
}

func CreateMoviePageFromRow(
	row *database.GetMovieRow,
	genres []string,
	directors []string,
	writers []string,
	actors []string,
) *MoviePage {
	return &MoviePage{
		Title:          row.Title,
		ImdbLink:       row.ImdbLink,
		Genres:         genres,
		Directors:      directors,
		Actors:         actors,
		Writers:        writers,
		Year:           int(row.Year),
		Rating:         row.Rated.String,
		Released:       row.Released.String,
		RuntimeMinutes: int(row.RuntimeMinutes.Int64),
		Plot:           row.Plot.String,
		Country:        row.Country.String,
		Language:       row.Language.String,
		BoxOffice:      row.BoxOffice.String,
		Production:     row.Production.String,
		CallFelissa:    row.CallFelissa != 0,
		Slasher:        row.Slasher != 0,
		Zombies:        row.Zombies != 0,
		Beast:          row.Beast != 0,
		Godzilla:       row.Godzilla != 0,
		WallpaperFu:    row.WallpaperFu.Bool,
	}
}

func CreateMoviePage(
	omdbResponse *OmdbMovieResponse, movieWatch *EnrichedMovieWatchRow,
) (*MoviePage, error) {
	year, err := strconv.Atoi(omdbResponse.Year)
	if err != nil {
		return nil, fmt.Errorf(
			"error converting %v to int for year: %v", omdbResponse.Year, err,
		)
	}

	releasedDate, err := ParseReleased(omdbResponse.Released)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing date %v: %v", omdbResponse.Released, err,
		)
	}
	runtime, err := ParseRuntime(omdbResponse.Runtime)
	if err != nil {
		log.Printf("Unable to parse %v, setting to null", omdbResponse.Runtime)
	}

	genres := SplitOnCommaAndTrim(omdbResponse.Genre)

	directors := SplitOnCommaAndTrim(omdbResponse.Director)

	writers := SplitOnCommaAndTrim(omdbResponse.Writer)

	actors := SplitOnCommaAndTrim(omdbResponse.Actors)

	return &MoviePage{
		Title:          omdbResponse.Title,
		ImdbLink:       fmt.Sprintf("https://www.imdb.com/title/%v/", omdbResponse.ImdbID),
		Genres:         genres,
		Directors:      directors,
		Writers:        writers,
		Actors:         actors,
		Year:           year,
		Rating:         omdbResponse.Rated,
		Released:       releasedDate,
		RuntimeMinutes: runtime,
		Plot:           omdbResponse.Plot,
		Country:        omdbResponse.Country,
		Language:       omdbResponse.Language,
		BoxOffice:      omdbResponse.BoxOffice,
		Production:     omdbResponse.Production,
		CallFelissa:    movieWatch.CallFelissa,
		Slasher:        movieWatch.Slasher,
		Zombies:        movieWatch.Zombies,
		Beast:          movieWatch.Beast,
		Godzilla:       movieWatch.Godzilla,
		WallpaperFu:    movieWatch.WallpaperFu,
	}, nil
}

type MovieReviewParser struct {
	MovieTitleExtractor  *regexp.Regexp
	MovieLikedExtractor  *regexp.Regexp
	MovieReviewExtractor *regexp.Regexp
}

func CreateMovieReviewParser() (*MovieReviewParser, error) {
	parser := MovieReviewParser{}

	movieTitleExtractor, err := regexp.Compile(
		`movie::\s*\[\[([a-zA-z0-9:\-/' ]+) \((tt\d{7,8})\)\]\]\s*\n`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for movie: %v", err)
	}
	parser.MovieTitleExtractor = movieTitleExtractor

	movieLikedExtractor, err := regexp.Compile(
		`liked::\s*(true|false)\s*\n`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for liked: %v", err)
	}
	parser.MovieLikedExtractor = movieLikedExtractor

	movieReviewExtractor, err := regexp.Compile(
		`(?s)## Review(.*)$`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for notes: %v", err)
	}
	parser.MovieReviewExtractor = movieReviewExtractor
	return &parser, nil
}

type MovieReviewPage struct {
	MovieTitle string
	ImdbId     string
	Liked      bool
	Review     string
}

func (p *MovieReviewParser) ParseMovieReviewPage(filename string) (
	*MovieReviewPage, error,
) {
	pageText, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v: %v", filename, err)
	}

	page := MovieReviewPage{}

	movieTitleMatch := p.MovieTitleExtractor.FindSubmatch(pageText)
	if matchLen := len(movieTitleMatch); matchLen != 3 {
		return nil, fmt.Errorf(
			"expected 3 matches for movie name, got %v", matchLen,
		)
	}
	page.MovieTitle = string(movieTitleMatch[1])
	page.ImdbId = string(movieTitleMatch[2])

	likedMatch := p.MovieLikedExtractor.FindSubmatch(pageText)
	if matchLen := len(likedMatch); matchLen != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for liked, got %v", matchLen,
		)
	}
	liked, err := strconv.ParseBool(string(likedMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing liked match %v: %v",
			string(likedMatch[1]),
			err,
		)
	}
	page.Liked = liked

	reviewMatch := p.MovieReviewExtractor.FindSubmatch(pageText)
	if matchLen := len(reviewMatch); matchLen != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for review, got %v", matchLen,
		)
	}
	page.Review = string(reviewMatch[1])

	return &page, nil
}

func (p *MovieReviewPage) CreateRow(movieUuid string) *MovieReviewRow {
	return &MovieReviewRow{
		Uuid:       uuid.New().String(),
		MovieUuid:  movieUuid,
		MovieTitle: p.MovieTitle,
		Review:     p.Review,
		Liked:      p.Liked,
	}
}
