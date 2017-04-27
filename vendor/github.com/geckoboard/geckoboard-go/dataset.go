package geckoboard

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	NumberFieldType     = "number"
	DateFieldType       = "date"
	DatetimeFieldType   = "datetime"
	StringFieldType     = "string"
	MoneyFieldType      = "money"
	PercentageFieldType = "percentage"
)

type DataSet struct {
	ID        string    `json:"id,omitempty"`
	Fields    Fields    `json:"fields"`
	UniqueBy  []string  `json:"unique_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Fields map[string]Field

type Field struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	CurrencyCode string `json:"currency_code"`
}

type Record map[string]interface{}

func (s DataSet) Delete(c *Client) error {
	resp, err := c.sendNewRequest("DELETE", fmt.Sprintf("/datasets/%s", s.ID), nil)
	if err != nil {
		return err
	}

	var body struct{}
	return json.NewDecoder(resp.Body).Decode(&body)
}

func (s DataSet) SendAll(c *Client, recs interface{}) error {
	data := struct {
		Data interface{} `json:"data"`
	}{Data: recs}

	resp, err := c.sendNewRequest("PUT", fmt.Sprintf("/datasets/%s/data", s.ID), data)
	if err != nil {
		return err
	}

	var body struct{}
	return json.NewDecoder(resp.Body).Decode(&body)
}

func (s DataSet) Append(c *Client, recs interface{}) error {
	data := struct {
		Data interface{} `json:"data"`
	}{Data: recs}

	resp, err := c.sendNewRequest("POST", fmt.Sprintf("/datasets/%s/data", s.ID), data)
	if err != nil {
		return err
	}

	var body struct{}
	return json.NewDecoder(resp.Body).Decode(&body)
}

func (s *DataSet) FindOrCreate(c *Client) error {
	resp, err := c.sendNewRequest("PUT", fmt.Sprintf("/datasets/%s", s.ID), s)
	if err != nil {
		return err
	}

	var s2 DataSet
	err = json.NewDecoder(resp.Body).Decode(&s2)
	if err != nil {
		return err
	}

	mergeDataSets(s, &s2)

	return nil
}

func mergeDataSets(dOut, dIn *DataSet) {
	dOut.Fields = dIn.Fields
	dOut.CreatedAt = dIn.CreatedAt
	dOut.UpdatedAt = dIn.UpdatedAt
}
