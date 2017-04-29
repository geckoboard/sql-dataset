package models

import (
	"database/sql"
	"io/ioutil"
	"os"
	"strings"
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
						SQL: "SELECT CAST(build_cost*100 AS INTEGER), created_at FROM builds GROUP BY app_name order by app_name",
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
					"build_cost": int64(1132),
					"day":        parseTime("2017-03-23T16:44:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": int64(54),
					"day":        parseTime("2017-03-21T11:12:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": int64(0),
					"day":        parseTime("2017-03-23T16:22:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": int64(111),
					"day":        parseTime("2017-04-23T12:32:00Z", t).Format(dateFormat),
				},
				{
					"build_cost": int64(264),
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
						SQL: "SELECT app_name, CAST(SUM(build_cost)*100 AS INTEGER), updated_at FROM builds GROUP BY app_name order by app_name",
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
					"build_cost": int64(1132),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "everdeen",
					"build_cost": int64(54),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(48),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "react",
					"build_cost": int64(111),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "westworld",
					"build_cost": int64(264),
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

func TestBuildDatasetPostgresDriver(t *testing.T) {
	// Setup the postgres and run the insert
	env, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		t.Errorf("This test requires real postgres db using env:POSTGRES_URL ensure the db exists" +
			" eg. postgres://postgres:postgres@localhost:5432/testDbName?sslmode=disable")

		return
	}

	db, err := sql.Open("postgres", env)
	contents, err := ioutil.ReadFile("fixtures/insert_postgres.sql")
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
			err: `Scan failed: sql: Scan error on column index 0: converting driver.Value type string ("") to a int64: invalid syntax`,
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
						SQL: "SELECT app_name, build_cost, created_at FROM builds",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
			},
			out: nil,
			err: "Scan failed: sql: expected 3 destination arguments in Scan, not 2",
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
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
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
					Driver: PostgresDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, SUM(CAST(build_cost*100 AS INTEGER)), created_at FROM builds GROUP BY app_name, created_at ORDER BY app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Cost", Type: MoneyType},
							{Name: "Day", Type: DateType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":        "",
					"build_cost": int64(1132),
					"day":        parseTime("2017-03-23T16:44:00Z", t).Format(dateFormat),
				},
				{
					"app":        "everdeen",
					"build_cost": int64(54),
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(24),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(74),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(0),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "react",
					"build_cost": int64(111),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "westworld",
					"build_cost": int64(264),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
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
						SQL: "SELECT app_name, updated_at FROM builds ORDER BY app_name LIMIT 4",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Day", Type: DatetimeType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "",
					"day": parseTime("2017-03-23T16:45:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "everdeen",
					"day": parseTime("2017-04-23T11:14:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:42:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": nil,
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

func TestBuildDatasetMysqlDriver(t *testing.T) {
	// Setup the postgres and run the insert
	env, ok := os.LookupEnv("MYSQL_URL")
	if !ok {
		t.Errorf("This test requires real mysql db using env:MYSQL_URL ensure the db exists" +
			" eg. [username[:password]@][protocol[(address)]]/dbname[?parseTime=true]")

		return
	}

	db, err := sql.Open("mysql", env)
	contents, err := ioutil.ReadFile("fixtures/insert_mysql.sql")
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
				DatabaseConfig: &DatabaseConfig{
					Driver: MysqlDriver,
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
			err: `Scan failed: sql: Scan error on column index 0: converting driver.Value type []uint8 ("") to a int64: invalid syntax`,
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MysqlDriver,
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
					Driver: MysqlDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, build_cost, created_at FROM builds",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Count", Type: NumberType},
						},
					},
				},
			},
			out: nil,
			err: "Scan failed: sql: expected 3 destination arguments in Scan, not 2",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MysqlDriver,
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
					Driver: MysqlDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, SUM(CAST(build_cost*100 AS SIGNED INTEGER)), created_at FROM builds GROUP BY app_name, created_at ORDER BY app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Cost", Type: MoneyType},
							{Name: "Day", Type: DateType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":        "",
					"build_cost": int64(1132),
					"day":        parseTime("2017-03-23T16:44:00Z", t).Format(dateFormat),
				},
				{
					"app":        "everdeen",
					"build_cost": int64(54),
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(74),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(0),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": int64(24),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "react",
					"build_cost": int64(111),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "westworld",
					"build_cost": int64(264),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MysqlDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, SUM(build_cost/10), created_at FROM builds GROUP BY app_name, created_at ORDER BY app_name",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Build Cost", Type: PercentageType, FloatPrecision: 32},
							{Name: "Day", Type: DateType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app":        "",
					"build_cost": float32(1.132),
					"day":        parseTime("2017-03-23T16:44:00Z", t).Format(dateFormat),
				},
				{
					"app":        "everdeen",
					"build_cost": float32(0.054),
					"day":        parseTime("2017-03-21T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": float32(0.074),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": float32(0),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "geckoboard-ruby",
					"build_cost": float32(0.024),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "react",
					"build_cost": float32(0.111),
					"day":        parseTime("2017-04-23T00:00:00Z", t).Format(dateFormat),
				},
				{
					"app":        "westworld",
					"build_cost": float32(0.264),
					"day":        parseTime("2017-03-23T00:00:00Z", t).Format(dateFormat),
				},
			},
			err: "",
		},
		{
			config: Config{
				DatabaseConfig: &DatabaseConfig{
					Driver: MysqlDriver,
					URL:    env,
				},
				Datasets: []Dataset{
					{
						SQL: "SELECT app_name, updated_at FROM builds ORDER BY app_name LIMIT 4",
						Fields: []Field{
							{Name: "App", Type: StringType},
							{Name: "Day", Type: DatetimeType},
						},
					},
				},
			},
			out: []map[string]interface{}{
				{
					"app": "",
					"day": parseTime("2017-03-23T16:45:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "everdeen",
					"day": parseTime("2017-04-23T11:14:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-04-23T13:42:00Z", t).Format(time.RFC3339),
				},
				{
					"app": "geckoboard-ruby",
					"day": parseTime("2017-03-23T17:11:00Z", t).Format(time.RFC3339),
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
