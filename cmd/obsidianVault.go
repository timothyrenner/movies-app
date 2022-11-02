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

	titleExtractor, err := regexp.Compile(`\[\[([a-zA-Z0-9:\-/() ]+) \(tt\d{7,8}\)\]\]`)
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

type MovieParser struct {
	TitleExtractor          *regexp.Regexp
	ImdbLinkExtractor       *regexp.Regexp
	GenreExtractor          *regexp.Regexp
	DirectorExtractor       *regexp.Regexp
	ActorExtractor          *regexp.Regexp
	WriterExtractor         *regexp.Regexp
	YearExtractor           *regexp.Regexp
	RatedExtractor          *regexp.Regexp
	ReleasedExtractor       *regexp.Regexp
	RuntimeMinutesExtractor *regexp.Regexp
	PlotExtractor           *regexp.Regexp
	CountryExtractor        *regexp.Regexp
	LanguageExtractor       *regexp.Regexp
	BoxOfficeExtractor      *regexp.Regexp
	ProductionExtractor     *regexp.Regexp
	CallFelissaExtractor    *regexp.Regexp
	SlasherExtractor        *regexp.Regexp
	ZombiesExtractor        *regexp.Regexp
	BeastExtractor          *regexp.Regexp
	GodzillaExtractor       *regexp.Regexp
	WallpaperFuExtractor    *regexp.Regexp
}

func CreateMovieParser() (*MovieParser, error) {
	parser := MovieParser{}

	titleExtractor, err := regexp.Compile(`title::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling title extractor: %v", err)
	}
	parser.TitleExtractor = titleExtractor

	imdbLinkExtractor, err := regexp.Compile(`imdb_link::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the imdb link extractor: %v", err)
	}
	parser.ImdbLinkExtractor = imdbLinkExtractor

	genreExtractor, err := regexp.Compile(`genre::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the genre extractor: %v", err)
	}
	parser.GenreExtractor = genreExtractor

	directorExtractor, err := regexp.Compile(`director::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the director extractor: %v", err)
	}
	parser.DirectorExtractor = directorExtractor

	actorExtractor, err := regexp.Compile(`actor::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the actor extractor: %v", err)
	}
	parser.ActorExtractor = actorExtractor

	writerExtractor, err := regexp.Compile(`writer::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the write extractor: %v", err)
	}
	parser.WriterExtractor = writerExtractor

	yearExtractor, err := regexp.Compile(`year::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the year extractor: %v", err)
	}
	parser.YearExtractor = yearExtractor

	ratedExtractor, err := regexp.Compile(`rated::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the rated extractor: %v", err)
	}
	parser.RatedExtractor = ratedExtractor

	releasedExtractor, err := regexp.Compile(`released::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the released extractor: %v", err)
	}
	parser.ReleasedExtractor = releasedExtractor

	runtimeMinutesExtractor, err := regexp.Compile(`runtime_minutes::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the runtime minutes extractor: %v", err)
	}
	parser.RuntimeMinutesExtractor = runtimeMinutesExtractor

	plotExtractor, err := regexp.Compile(`plot::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the plot extractor: %v", err)
	}
	parser.PlotExtractor = plotExtractor

	countryExtractor, err := regexp.Compile(`country::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the country extractor: %v", err)
	}
	parser.CountryExtractor = countryExtractor

	languageExtractor, err := regexp.Compile(`language::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the language extractor: %v", err)
	}
	parser.LanguageExtractor = languageExtractor

	boxOfficeExtractor, err := regexp.Compile(`box_office::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the box office extractor: %v", err)
	}
	parser.BoxOfficeExtractor = boxOfficeExtractor

	productionExtractor, err := regexp.Compile(`production::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the production extractor: %v", err)
	}
	parser.ProductionExtractor = productionExtractor

	callFelissaExtractor, err := regexp.Compile(`call_felissa::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the call felissa extractor: %v", err)
	}
	parser.CallFelissaExtractor = callFelissaExtractor

	slasherExtractor, err := regexp.Compile(`slasher::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the slasher extractor: %v", err)
	}
	parser.SlasherExtractor = slasherExtractor

	zombiesExtractor, err := regexp.Compile(`zombies::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the zombies extractor: %v", err)
	}
	parser.ZombiesExtractor = zombiesExtractor

	beastExtractor, err := regexp.Compile(`beast::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the beast extractor: %v", err)
	}
	parser.BeastExtractor = beastExtractor

	godzillaExtractor, err := regexp.Compile(`godzilla::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the godzilla extractor: %v", err)
	}
	parser.GodzillaExtractor = godzillaExtractor

	wallpaperFuExtractor, err := regexp.Compile(`wallpaper_fu::\s*(.*)\s*\n`)
	if err != nil {
		return nil, fmt.Errorf("error compiling the wallpaper fu extractor: %v", err)
	}
	parser.WallpaperFuExtractor = wallpaperFuExtractor

	return &parser, nil
}

