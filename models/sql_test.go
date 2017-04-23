package models

import (
	"testing"
	"time"
)

func TestBuildDataset(t *testing.T) {
	testCases := []struct {
		config       Config
		fieldKeyType map[string]FieldType
		out          []map[string]interface{}
		err          string
	}{
		{
			config: Config{
				GeckoboardAPIKey: "1234-12345",
				RefreshTimeSec:   120,
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "models/fixtures/nonexisting",
				},
				Datasets: []Dataset{
					{
						Name:       "users.count",
						UpdateType: Replace,
						SQL:        "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: MoneyType},
						},
					},
				},
			},
			out: nil,
			err: "Database query failed: unable to open database file",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "App", Type: NumberType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
			},
			out: nil,
			err: `Scan failed: sql: Scan error on column index 0: converting driver.Value type []uint8 ("") to a int64: invalid syntax`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, create_at FROM builds order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Date", Type: DatetimeType},
						},
					},
				},
			},
			out: nil,
			err: `Database query failed: no such column: create_at`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, build_cost, created_at FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
			},
			out: nil,
			err: `Scan failed: sql: expected 3 destination arguments in Scan, not 2`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, count(*) FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
			},
			fieldKeyType: map[string]FieldType{
				"app":         StringType,
				"build_count": NumberType,
			},
			out: []map[string]interface{}{
				{
					"app":         "",
					"build_count": int64(1),
				},
				{
					"app":         "everdeen",
					"build_count": int64(1),
				},
				{
					"app":         "geckoboard-ruby",
					"build_count": int64(3),
				},
				{
					"app":         "react",
					"build_count": int64(1),
				},
				{
					"app":         "westworld",
					"build_count": int64(1),
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT build_cost, created_at FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "Build Cost", Type: MoneyType},
							{Name: "Day", Type: DateType},
						},
					},
				},
			},
			fieldKeyType: map[string]FieldType{
				"build_cost": MoneyType,
				"day":        DatetimeType,
			},
			out: []map[string]interface{}{
				{
					"build_cost": 11.32,
					"day":        parseTime("2017-03-23T16:44:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": 0.54,
					"day":        parseTime("2017-03-21T11:12:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": float64(0),
					"day":        parseTime("2017-03-23T16:22:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": 1.11,
					"day":        parseTime("2017-04-23T12:32:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": 2.64,
					"day":        parseTime("2017-03-23T15:11:00Z", t).Format(dateFormat),
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, SUM(build_cost), updated_at FROM builds GROUP BY app_name order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Cost", Type: MoneyType},
							{Name: "Day", Type: DateType},
						},
					},
				},
			},
			fieldKeyType: map[string]FieldType{
				"app":        StringType,
				"build_cost": MoneyType,
				"day":        DatetimeType,
			},
			out: []map[string]interface{}{
				{
					"app":        "",
					"build_cost": 11.32,
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "everdeen",
					"build_cost": 0.54,
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": float64(0.48),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "react",
					"build_cost": 1.11,
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "westworld",
					"build_cost": 2.64,
					"day":        nil,
				},
			},
			err: "",
		},
	}

	for idx, tc := range testCases {
		out, err := tc.config.Datasets[0].BuildDataset(tc.config.DatabaseConfig)

		if tc.err == "" && err != nil {
			t.Errorf("[%d] Expected no error but got %s", idx, err)
		}

		if err != nil && tc.err != err.Error() {
			t.Errorf("[%d] Expected error %s but got %s", idx, tc.err, err)
		}

		for i, mp := range out {
			for k, v := range mp {
				if tc.out[i][k] != v {
					t.Errorf("[%d-%d] Expected key '%s' to have value %v but got %v", idx, i, k, tc.out[i][k], v)
				}
			}
		}
	}
}

func parseTime(str string, t *testing.T) time.Time {
	tme, err := time.Parse(time.RFC3339, str)

	if err != nil {
		t.Fatal(err)
	}

	return tme

}
