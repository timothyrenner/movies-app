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

type GristMovieWatchRecords struct {
	Records []GristMovieWatchRecord `json:"records"`
}

type GristRecord struct {
	Id int `json:"id"`
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
	FirstTime   bool     `json:"First_Time,omitempty"`
	Watched     int      `json:"Watched,omitempty"`
	JoeBob      bool     `json:"Joe_Bob,omitempty"`
	CallFelissa bool     `json:"Call_Felissa,omitempty"`
	Beast       bool     `json:"Beast,omitempty"`
	Godzilla    bool     `json:"Godzilla,omitempty"`
	Zombies     bool     `json:"Zombies,omitempty"`
	Slasher     bool     `json:"Slasher,omitempty"`
	Service     []string `json:"Service,omitempty"`
	Uuid        string   `json:"uuid,omitempty"`
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

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

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
	request.Header.Set("Content-Type", "application/json")
	_, err = c.client.Do(request)
	if err != nil {
		return fmt.Errorf("encountered error making request: %v", err)
	}

	return nil
}
