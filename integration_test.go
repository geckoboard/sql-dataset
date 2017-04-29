package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	gb "github.com/geckoboard/geckoboard-go"
	"github.com/geckoboard/sql-dataset/models"
)

type GBRequest struct {
	Path string
	Body string
}

func TestEndToEndFlow(t *testing.T) {
	testCases := []struct {
		config models.Config
		gbHits int
		gbReqs []GBRequest
	}{
		{
			config: models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.SQLiteDriver,
					URL:    filepath.Join("models", "fixtures", "db.sqlite"),
				},
				Datasets: []models.Dataset{
					{
						Name: "app.counts",
						SQL:  "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Build Count", Type: models.NumberType},
						},
					},
					{
						Name: "app.build.costs",
						SQL:  "SELECT app_name, CAST(build_cost*100 AS INTEGER) FROM builds GROUP BY app_name order by app_name",
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Build Cost", Type: models.MoneyType, CurrencyCode: "USD"},
						},
					},
				},
			},
			gbReqs: []GBRequest{
				{
					Path: "/datasets/app.counts",
					Body: `{"id":"app.counts","fields":{"app":{"name":"App","type":"string","currency_code":""},"build_count":{"name":"Build Count","type":"number","currency_code":""}},"unique_by":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/app.counts/data",
					Body: `{"data":[{"app":"","build_count":1},{"app":"everdeen","build_count":1},{"app":"geckoboard-ruby","build_count":3},{"app":"react","build_count":1},{"app":"westworld","build_count":1}]}`,
				},
				{
					Path: "/datasets/app.build.costs",
					Body: `{"id":"app.build.costs","fields":{"app":{"name":"App","type":"string","currency_code":""},"build_cost":{"name":"Build Cost","type":"money","currency_code":"USD"}},"unique_by":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/app.build.costs/data",
					Body: `{"data":[{"app":"","build_cost":1132},{"app":"everdeen","build_cost":54},{"app":"geckoboard-ruby","build_cost":0},{"app":"react","build_cost":111},{"app":"westworld","build_cost":264}]}`,
				},
			},
		},
	}

	for i, tc := range testCases {

		gbWS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc.gbHits++

			if tc.gbHits-1 >= len(tc.gbReqs) {
				t.Errorf("[%d] Got unexpected extra request for geckoboard unable to process", i)
				return
			}

			tcReq := tc.gbReqs[tc.gbHits-1]

			if tcReq.Path != r.URL.Path {
				t.Errorf("[%d] Expected geckoboard request path %s but got %s", i, tcReq.Path, r.URL.Path)
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal("Failed to consume body", err)
			}

			if strings.TrimSpace(string(b)) != tcReq.Body {
				t.Errorf("Expected geckoboard request body %s but got %s", tcReq.Body, string(b))
			}

			fmt.Fprintf(w, `{}`)
		}))

		gbClient = gb.New(gb.Config{Key: "fakeKey", URL: gbWS.URL})

		if bol := processAllDatasets(&tc.config); bol {
			t.Errorf("[%d] Expected no errors but processAllDatasets suggested it errored", i)
		}

		gbWS.Close()
	}
}