func (p *MovieParser) ParsePage(fileName string) (*MoviePage, error) {
	pageText, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v: %v", fileName, err)
	}

	page := MoviePage{}

	movieTitleMatch := p.TitleExtractor.FindSubmatch(pageText)
	if len(movieTitleMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for movie title, got %v", len(movieTitleMatch),
		)
	}
	page.Title = strings.TrimSpace(string(movieTitleMatch[1]))

	imdbLinkMatch := p.ImdbLinkExtractor.FindSubmatch(pageText)
	if len(imdbLinkMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for imdb link, got %v", len(imdbLinkMatch),
		)
	}
	page.ImdbLink = strings.TrimSpace(string(imdbLinkMatch[1]))

	genreMatch := p.GenreExtractor.FindSubmatch(pageText)
	if len(genreMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for genre, got %v", len(genreMatch),
		)
	}
	page.Genres = SplitOnCommaAndTrim(string(genreMatch[1]))

	directorMatch := p.DirectorExtractor.FindSubmatch(pageText)
	if len(directorMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for director, got %v", len(directorMatch),
		)
	}
	page.Directors = SplitOnCommaAndTrim(string(directorMatch[1]))

	actorMatch := p.ActorExtractor.FindSubmatch(pageText)
	if len(actorMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for actor, got %v", len(actorMatch),
		)
	}
	page.Actors = SplitOnCommaAndTrim(string(actorMatch[1]))

	writerMatch := p.WriterExtractor.FindSubmatch(pageText)
	if len(writerMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for writer, got %v", len(writerMatch),
		)
	}
	page.Writers = SplitOnCommaAndTrim(string(writerMatch[1]))

	yearMatch := p.YearExtractor.FindSubmatch(pageText)
	if len(yearMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for year, got %v", len(yearMatch),
		)
	}
	year, err := strconv.Atoi(string(yearMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error converting %v to int: %v", string(yearMatch[1]), err,
		)
	}
	page.Year = year

	runtimeMinutesMatch := p.RuntimeMinutesExtractor.FindSubmatch(pageText)
	if len(runtimeMinutesMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for runtime minutes, got %v",
			len(runtimeMinutesMatch),
		)
	}
	runtime, err := ParseRuntime(string(runtimeMinutesMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing runtime %v: %v", string(runtimeMinutesMatch[1]), err,
		)
	}
	page.RuntimeMinutes = runtime

	ratingMatch := p.RatedExtractor.FindSubmatch(pageText)
	if len(ratingMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for rating, got %v", err,
		)
	}
	page.Rating = string(ratingMatch[1])

	releasedMatch := p.ReleasedExtractor.FindSubmatch(pageText)
	if len(releasedMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for released, got %v", err,
		)
	}
	page.Released = string(releasedMatch[1])

	plotMatch := p.PlotExtractor.FindSubmatch(pageText)
	if len(plotMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for plot, got %v", err,
		)
	}
	page.Plot = string(plotMatch[1])

	countryMatch := p.CountryExtractor.FindSubmatch(pageText)
	if len(countryMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for country, got %v", err,
		)
	}
	page.Country = string(countryMatch[1])

	languageMatch := p.LanguageExtractor.FindSubmatch(pageText)
	if len(languageMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for language, got %v", err,
		)
	}
	page.Language = string(languageMatch[1])

	boxOfficeMatch := p.BoxOfficeExtractor.FindSubmatch(pageText)
	if len(boxOfficeMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for box office, got %v", err,
		)
	}
	page.BoxOffice = string(boxOfficeMatch[1])

	productionMatch := p.ProductionExtractor.FindSubmatch(pageText)
	if len(productionMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for production, got %v", err,
		)
	}
	page.Production = string(productionMatch[1])

	callFelissaMatch := p.CallFelissaExtractor.FindSubmatch(pageText)
	if len(callFelissaMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for call felissa, got %v", err,
		)
	}
	callFelissa, err := strconv.ParseBool(string(callFelissaMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(callFelissaMatch[1]), err,
		)
	}
	page.CallFelissa = callFelissa

	slasherMatch := p.SlasherExtractor.FindSubmatch(pageText)
	if len(slasherMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for slasher, got %v", err,
		)
	}
	slasher, err := strconv.ParseBool(string(slasherMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(slasherMatch[1]), err,
		)
	}
	page.Slasher = slasher

	zombiesMatch := p.ZombiesExtractor.FindSubmatch(pageText)
	if len(zombiesMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for zombies, got %v", err,
		)
	}
	zombies, err := strconv.ParseBool(string(zombiesMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(zombiesMatch[1]), err,
		)
	}
	page.Zombies = zombies

	beastMatch := p.BeastExtractor.FindSubmatch(pageText)
	if len(beastMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for beast, got %v", err,
		)
	}
	beast, err := strconv.ParseBool(string(beastMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(beastMatch[1]), err,
		)
	}
	page.Beast = beast

	godzillaMatch := p.GodzillaExtractor.FindSubmatch(pageText)
	if len(godzillaMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for godzilla, got %v", err,
		)
	}
	godzilla, err := strconv.ParseBool(string(godzillaMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(godzillaMatch[1]), err,
		)
	}
	page.Godzilla = godzilla

	wallpaperFuMatch := p.WallpaperFuExtractor.FindSubmatch(pageText)
	if len(wallpaperFuMatch) != 2 {
		return nil, fmt.Errorf(
			"expected 2 matches for wallpaper_fu, got %v", err,
		)
	}
	wallpaperFu, err := strconv.ParseBool(string(wallpaperFuMatch[1]))
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing %v as bool: %v", string(wallpaperFuMatch[1]), err,
		)
	}
	page.WallpaperFu = wallpaperFu

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
		WallpaperFu:    row.WallpaperFu.Bool,
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
