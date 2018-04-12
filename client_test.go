package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geckoboard/sql-dataset/models"
)

type request struct {
	Method  string
	Headers map[string]string
	Path    string
	Body    string
}

type response struct {
	code int
	body string
}

const (
	apiKey  = "ap1K3Y"
	expAuth = "Basic YXAxSzNZOg=="
)

func TestFindOrCreateDataset(t *testing.T) {
	userAgent = "SQL-Dataset/test-fake"

	testCases := []struct {
		dataset  models.Dataset
		request  request
		response *response
		err      string
	}{
		{
			dataset: models.Dataset{
				Name: "active.users.count",
				Fields: []models.Field{
					{
						Name:         "mrr",
						Type:         models.MoneyType,
						CurrencyCode: "USD",
					},
					{
						Name: "signups",
						Type: models.NumberType,
					},
				},
			},
			request: request{
				Method: http.MethodPut,
				Path:   "/datasets/active.users.count",
				Body:   `{"id":"active.users.count","fields":{"mrr":{"type":"money","name":"mrr","currency_code":"USD"},"signups":{"type":"number","name":"signups"}}}`,
			},
		},
		{
			dataset: models.Dataset{
				Name: "active.users.count",
				Fields: []models.Field{
					{
						Name: "day",
						Type: models.DatetimeType,
					},
					{
						Name:     "mrr Percent",
						Type:     models.PercentageType,
						Optional: true,
					},
				},
			},
			request: request{
				Method: http.MethodPut,
				Path:   "/datasets/active.users.count",
				Body:   `{"id":"active.users.count","fields":{"day":{"type":"datetime","name":"day"},"mrr_percent":{"type":"percentage","name":"mrr Percent","optional":true}}}`,
			},
		},
		{
			dataset: models.Dataset{
				Name:     "builds.count.by.day",
				UniqueBy: []string{"day"},
				Fields: []models.Field{
					{
						Name:     "Builds",
						Type:     models.NumberType,
						Optional: true,
					},
					{
						Name: "day",
						Type: models.DateType,
					},
				},
			},
			request: request{
				Method: http.MethodPut,
				Path:   "/datasets/builds.count.by.day",
				Body:   `{"id":"builds.count.by.day","unique_by":["day"],"fields":{"builds":{"type":"number","name":"Builds","optional":true},"day":{"type":"date","name":"day"}}}`,
			},
		},
		{
			// Verify 50x just returns generic error
			dataset: models.Dataset{
				Name: "active.users.count",
				Fields: []models.Field{
					{
						Name:         "mrr",
						Type:         models.MoneyType,
						CurrencyCode: "USD",
					},
					{
						Name: "signups",
						Type: models.NumberType,
					},
				},
			},
			response: &response{
				code: 500,
				body: "<html>Internal error</html>",
			},
			err: errUnexpectedResponse.Error(),
		},
		{
			// Verify 40x response unmarshalled and return message
			dataset: models.Dataset{
				Name: "active.users.count",
				Fields: []models.Field{
					{
						Name:         "mrr",
						Type:         models.MoneyType,
						CurrencyCode: "USD",
					},
					{
						Name: "signups",
						Type: models.NumberType,
					},
				},
			},
			response: &response{
				code: 400,
				body: `{"error":{"type":"ErrResourceInvalid","message":"Field name too short"}}`,
			},
			err: fmt.Sprintf(errInvalidPayload, "Field name too short"),
		},
		{
			// Verify 201 response correctly handled
			dataset: models.Dataset{
				Name: "active.users.count",
				Fields: []models.Field{
					{
						Name:         "mrr",
						Type:         models.MoneyType,
						CurrencyCode: "USD",
					},
					{
						Name: "signups",
						Type: models.NumberType,
					},
				},
			},
			response: &response{
				code: 201,
				body: `{}`,
			},
		},
	}

	for _, tc := range testCases {
		reqCount := 0

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount++

			if tc.response != nil {
				w.WriteHeader(tc.response.code)
				fmt.Fprintf(w, tc.response.body)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth != expAuth {
				t.Errorf("Expected authorization header '%s' but got '%s'", expAuth, auth)
			}

			ua := r.Header.Get("User-Agent")
			if ua != userAgent {
				t.Errorf("Expected user header '%s' but got '%s'", userAgent, ua)
			}

			if r.URL.Path != tc.request.Path {
				t.Errorf("Expected path '%s' but got '%s'", tc.request.Path, r.URL.Path)
			}

			if r.Method != tc.request.Method {
				t.Errorf("Expected method '%s' but got '%s'", tc.request.Method, r.Method)
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Errored while consuming body %s", err)
			}

			body := strings.Trim(string(b), "\n")

			if body != tc.request.Body {
				t.Errorf("Expected body '%s' but got '%s'", tc.request.Body, body)
			}
		}))

		gbHost = server.URL

		c := NewClient(apiKey)
		err := c.FindOrCreateDataset(&tc.dataset)

		if err != nil && tc.err == "" {
			t.Errorf("Expected no error but got '%s'", err)
		}

		if err == nil && tc.err != "" {
			t.Errorf("Expected error '%s' but got none", tc.err)
		}

		if err != nil && err.Error() != tc.err {
			t.Errorf("Expected error '%s' but got '%s'", tc.err, err)
		}

		if reqCount != 1 {
			t.Errorf("Expected one request but got %d", reqCount)
		}

		server.Close()
	}
}

