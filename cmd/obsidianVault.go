package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/timothyrenner/movies-app/database"
)

type MovieWatchParser struct {
	DataExtractor  *regexp.Regexp
	NotesExtractor *regexp.Regexp
	TitleExtractor *regexp.Regexp
}

func CreateMovieWatchParser() (*MovieWatchParser, error) {
	parser := MovieWatchParser{}
	// Time for some regex fu.
	// But not too much.
	dataExtractor, err := regexp.Compile(`## Data\n((?:.|\n)*)\n## Tags`)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for data: %v", err)
	}
	parser.DataExtractor = dataExtractor

	titleExtractor, err := regexp.Compile(`\[\[([a-zA-Z0-9:\-/()', ]+) \(tt\d{7,8}\)\]\]`)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for title: %v", err)
	}
	parser.TitleExtractor = titleExtractor

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

	movieDataMatch := p.DataExtractor.FindSubmatch(pageText)
	if len(movieDataMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for movie data: got %v", len(movieDataMatch),
		)
	}

	dataLines := strings.Split(string(movieDataMatch[1]), "\n")

	for ii := range dataLines {

		splitLine := strings.Split(dataLines[ii], "::")
		tag := splitLine[0]
		// Usually this will add nothing. On the edge case where there are two
		// colons it will prevent data truncation.
		var data string
		if len(splitLine) > 1 {
			data = strings.Join(splitLine[1:], "::")
			data = strings.TrimSpace(data)
		}

		switch tag {
		case "name":
			titleMatch := p.TitleExtractor.FindSubmatch([]byte(data))
			if len(titleMatch) != 2 {
				return nil, fmt.Errorf(
					"should be one title submatch for %v, got %v",
					data, len(titleMatch),
				)
			}
			page.Title = string(titleMatch[1])
		case "watched":
			watched := strings.Trim(data, "]")
			watched = strings.Trim(watched, "[")
			page.Watched = watched
		case "imdb_link":
			page.ImdbLink = data
		case "imdb_id":
			page.ImdbId = data
		case "service":
			page.Service = data
		case "first_time":
			firstTime, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing first time %v: %v", data, err,
				)
			}
			page.FirstTime = firstTime
		case "joe_bob":
			joeBob, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing joe bob %v: %v", data, err,
				)
			}
			page.JoeBob = joeBob
		case "slasher":
			slasher, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing slasher %v: %v", data, err,
				)
			}
			page.Slasher = slasher
		case "call_felissa":
			callFelissa, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing call felissa %v: %v", data, err,
				)
			}
			page.CallFelissa = callFelissa
		case "beast":
			beast, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing beast %v: %v", data, err,
				)
			}
			page.Beast = beast
		case "zombies":
			zombies, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing zombies %v: %v", data, err,
				)
			}
			page.Zombies = zombies
		case "godzilla":
			godzilla, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing godzilla %v: %v", data, err,
				)
			}
			page.Godzilla = godzilla
		case "wallpaper_fu":
			wallpaperFu, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing wallpaper fu %v: %v", data, err,
				)
			}
			page.WallpaperFu = wallpaperFu
		}
	}

	page.FileTitle = cleanTitle(page.Title)

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
		Title:       row.MovieTitle,
		FileTitle:   cleanTitle(row.MovieTitle),
		Watched:     row.Watched,
		ImdbLink:    row.ImdbLink,
		ImdbId:      row.ImdbID,
		FirstTime:   row.FirstTime != 0,
		JoeBob:      row.JoeBob != 0,
		CallFelissa: row.CallFelissa != 0,
		Beast:       row.Beast != 0,
		Godzilla:    row.Godzilla != 0,
		Zombies:     row.Zombies != 0,
		Slasher:     row.Slasher != 0,
		WallpaperFu: row.WallpaperFu != 0,
		Service:     row.Service,
		Notes:       row.Notes.String,
	}
}

type MovieParser struct {
	DataExtractor   *regexp.Regexp
	ImdbIDExtractor *regexp.Regexp
}

func CreateMovieParser() (*MovieParser, error) {
	parser := MovieParser{}

	dataExtractor, err := regexp.Compile(`## Data\n((?:.|\n)*)\n## Tags`)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for data: %v", err)
	}
	parser.DataExtractor = dataExtractor

	return &parser, nil
}

