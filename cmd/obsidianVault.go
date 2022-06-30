package cmd

import "time"

var MOVIE_WATCH_TEMPLATE = `
# {{.Title}}: {{.Watched}}

## Data
name:: [[{{.Title}}]]
watched:: [[{{.Watched}}]]
imdb_id:: {{.ImdbId}}
first_time:: {{.FirstTime}}
joe_bob:: {{.JoeBob}}
service:: {{.Service}}

## Tags
#movie-watch
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
	Service     string
}

func (r *MovieWatchRow) CreatePage() *MovieWatchPage {
	// ! There's a time zone issue here. Whatever Grist is returning is not
	// ! GMT, it's local time. The watch date for every movie watch is one day
	// ! off.
	// ! Another option is to just switch the column from an int to a formatted
	// ! string and correct it in a script. That needs to be done anyway, why
	// ! not do it now?

	// ! The code has been adjusted, theoretically should work.
	watched := time.Unix(int64(r.Watched)+5*60*60, 0).Format("2006-01-02")
	return &MovieWatchPage{
		Title:     r.MovieTitle,
		FileTitle: cleanTitle(r.MovieTitle),
		Watched:   watched,
		ImdbId:    r.ImdbId,
		FirstTime: r.FirstTime,
		JoeBob:    r.JoeBob,
		Service:   r.Service,
	}
}

func (p *MovieWatchPage) CreateRow() *MovieWatchRow {
	return nil // TODO: Implement.
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
runtime_minutes:: {{.RuntimeMinutes}}
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
}

func (r *MovieRow) createPage(
	genres []string, directors []string, writers []string, actors []string,
) *MoviePage {
	return &MoviePage{
		Title:          r.Title,
		ImdbLink:       r.ImdbLink,
		Genres:         genres,
		Directors:      directors,
		Actors:         actors,
		Writers:        writers,
		Year:           r.Year,
		Rating:         r.Rated.String,
		Released:       r.Released.String,
		RuntimeMinutes: r.RuntimeMinutes,
		Plot:           r.Plot.String,
		Country:        r.Country.String,
		Language:       r.Language.String,
		BoxOffice:      r.BoxOffice.String,
		Production:     r.Production.String,
		CallFelissa:    r.CallFelissa,
		Slasher:        r.Slasher,
		Zombies:        r.Zombies,
		Beast:          r.Beast,
		Godzilla:       r.Godzilla,
	}
}
