package client

import "net/http"

type Client struct {
	BaseURL string
	HTTP    *http.Client

	// Account *AccountClient
}

func NextClient(baseURL string) *Client {
	httpClient := &http.Client{}

	c := &Client{
		BaseURL: baseURL,
		HTTP:    httpClient,
	}

	return c
}
