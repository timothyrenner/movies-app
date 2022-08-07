package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type OmdbClient struct {
	client  http.Client
	key     string
	rootUrl string
}

type OmdbMovieResponse struct {
	Title      string   `json:"Title"`
	Year       string   `json:"Year"`
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Poster     string   `json:"Poster"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbVotes  string   `json:"imdbVotes"`
	ImdbID     string   `json:"imdbID"`
	Type       string   `json:"Type"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`
	Response   string   `json:"Response"`
}

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

func NewOmdbClient(key string) *OmdbClient {
	client := OmdbClient{}
	client.rootUrl = "http://omdbapi.com"
	client.key = key
	client.client = http.Client{
		Timeout: 5 * time.Second,
	}
	return &client
}

func (c *OmdbClient) GetMovie(movieId string) (*OmdbMovieResponse, error) {
	params := url.Values{}
	params.Set("apikey", c.key)
	params.Set("i", movieId)

	url := fmt.Sprintf("%v/?%v", c.rootUrl, params.Encode())

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var movieResponse OmdbMovieResponse
	if err = json.Unmarshal(body, &movieResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	return &movieResponse, nil
}
