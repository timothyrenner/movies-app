package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GristClient struct {
	client  http.Client
	header  http.Header
	rootUrl string
}

type GristRecords struct {
	Records []GristRecord `json:"records"`
}

type GristRecord struct {
	Id int `json:"id,omitempty"`
}
type GristMovieWatchRecords struct {
	Records []GristMovieWatchRecord `json:"records"`
}

type GristMovieWatchRecord struct {
	GristRecord
	Fields GristMovieWatchFields `json:"fields"`
}

func (r *GristMovieWatchRecord) ImdbId() string {
	urlComponents := strings.Split(r.Fields.ImdbLink, "/")
	if urlComponents[len(urlComponents)-1] != "" {
		log.Panicf("Encountered error getting ID from %v", r.Fields.ImdbLink)
	}
	return urlComponents[len(urlComponents)-2]
}

type GristMovieWatchFields struct {
	Name        string   `json:"name,omitempty"`
	ImdbLink    string   `json:"IMDB_Link,omitempty"`
	ImdbId      string   `json:"IMDB_ID,omitempty"`
	FirstTime   bool     `json:"First_Time,omitempty"`
	Watched     int      `json:"Watched,omitempty"`
	JoeBob      bool     `json:"Joe_Bob,omitempty"`
	CallFelissa bool     `json:"Call_Felissa,omitempty"`
	Beast       bool     `json:"Beast,omitempty"`
	Godzilla    bool     `json:"Godzilla,omitempty"`
	Zombies     bool     `json:"Zombies,omitempty"`
	Slasher     bool     `json:"Slasher,omitempty"`
	Service     []string `json:"Service,omitempty"`
	Movie       int64    `json:"Movie,omitempty"`
}

type GristMovieRecords struct {
	Records []GristMovieRecord `json:"records"`
}

type GristMovieRecord struct {
	GristRecord
	Fields GristMovieFields `json:"fields"`
}

type GristMovieFields struct {
	Title       string   `json:"Title,omitempty"`
	ImdbLink    string   `json:"IMDB_Link,omitempty"`
	Year        int      `json:"Year,omitempty"`
	Rated       string   `json:"Rated,omitempty"`
	Released    string   `json:"Released,omitempty"`
	Runtime     int      `json:"Runtime,omitempty"`
	Plot        string   `json:"Plot,omitempty"`
	Country     string   `json:"Country,omitempty"`
	Language    string   `json:"Language,omitempty"`
	BoxOffice   string   `json:"BoxOffice,omitempty"`
	Production  string   `json:"Production,omitempty"`
	CallFelissa bool     `json:"Call_Felissa,omitempty"`
	Slasher     bool     `json:"Slasher,omitempty"`
	Zombies     bool     `json:"Zombies,omitempty"`
	Beast       bool     `json:"Beast,omitempty"`
	Godzilla    bool     `json:"Godzilla,omitempty"`
	Genre       []string `json:"Genre,omitempty"`
	Actor       []string `json:"Actor,omitempty"`
	Director    []string `json:"Director,omitempty"`
	Writer      []string `json:"Writer,omitempty"`
	Rating      []any    `json:"Rating,omitempty"`
}

type GristMovieRatingRecords struct {
	Records []GristMovieRatingRecord `json:"records"`
}

type GristMovieRatingRecord struct {
	GristRecord
	Fields GristMovieRatingFields `json:"fields"`
}

type GristMovieRatingFields struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

func NewGristClient(key string) *GristClient {
	client := GristClient{}
	client.header = http.Header{}
	client.header.Add("Authorization", fmt.Sprintf("Bearer %v", key))
	client.rootUrl = "https://docs.getgrist.com/api"
	client.client = http.Client{
		Timeout: 5 * time.Second,
	}
	return &client
}

func (c *GristClient) GetMovieWatchRecords(
	documentId string,
	tableId string,
	filter *map[string]any,
	sort string,
	limit int,
) (*GristMovieWatchRecords, error) {
	params := url.Values{}
	if filter != nil {
		bytes, err := json.Marshal(*filter)
		if err != nil {
			return nil, fmt.Errorf("error serializing filter: %v", err)
		}
		params.Set("filter", string(bytes))
	}

	if sort != "" {
		params.Set("sort", sort)
	}

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	url := fmt.Sprintf(
		"%v/docs/%v/tables/%v/records?%v",
		c.rootUrl,
		documentId,
		tableId,
		params.Encode(),
	)

	log.Printf("Making call to grist: %v", url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header = c.header

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var movieWatchResponse GristMovieWatchRecords
	if err = json.Unmarshal(body, &movieWatchResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &movieWatchResponse, nil
}

func (c *GristClient) UpdateMovieWatchRecords(
	documentId string,
	tableId string,
	records *GristMovieWatchRecords,
) error {
	url := fmt.Sprintf(
		"%v/docs/%v/tables/%v/records",
		c.rootUrl,
		documentId,
		tableId,
	)

	payload, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("encountered error marshalling records: %v", err)
	}
	request, err := http.NewRequest(
		http.MethodPatch, url, bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("encountered error creating request: %v", err)
	}
	request.Header = c.header
	request.Header.Set("Content-Type", "application/json")
	_, err = c.client.Do(request)
	if err != nil {
		return fmt.Errorf("encountered error making request: %v", err)
	}

	return nil
}

func (c *GristClient) CreateMovieRatingRecords(
	documentId string,
	tableId string,
	records *GristMovieRatingRecords,
) (*GristRecords, error) {
	url := fmt.Sprintf(
		"%v/docs/%v/tables/%v/records",
		c.rootUrl,
		documentId,
		tableId,
	)
	payload, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("encountered error marshalling records: %v", err)
	}
	request, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("encountered error creating request: %v", err)
	}
	request.Header = c.header
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("encountered error making request: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var movieRatingResponse GristRecords
	if err = json.Unmarshal(body, &movieRatingResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &movieRatingResponse, nil
}

func (c *GristClient) CreateMovieRecords(
	documentId string,
	tableId string,
	records *GristMovieRecords,
) (*GristRecords, error) {
	url := fmt.Sprintf(
		"%v/docs/%v/tables/%v/records",
		c.rootUrl,
		documentId,
		tableId,
	)
	payload, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("encountered error marshalling records: %v", err)
	}
	request, err := http.NewRequest(
		http.MethodPost, url, bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header = c.header
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("encountered error making request: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(
			"got status code %v: %v", response.StatusCode, string(body),
		)
	}
	var movieResponse GristRecords
	if err = json.Unmarshal(body, &movieResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &movieResponse, nil

}