// Preset the dataset rows and test we batch correctly return errors on status codes
func TestSendAllData(t *testing.T) {
	testCases := []struct {
		dataset  models.Dataset
		data     models.DatasetRows
		requests []request
		maxRows  int
		response *response
		err      string
	}{
		{
			// Error with 40x
			dataset: models.Dataset{
				Name:       "app.build.costs",
				UpdateType: models.Replace,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Run time", Type: models.NumberType},
				},
			},
			data: models.DatasetRows{
				{
					"app":      "acceptance",
					"run_time": 4421,
				},
			},
			requests: []request{{}},
			response: &response{
				code: 400,
				body: `{"error": {"type":"ErrMissingData", "message": "Missing data for 'app'"}}`,
			},
			err: fmt.Sprintf(errInvalidPayload, "Missing data for 'app'"),
		},
		{
			// Error with 50x
			dataset: models.Dataset{
				Name:       "app.build.costs",
				UpdateType: models.Replace,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Run time", Type: models.NumberType},
				},
			},
			data: models.DatasetRows{
				{
					"app":      "acceptance",
					"run_time": 4421,
				},
			},
			requests: []request{{}},
			response: &response{
				code: 500,
				body: "<html>Internal Server error</html>",
			},
			err: errUnexpectedResponse.Error(),
		},
		{
			//Replace dataset with no data
			dataset: models.Dataset{
				Name:       "app.no.data",
				UpdateType: models.Replace,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Percent", Type: models.PercentageType},
				},
			},
			data: models.DatasetRows{},
			requests: []request{
				{
					Method: http.MethodPut,
					Path:   "/datasets/app.no.data/data",
					Body:   `{"data":[]}`,
				},
			},
		},
		{
			//Replace dataset under the batch rows limit doesn't error
			dataset: models.Dataset{
				Name:       "app.reliable.percent",
				UpdateType: models.Replace,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Percent", Type: models.PercentageType},
				},
			},
			data: models.DatasetRows{
				{
					"app":     "acceptance",
					"percent": 0.43,
				},
				{
					"app":     "redis",
					"percent": 0.22,
				},
				{
					"app":     "api",
					"percent": 0.66,
				},
			},
			requests: []request{
				{
					Method: http.MethodPut,
					Path:   "/datasets/app.reliable.percent/data",
					Body:   `{"data":[{"app":"acceptance","percent":0.43},{"app":"redis","percent":0.22},{"app":"api","percent":0.66}]}`,
				},
			},
		},
		{
			//Append with no data makes no requests
			dataset: models.Dataset{
				Name:       "append.no.data",
				UpdateType: models.Append,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Count", Type: models.NumberType},
				},
			},
			data:     models.DatasetRows{},
			requests: []request{},
		},
		{
			//Append dataset under the batch rows limit
			dataset: models.Dataset{
				Name:       "app.builds.count",
				UpdateType: models.Append,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Count", Type: models.NumberType},
				},
			},
			data: models.DatasetRows{
				{
					"app":   "acceptance",
					"count": 88,
				},
				{
					"app":   "redis",
					"count": 55,
				},
				{
					"app":   "api",
					"count": 214,
				},
			},
			requests: []request{
				{
					Method: http.MethodPost,
					Path:   "/datasets/app.builds.count/data",
					Body:   `{"data":[{"app":"acceptance","count":88},{"app":"redis","count":55},{"app":"api","count":214}]}`,
				},
			},
		},
		{
			//Replace dataset over the batch rows limit sends first 3 and errors
			dataset: models.Dataset{
				Name:       "app.build.costs",
				UpdateType: models.Replace,
				Fields: []models.Field{
					{Name: "App", Type: models.StringType},
					{Name: "Cost", Type: models.MoneyType},
				},
			},
			data: models.DatasetRows{
				{
					"app":  "acceptance",
					"cost": 4421,
				},
				{
					"app":  "redis",
					"cost": 221,
				},
				{
					"app":  "api",
					"cost": 212,
				},
				{
					"app":  "integration",
					"cost": 121,
				},
			},
			requests: []request{
				{
					Method: http.MethodPut,
					Path:   "/datasets/app.build.costs/data",
					Body:   `{"data":[{"app":"acceptance","cost":4421},{"app":"redis","cost":221},{"app":"api","cost":212}]}`,
				},
			},
			maxRows: 3,
			err:     fmt.Sprintf(errMoreRowsToSend, 4, 3),
		},
		{
			//Append dataset sends all data in batches
			dataset: models.Dataset{
				Name:       "animal.run.time",
				UpdateType: models.Append,
				Fields: []models.Field{
					{Name: "animal", Type: models.StringType},
					{Name: "Run time", Type: models.NumberType},
				},
			},
			data: models.DatasetRows{
				{
					"animal":   "worm",
					"run_time": 621,
				},
				{
					"animal":   "snail",
					"run_time": 521,
				},
				{
					"animal":   "duck",
					"run_time": 41,
				},
				{
					"animal":   "geese",
					"run_time": 44,
				},
				{
					"animal":   "terrapin",
					"run_time": 444,
				},
				{
					"animal":   "bird",
					"run_time": 22,
				},
			},
			requests: []request{
				{
					Method: http.MethodPost,
					Path:   "/datasets/animal.run.time/data",
					Body:   `{"data":[{"animal":"worm","run_time":621},{"animal":"snail","run_time":521},{"animal":"duck","run_time":41}]}`,
				},
				{
					Method: http.MethodPost,
					Path:   "/datasets/animal.run.time/data",
					Body:   `{"data":[{"animal":"geese","run_time":44},{"animal":"terrapin","run_time":444},{"animal":"bird","run_time":22}]}`,
				},
			},
			maxRows: 3,
		},
		{
			//Append dataset sends all data in batches remaining one
			dataset: models.Dataset{
				Name:       "animal.run.time",
				UpdateType: models.Append,
				Fields: []models.Field{
					{Name: "Animal", Type: models.StringType},
					{Name: "Run time", Type: models.NumberType},
				},
			},
			data: models.DatasetRows{
				{
					"animal":   "worm",
					"run_time": 621,
				},
				{
					"animal":   "snail",
					"run_time": 521,
				},
				{
					"animal":   "duck",
					"run_time": 41,
				},
				{
					"animal":   "geese",
					"run_time": 44,
				},
				{
					"animal":   "terrapin",
					"run_time": 444,
				},
				{
					"animal":   "bird",
					"run_time": 22,
				},
				{
					"animal":   "squirrel",
					"run_time": 88,
				},
			},
			requests: []request{
				{
					Method: http.MethodPost,
					Path:   "/datasets/animal.run.time/data",
					Body:   `{"data":[{"animal":"worm","run_time":621},{"animal":"snail","run_time":521},{"animal":"duck","run_time":41}]}`,
				},
				{
					Method: http.MethodPost,
					Path:   "/datasets/animal.run.time/data",
					Body:   `{"data":[{"animal":"geese","run_time":44},{"animal":"terrapin","run_time":444},{"animal":"bird","run_time":22}]}`,
				},
				{
					Method: http.MethodPost,
					Path:   "/datasets/animal.run.time/data",
					Body:   `{"data":[{"animal":"squirrel","run_time":88}]}`,
				},
			},
			maxRows: 3,
		},
	}

	for _, tc := range testCases {
		reqCount := 0

		if tc.maxRows != 0 {
			maxRows = tc.maxRows
		} else {
			maxRows = 500
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount++

			if tc.response != nil {
				w.WriteHeader(tc.response.code)
				fmt.Fprintf(w, tc.response.body)
				return
			}

			if reqCount > len(tc.requests) {
				t.Errorf("Got unexpected extra requests")
				return
			}

			tcReq := tc.requests[reqCount-1]

			auth := r.Header.Get("Authorization")
			if auth != expAuth {
				t.Errorf("Expected authorization header '%s' but got '%s'", expAuth, auth)
			}

			if r.URL.Path != tcReq.Path {
				t.Errorf("Expected path '%s' but got '%s'", tcReq.Path, r.URL.Path)
			}

			if r.Method != tcReq.Method {
				t.Errorf("Expected method '%s' but got '%s'", tcReq.Method, r.Method)
			}

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Errored while consuming body %s", err)
			}

			body := strings.Trim(string(b), "\n")

			if body != tcReq.Body {
				t.Errorf("Expected body '%s' but got '%s'", tcReq.Body, body)
			}
		}))

		gbHost = server.URL

		c := NewClient(apiKey)
		err := c.SendAllData(&tc.dataset, tc.data)

		if err != nil && tc.err == "" {
			t.Errorf("Expected no error but got '%s'", err)
		}

		if err == nil && tc.err != "" {
			t.Errorf("Expected error '%s' but got none", tc.err)
		}

		if err != nil && err.Error() != tc.err {
			t.Errorf("Expected error '%s' but got '%s'", tc.err, err)
		}

		if reqCount != len(tc.requests) {
			t.Errorf("Expected %d requests but got %d", len(tc.requests), reqCount)
		}

		server.Close()
	}
}
