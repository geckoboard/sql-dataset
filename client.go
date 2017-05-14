package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/geckoboard/sql-dataset/models"
)

type Client struct {
	apiKey string
	client *http.Client
}

type Error struct {
	Detail `json:"error"`
}

type Detail struct {
	Message string `json:"message"`
}

type DataPayload struct {
	Data models.DatasetRows `json:"data"`
}

var (
	gbHost    = "https://api.geckoboard.com"
	batchRows = 500

	errUnexpectedResponse = errors.New("Unexpected server error response from Geckoboard")
	errMoreRowsToSend     = "Sent the first %d rows, %d rows existed " +
		"to support sending more change dataset update type " +
		"from 'replace' to 'append' to support upto 5000 rows"
)

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{Timeout: time.Second * 10},
	}
}

func (c Client) FindOrCreateDataset(ds *models.Dataset) error {
	ds.BuildGBFields()
	resp, err := c.makeRequest(http.MethodPut, fmt.Sprintf("/datasets/%s", ds.Name), ds)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return handleResponse(resp)
}

func (c Client) sendData(ds *models.Dataset, data models.DatasetRows) (err error) {
	method := http.MethodPost

	if ds.UpdateType == models.Replace {
		method = http.MethodPut
	}

	resp, err := c.makeRequest(method, fmt.Sprintf("/datasets/%s/data", ds.Name), DataPayload{data})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return handleResponse(resp)
}

// SendAllData determines how to send the data to Geckoboard and returns an error
// if there is too much data for replace dataset and batches requests for append
func (c Client) SendAllData(ds *models.Dataset, data models.DatasetRows) (err error) {
	switch ds.UpdateType {
	case models.Replace:
		if len(data) > batchRows {
			err = c.sendData(ds, data[0:batchRows])
			if err == nil {
				err = fmt.Errorf(errMoreRowsToSend, batchRows, len(data))
			}
		} else {
			err = c.sendData(ds, data)
		}
	case models.Append:
		grps := len(data) / batchRows

		for i := 0; i <= grps; i++ {
			if i == grps {
				if (batchRows*i)+1 <= len(data) {
					err = c.sendData(ds, data[batchRows*i:])
				}
			} else {
				err = c.sendData(ds, data[batchRows*i:batchRows*(i+1)])
			}

			if err != nil {
				return err
			}
		}
	}

	return err
}

func (c Client) makeRequest(method, path string, body interface{}) (resp *http.Response, err error) {
	var b []byte

	if body != nil {
		if b, err = json.Marshal(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", gbHost, path), bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, "")

	return c.client.Do(req)
}

func handleResponse(resp *http.Response) error {
	switch res := resp.StatusCode; true {
	case res >= 200 && res < 300:
		return nil
	case res >= 400 && res < 500:
		var err Error
		json.NewDecoder(resp.Body).Decode(&err)
		return fmt.Errorf("response error: %s", err.Detail.Message)
	default:
		return errUnexpectedResponse
	}
}
