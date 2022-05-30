package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GristClient struct {
	client  http.Client
	header  http.Header
	rootUrl string
}

type GetMovieWatchRecordsResponse struct {
	Records []GristMovieWatchRecord `json:"records"`
}

type GristRecord struct {
	Id int `json:"id"`
}

type GristMovieWatchRecord struct {
	GristRecord
	Fields GristMovieWatchFields `json:"fields"`
}

type GristMovieWatchFields struct {
	Name        string   `json:"name"`
	ImdbLink    string   `json:"IMDB_Link"`
	FirstTime   bool     `json:"First_Time"`
	Watched     int      `json:"Watched"`
	JoeBob      bool     `json:"Joe_Bob"`
	CallFelissa bool     `json:"Call_Felissa"`
	Beast       bool     `json:"Beast"`
	Godzilla    bool     `json:"Godzilla"`
	Zombies     bool     `json:"Zombies"`
	Slasher     bool     `json:"Slasher"`
	Service     []string `json:"Service"`
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

func (c *GristClient) GetRecords(
	documentId string,
	tableId string,
	filter *map[string]any,
	sort string,
	limit int,
) (*GetMovieWatchRecordsResponse, error) {
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

	var movieWatchResponse GetMovieWatchRecordsResponse
	if err = json.Unmarshal(body, &movieWatchResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &movieWatchResponse, nil
}
