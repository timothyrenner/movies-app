package cmd

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func TestGetMovie(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := OmdbClient{
		client:  http.Client{},
		key:     "abc123",
		rootUrl: "http://omdbapi.com",
	}

	httpmock.RegisterResponder("GET", `http://omdbapi.com/?apikey=abc123&i=tt0105347`,
		httpmock.NewStringResponder(
			200,
			`{
			"Title": "Seedpeople",
			"Year": "1992",
			"Rated": "R",
			"Released": "21 Oct 1992",
			"Runtime": "87 min",
			"Genre": "Horror, Sci-Fi",
			"Director": "Peter Manoogian",
			"Writer": "Charles Band, Jackson Barr",
			"Actors": "Sam Hennings, Andrea Roth, Dane Witherspoon",
			"Plot": "The citizens of Comet Valley are being taken over by seeds from an alien plant that has taken root there. A sheriff investigates the strange goings-on.",
			"Language": "English",
			"Country": "United States",
			"Awards": "N/A",
			"Poster": "https://m.media-amazon.com/images/M/MV5BZDE1MTNjZTgtOTkzYS00MDFmLTlmMWYtNDFiYmJiZmVlZGVjXkEyXkFqcGdeQXVyMjAxMjEzNzU@._V1_SX300.jpg",
			"Ratings": [
				{
					"Source": "Internet Movie Database",
					"Value": "4.2/10"
				}
			],
			"Metascore": "N/A",
			"imdbRating": "4.2",
			"imdbVotes": "1,011",
			"imdbID": "tt0105347",
			"Type": "movie",
			"DVD": "22 Jan 2013",
			"BoxOffice": "N/A",
			"Production": "N/A",
			"Website": "N/A",
			"Response": "True"
		}`,
		),
	)

	movie, err := client.GetMovie("tt0105347")
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	truth := &OmdbMovieResponse{
		Title:    "Seedpeople",
		Year:     "1992",
		Rated:    "R",
		Released: "21 Oct 1992",
		Runtime:  "87 min",
		Genre:    "Horror, Sci-Fi",
		Director: "Peter Manoogian",
		Writer:   "Charles Band, Jackson Barr",
		Actors:   "Sam Hennings, Andrea Roth, Dane Witherspoon",
		Plot:     "The citizens of Comet Valley are being taken over by seeds from an alien plant that has taken root there. A sheriff investigates the strange goings-on.",
		Language: "English",
		Country:  "United States",
		Awards:   "N/A",
		Poster:   "https://m.media-amazon.com/images/M/MV5BZDE1MTNjZTgtOTkzYS00MDFmLTlmMWYtNDFiYmJiZmVlZGVjXkEyXkFqcGdeQXVyMjAxMjEzNzU@._V1_SX300.jpg",
		Ratings: []Rating{
			{Source: "Internet Movie Database", Value: "4.2/10"},
		},
		Metascore:  "N/A",
		ImdbRating: "4.2",
		ImdbVotes:  "1,011",
		ImdbID:     "tt0105347",
		Type:       "movie",
		DVD:        "22 Jan 2013",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
	if !cmp.Equal(truth, movie) {
		t.Errorf("Expected %v, got %v", *truth, *movie)
	}
}
