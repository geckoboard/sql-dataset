package geckoboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	config Config
	client *http.Client
}

func New(config Config) *Client {
	cfg := defaultConfig()
	cfg.mergeIn(config)

	return &Client{
		client: http.DefaultClient,
		config: cfg,
	}
}

func (c Client) sendNewRequest(method, path string, body interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

func (c Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	url := fmt.Sprintf("%s%s", c.config.URL, path)
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.config.Key, "")

	return req, nil
}

func (c Client) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		return resp, nil
	default:
		var gbe Error
		json.NewDecoder(resp.Body).Decode(&gbe)
		return resp, gbe
	}
}