func (p *MovieParser) ParsePage(fileName string) (*MoviePage, error) {
	pageText, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v: %v", fileName, err)
	}

	page := MoviePage{}

	movieDataMatch := p.DataExtractor.FindSubmatch(pageText)
	if len(movieDataMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for movie data, got %v", len(movieDataMatch),
		)
	}

	dataLines := strings.Split(string(movieDataMatch[1]), "\n")
	for ii := range dataLines {
		splitLine := strings.Split(dataLines[ii], "::")
		tag := splitLine[0]
		// Usually this will add nothing. On the edge case where there are two
		// colons in the data it will prevent data truncation.
		var data string
		if len(splitLine) > 1 {
			data = strings.Join(splitLine[1:], "::")
			data = strings.TrimSpace(data)
		}

		switch tag {
		case "title":
			page.Title = data
		case "imdb_link":
			page.ImdbLink = data
		case "":
			// Do nothing, this is a blank line.
		case "genre":
			page.Genres = SplitOnCommaAndTrim(data)
		case "director":
			page.Directors = SplitOnCommaAndTrim(data)
		case "actor":
			page.Actors = SplitOnCommaAndTrim(data)
		case "writer":
			page.Writers = SplitOnCommaAndTrim(data)
		case "year":
			year, err := strconv.Atoi(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing year %v as int: %v", data, err,
				)
			}
			page.Year = year
		case "rated":
			page.Rating = data
		case "released":
			page.Released = data
		case "runtime_minutes":
			runtimeMinutes, err := strconv.Atoi(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing runtime %v as int: %v", data, err,
				)
			}
			page.RuntimeMinutes = runtimeMinutes
		case "plot":
			page.Plot = data
		case "country":
			page.Country = data
		case "language":
			page.Language = data
		case "box_office":
			page.BoxOffice = data
		case "production":
			page.Production = data
		case "call_felissa":
			callFelissa, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing call felissa %v as bool: %v", data, err,
				)
			}
			page.CallFelissa = callFelissa
		case "slasher":
			slasher, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing slasher %v as bool: %v", data, err,
				)
			}
			page.Slasher = slasher
		case "zombies":
			zombies, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing zombies %v as bool: %v", data, err,
				)
			}
			page.Zombies = zombies
		case "beast":
			beast, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing beast %v as bool: %v", data, err,
				)
			}
			page.Beast = beast
		case "godzilla":
			godzilla, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing godzilla %v as bool: %v", data, err,
				)
			}
			page.Godzilla = godzilla
		case "wallpaper_fu":
			wallpaperFu, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"error parsing wallpaper fu %v as bool: %v", data, err,
				)
			}
			page.WallpaperFu = wallpaperFu
		}
	}

	return &page, nil
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
	row *database.Movie,
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
		WallpaperFu:    row.WallpaperFu != 0,
	}
}

func CreateMoviePage(
	omdbResponse *OmdbMovieResponse, movieWatch *MovieWatchPage,
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
	DataExtractor   *regexp.Regexp
	TitleExtractor  *regexp.Regexp
	ReviewExtractor *regexp.Regexp
}

func CreateMovieReviewParser() (*MovieReviewParser, error) {
	parser := MovieReviewParser{}

	dataExtractor, err := regexp.Compile(
		`# Review:.*\n((?:.|\n)*)\n## Review`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for data: %v", err)
	}
	parser.DataExtractor = dataExtractor

	titleExtractor, err := regexp.Compile(
		`\[\[([a-zA-z0-9:\-/' ]+) \((tt\d{7,8})\)\]\]`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for movie: %v", err)
	}
	parser.TitleExtractor = titleExtractor

	reviewExtractor, err := regexp.Compile(
		`(?s)## Review(.*)$`,
	)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex for notes: %v", err)
	}
	parser.ReviewExtractor = reviewExtractor
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

	reviewDataMatch := p.DataExtractor.FindSubmatch(pageText)
	if len(reviewDataMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for review data, got %v", len(reviewDataMatch),
		)
	}
	dataLines := strings.Split(string(reviewDataMatch[1]), "\n")
	for ii := range dataLines {
		splitLine := strings.Split(dataLines[ii], "::")
		tag := splitLine[0]
		var data string
		if len(splitLine) > 1 {
			data = strings.Join(splitLine[1:], "::")
			data = strings.TrimSpace(data)
		}

		switch tag {
		case "movie":
			titleMatch := p.TitleExtractor.FindSubmatch([]byte(data))
			if len(titleMatch) != 3 {
				return nil, fmt.Errorf(
					"expected %v to give 2 groups, got %v", data, len(titleMatch),
				)
			}
			page.MovieTitle = string(titleMatch[1])
			page.ImdbId = string(titleMatch[2])
		case "liked":
			liked, err := strconv.ParseBool(data)
			if err != nil {
				return nil, fmt.Errorf(
					"unabled to parse liked %v as bool: %v", data, err,
				)
			}
			page.Liked = liked
		}
	}

	reviewMatch := p.ReviewExtractor.FindSubmatch(pageText)
	if matchLen := len(reviewMatch); matchLen != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for review, got %v", matchLen,
		)
	}
	page.Review = string(reviewMatch[1])

	return &page, nil
}
