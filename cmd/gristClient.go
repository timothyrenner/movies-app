package cmd

import (
	"fmt"
	"net/http"
	"time"
)

type GristClient struct {
	client  http.Client
	header  http.Header
	rootUrl string
}

func NewGristClient(key string) *GristClient {
	client := GristClient{}
	client.header.Add("Authorization", fmt.Sprintf("Bearer %v", key))
	client.rootUrl = "https://docs.getgrist.com/api"
	client.client = http.Client{
		Timeout: 5 * time.Second,
	}
	return &client
}
