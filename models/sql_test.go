package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBuildDatasetSQLiteDriver(t *testing.T) {
	testCases := []struct {
		config Config
		out    []map[string]interface{}
		err    string
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
			err: `Scan failed: sql: Scan error on column index 0: strconv.ParseInt: parsing "everdeen": invalid syntax`,
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
			// StringType and Number as an int64
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
			out: []map[string]interface{}{
				{
					"app":         "",
					"build_count": int64(2),
				},
				{
					"app":         "everdeen",
					"build_count": int64(2),
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
			// Date only with money type (sqlite lib doesn't support DATE(col) in select as time.Time)
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT created_at, CAST(build_cost*100 AS INTEGER) FROM builds order by created_at",
						Fields: []Field{
							{Name: "Day", Type: DateType},
							{Name: "Build Cost", Type: MoneyType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(54),
				},
				{
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(144),
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(264),
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": 0,
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": 0,
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(1132),
				},
				{
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(111),
				},
				{
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(24),
				},
				{
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(92),
				},
			},
			err: "",
		},
		{
			// Datetime type example
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at FROM builds order by created_at;",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Day", Type: DatetimeType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:13:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "westworld",
					"day": parseTime("2017-03-23T15:11:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-03-23T16:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:22:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:44:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "react",
					"day": parseTime("2017-04-23T12:32:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:42:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:43:00Z", t).Format(time.RFC3339),
				},
			},
			err: "",
		},
		{
			// PercentageType with stringType
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, CAST(percent_passed/100.00 AS FLOAT) FROM builds order by app_name, created_at",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Percentage Completed", Type: PercentageType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "",
					"percentage_completed": 0.01,
				},
				{
					"app": "",
					"percentage_completed": 0.34,
				},
				{
					"app": "everdeen",
					"percentage_completed": 0.8,
				},
				{
					"app": "everdeen",
					"percentage_completed": 1.0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.24,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.55,
				},
				{
					"app": "react",
					"percentage_completed": 0.95,
				},
				{
					"app": "westworld",
					"percentage_completed": 0,
				},
			},
			err: "",
		},
		{
			// NumberType as float64 and date only type
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: SQLiteDriver,
					URL:    "fixtures/db.sqlite",
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at, run_time FROM builds where app_name <> '' order by app_name, created_at",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.31882276212,
				},
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 144.31838122382,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.21882232124,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 77.21381276421,
				},
				{
					"app":      "react",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 118.18382961212,
				},
				{
					"app":      "westworld",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 321.93774373,
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
						SQL: "SELECT app_name, null, ROUND(run_time, 7) FROM builds order by created_at limit 1",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     nil,
					"run_time": 0.3188228,
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

		if len(out) != len(tc.out) {
			t.Errorf("[%d] Expected slice size %d but got %d", idx, len(tc.out), len(out))
			continue
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

func TestBuildDatasetPostgresDriver(t *testing.T) {
	// Setup the postgres and run the insert
	env, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		t.Errorf("This test requires real postgres db using env:POSTGRES_URL ensure the db exists" +
			" eg. postgres://postgres:postgres@localhost:5432/testDbName?sslmode=disable")

		return
	}

	db, err := sql.Open("postgres", env)
	contents, err := ioutil.ReadFile("fixtures/postgres.sql")
	if err != nil {
		t.Fatal(err)
	}

	queries := strings.Split(string(contents), ";")

	for _, query := range queries {
		if _, err = db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		config Config
		out    []map[string]interface{}
		err    string
	}{
		{
			config: Config{
				GeckoboardAPIKey: "1234-12345",
				RefreshTimeSec:   120,
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    "postgres://postgresql:postgres@fakehost:5432",
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
			err: "Database query failed: dial tcp: lookup fakehost: no such host",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
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
			err: `Scan failed: sql: Scan error on column index 0: can't convert string "" to number`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
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
			err: `Database query failed: pq: column "create_at" does not exist`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
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
			err: `Database query failed: pq: column "builds.build_cost" must appear in the GROUP BY clause or be used in an aggregate function`,
		},
		{
			// StringType and Number as an int64
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
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
			out: []map[string]interface{}{
				{
					"app":         "",
					"build_count": int64(2),
				},
				{
					"app":         "everdeen",
					"build_count": int64(2),
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
			// Date only with money type grouping by date in sql
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT DATE(created_at) dte, SUM(CAST(build_cost*100 AS INTEGER)) FROM builds GROUP BY DATE(created_at) order by DATE(created_at)",
						Fields: []Field{
							{Name: "Day", Type: DateType},
							{Name: "Build Cost", Type: MoneyType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(198),
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(1396),
				},
				{
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(227),
				},
			},
			err: "",
		},
		{
			// Datetime type example
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at FROM builds order by created_at;",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Day", Type: DatetimeType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:13:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "westworld",
					"day": parseTime("2017-03-23T15:11:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-03-23T16:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:22:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:44:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "react",
					"day": parseTime("2017-04-23T12:32:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:42:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:43:00Z", t).Format(time.RFC3339),
				},
			},
			err: "",
		},
		{
			// PercentageType with stringType
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, CAST(percent_passed/100.00 AS FLOAT) FROM builds order by app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Percentage Completed", Type: PercentageType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "",
					"percentage_completed": 0.01,
				},
				{
					"app": "",
					"percentage_completed": 0.34,
				},
				{
					"app": "everdeen",
					"percentage_completed": 0.8,
				},
				{
					"app": "everdeen",
					"percentage_completed": 1.0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.55,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.24,
				},
				{
					"app": "react",
					"percentage_completed": 0.95,
				},
				{
					"app": "westworld",
					"percentage_completed": 0,
				},
			},
			err: "",
		},
		{
			// NumberType as float64 and date only type
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at, run_time FROM builds where app_name <> '' order by app_name, created_at",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.31882276212,
				},
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 144.31838122382,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.21882232124,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 77.21381276421,
				},
				{
					"app":      "react",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 118.18382961212,
				},
				{
					"app":      "westworld",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 321.93774373,
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, null, ROUND(CAST(run_time AS NUMERIC), 7) FROM builds order by created_at limit 1",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     nil,
					"run_time": 0.3188228,
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

		if len(out) != len(tc.out) {
			t.Errorf("[%d] Expected slice size %d but got %#v", idx, len(tc.out), out)
			continue
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

func TestBuildDatasetMySQLDriver(t *testing.T) {
	// Setup the postgres and run the insert
	env, ok := os.LookupEnv("MYSQL_URL")
	if !ok {
		t.Errorf("This test requires real mysql db using env:MYSQL_URL ensure the db exists" +
			" eg. [username[:password]@][protocol[(address)]]/dbname[?parseTime=true]")

		return
	}

	db, err := sql.Open("mysql", env)
	contents, err := ioutil.ReadFile("fixtures/mysql.sql")
	if err != nil {
		t.Fatal(err)
	}

	queries := strings.Split(string(contents), ";")

	for _, query := range queries {
		if _, err = db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		config Config
		out    []map[string]interface{}
		err    string
	}{
		{
			config: Config{
				GeckoboardAPIKey: "1234-12345",
				RefreshTimeSec:   120,
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    "root@tcp(fakehost:3306)/testdb",
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
			err: "Database query failed: dial tcp: lookup fakehost: no such host",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
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
			err: `Scan failed: sql: Scan error on column index 0: strconv.ParseInt: parsing "everdeen": invalid syntax`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
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
			err: `Database query failed: Error 1054: Unknown column 'create_at' in 'field list'`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
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
			err: `Database query failed: Error 1055: Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'testdb.builds.build_cost' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by`,
		},
		{
			// StringType and Number as an int64
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
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
			out: []map[string]interface{}{
				{
					"app":         "",
					"build_count": int64(2),
				},
				{
					"app":         "everdeen",
					"build_count": int64(2),
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
			// Date only with money type grouping by date in sql
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT DATE(created_at) dte, SUM(CAST(build_cost*100 AS SIGNED INTEGER)) FROM builds GROUP BY DATE(created_at) order by DATE(created_at)",
						Fields: []Field{
							{Name: "Day", Type: DateType},
							{Name: "Build Cost", Type: MoneyType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(198),
				},
				{
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(1396),
				},
				{
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"build_cost": int64(227),
				},
			},
			err: "",
		},
		{
			// Datetime type example
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at FROM builds order by created_at;",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Day", Type: DatetimeType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "everdeen",
					"day": parseTime("2017-03-21T11:13:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "westworld",
					"day": parseTime("2017-03-23T15:11:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-03-23T16:12:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:22:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "",
					"day": parseTime("2017-03-23T16:44:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "react",
					"day": parseTime("2017-04-23T12:32:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:42:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:43:00Z", t).Format(time.RFC3339),
				},
			},
			err: "",
		},
		{
			// PercentageType with stringType
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, CAST(percent_passed/100.00 AS DECIMAL(3,2)) FROM builds order by app_name, created_at",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Percentage Completed", Type: PercentageType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "",
					"percentage_completed": 0.01,
				},
				{
					"app": "",
					"percentage_completed": 0.34,
				},
				{
					"app": "everdeen",
					"percentage_completed": 0.8,
				},
				{
					"app": "everdeen",
					"percentage_completed": 1.0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.24,
				},
				{
					"app": "geckoboard-ruby",
					"percentage_completed": 0.55,
				},
				{
					"app": "react",
					"percentage_completed": 0.95,
				},
				{
					"app": "westworld",
					"percentage_completed": 0,
				},
			},
			err: "",
		},
		{
			// NumberType as float64 and date only type
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, created_at, run_time FROM builds where app_name <> '' order by app_name, created_at",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.31882276212,
				},
				{
					"app":      "everdeen",
					"date":     parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
					"run_time": 144.31838122382,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 0.21882232124,
				},
				{
					"app":      "geckoboard-ruby",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 77.21381276421,
				},
				{
					"app":      "react",
					"date":     parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 118.18382961212,
				},
				{
					"app":      "westworld",
					"date":     parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
					"run_time": 321.93774373,
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, null, ROUND(run_time, 7) FROM builds order by created_at limit 1",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Date", Type: DateType},
							{Name: "Run time", Type: NumberType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":      "everdeen",
					"date":     nil,
					"run_time": 0.3188228,
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

		if len(out) != len(tc.out) {
			fmt.Printf("%#v\n", out)
			t.Errorf("[%d] Expected slice size %d but got %d", idx, len(tc.out), len(out))
			continue
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
