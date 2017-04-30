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

var originalBatchRows = 500

type GBRequest struct {
	Path string
	Body string
}

func TestEndToEndFlow(t *testing.T) {
	testCases := []struct {
		config      models.Config
		batchRows   int
		expectError bool
		gbHits      int
		gbReqs      []GBRequest
	}{
		{
			config: models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.SQLiteDriver,
					URL:    filepath.Join("models", "fixtures", "db.sqlite"),
				},
				Datasets: []models.Dataset{
					{
						Name:       "app.counts",
						SQL:        "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						UpdateType: models.Append,
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Build Count", Type: models.NumberType},
						},
					},
					{
						Name:       "app.build.costs",
						UpdateType: models.Append,
						SQL:        "SELECT app_name, CAST(build_cost*100 AS INTEGER) FROM builds GROUP BY app_name order by app_name",
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
					Body: `{"data":[{"app":"","build_count":2},{"app":"everdeen","build_count":2},{"app":"geckoboard-ruby","build_count":3},{"app":"react","build_count":1},{"app":"westworld","build_count":1}]}`,
				},
				{
					Path: "/datasets/app.build.costs",
					Body: `{"id":"app.build.costs","fields":{"app":{"name":"App","type":"string","currency_code":""},"build_cost":{"name":"Build Cost","type":"money","currency_code":"USD"}},"unique_by":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/app.build.costs/data",
					Body: `{"data":[{"app":"","build_cost":1132},{"app":"everdeen","build_cost":144},{"app":"geckoboard-ruby","build_cost":0},{"app":"react","build_cost":111},{"app":"westworld","build_cost":264}]}`,
				},
			},
		},
		{
			// Replace update type sends the first batch only
			config: models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.SQLiteDriver,
					URL:    filepath.Join("models", "fixtures", "db.sqlite"),
				},
				Datasets: []models.Dataset{
					{
						Name:       "apps.run.time",
						SQL:        "SELECT app_name, run_time FROM builds ORDER BY app_name",
						UpdateType: models.Replace,
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Run time", Type: models.NumberType, FloatPrecision: 64},
						},
					},
				},
			},
			batchRows:   4,
			expectError: true,
			gbReqs: []GBRequest{
				{
					Path: "/datasets/apps.run.time",
					Body: `{"id":"apps.run.time","fields":{"app":{"name":"App","type":"string","currency_code":""},"run_time":{"name":"Run time","type":"number","currency_code":""}},"unique_by":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/apps.run.time/data",
					Body: `{"data":[{"app":"","run_time":0.12349876543},{"app":"","run_time":46.432763287},{"app":"everdeen","run_time":0.31882276212},{"app":"everdeen","run_time":144.31838122382}]}`,
				},
			},
		},
		{
			// Append update type sends multiple requests in batches of batchRow limit when more rows exist
			config: models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.SQLiteDriver,
					URL:    filepath.Join("models", "fixtures", "db.sqlite"),
				},
				Datasets: []models.Dataset{
					{
						Name:       "apps.run.time",
						SQL:        "SELECT app_name, run_time FROM builds ORDER BY app_name",
						UpdateType: models.Append,
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Run time", Type: models.NumberType, FloatPrecision: 64},
						},
					},
				},
			},
			batchRows: 4,
			gbReqs: []GBRequest{
				{
					Path: "/datasets/apps.run.time",
					Body: `{"id":"apps.run.time","fields":{"app":{"name":"App","type":"string","currency_code":""},"run_time":{"name":"Run time","type":"number","currency_code":""}},"unique_by":null,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/apps.run.time/data",
					Body: `{"data":[{"app":"","run_time":0.12349876543},{"app":"","run_time":46.432763287},{"app":"everdeen","run_time":0.31882276212},{"app":"everdeen","run_time":144.31838122382}]}`,
				},
				{
					Path: "/datasets/apps.run.time/data",
					Body: `{"data":[{"app":"geckoboard-ruby","run_time":0.21882232124},{"app":"geckoboard-ruby","run_time":77.21381276421},{"app":"geckoboard-ruby","run_time":0},{"app":"react","run_time":118.18382961212}]}`,
				},
				{
					Path: "/datasets/apps.run.time/data",
					Body: `{"data":[{"app":"westworld","run_time":321.93774373}]}`,
				},
			},
		},
		{
			// Unique by correctly used and sent - doesn't do validation used with the correct update type
			config: models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.SQLiteDriver,
					URL:    filepath.Join("models", "fixtures", "db.sqlite"),
				},
				Datasets: []models.Dataset{
					{
						Name:       "app.counts",
						SQL:        "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						UpdateType: models.Append,
						UniqueBy:   []string{"app_name"},
						Fields: []models.Field{
							{Name: "App", Type: models.StringType},
							{Name: "Build Count", Type: models.NumberType},
						},
					},
				},
			},
			gbReqs: []GBRequest{
				{
					Path: "/datasets/app.counts",
					Body: `{"id":"app.counts","fields":{"app":{"name":"App","type":"string","currency_code":""},"build_count":{"name":"Build Count","type":"number","currency_code":""}},"unique_by":["app_name"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
				},
				{
					Path: "/datasets/app.counts/data",
					Body: `{"data":[{"app":"","build_count":2},{"app":"everdeen","build_count":2},{"app":"geckoboard-ruby","build_count":3},{"app":"react","build_count":1},{"app":"westworld","build_count":1}]}`,
				},
			},
		},
	}

	for i, tc := range testCases {
		batchRows = originalBatchRows

		if tc.batchRows != 0 {
			batchRows = tc.batchRows
		}

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
				t.Errorf("[%d] Expected geckoboard request body %s but got %s", i, tcReq.Body, string(b))
			}

			fmt.Fprintf(w, `{}`)
		}))

		gbClient = gb.New(gb.Config{Key: "fakeKey", URL: gbWS.URL})

		bol := processAllDatasets(&tc.config)

		if tc.expectError != bol {
			t.Errorf("[%d] Expected hasErrors to be %t but got %t", i, tc.expectError, bol)
		}

		if tc.gbHits != len(tc.gbReqs) {
			t.Errorf("Expected %d requests but got %d", len(tc.gbReqs), tc.gbHits)
		}
		gbWS.Close()
	}
}
