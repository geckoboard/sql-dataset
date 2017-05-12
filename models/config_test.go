package models

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		config Config
		err    []string
	}{
		{
			Config{
				GeckoboardAPIKey: "",
			},
			[]string{
				"Geckoboard api key is required",
				"Database config is required",
			},
		},
		{
			Config{
				GeckoboardAPIKey: "",
				DatabaseConfig:   &DatabaseConfig{},
			},
			[]string{
				"Geckoboard api key is required",
				"Database driver is required",
			},
		},
		{
			Config{
				GeckoboardAPIKey: "123",
				DatabaseConfig: &DatabaseConfig{
					Driver: "pear",
					URL:    "pear://localhost/test",
				},
			},
			[]string{
				"Unsupported driver 'pear' only [mysql postgres sqlite3] are supported",
			},
		},
		{
			Config{
				GeckoboardAPIKey: "1234-12345",
				DatabaseConfig: &DatabaseConfig{
					Driver: PostgresDriver,
					URL:    "mysql://localhost/testdb",
				},
			},
			nil,
		},
		{
			Config{
				GeckoboardAPIKey: "1234-12345",
				RefreshTimeSec:   120,
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    "mysql://localhost/testdb",
				},
			},
			nil,
		},
		{
			Config{
				GeckoboardAPIKey: "1234-12345",
				RefreshTimeSec:   120,
				DatabaseConfig: &DatabaseConfig{
					Driver: MySQLDriver,
					URL:    "mysql://localhost/testdb",
				},
				Datasets: []Dataset{
					{
						Name:       "users.count",
						UpdateType: "wrong",
						SQL:        "fake sql",
						Fields:     []Field{{Name: "count", Type: "number"}},
					},
				},
			},
			[]string{
				"Dataset update type must be append or replace",
			},
		},
	}

	for i, tc := range testCases {
		err := tc.config.Validate()

		if tc.err == nil && err != nil {
			t.Errorf("[%d] Expected no error but got %s", i, err)
		}

		if tc.err != nil && err == nil {
			t.Errorf("[%d] Expected error %s but got none", i, tc.err)
		}

		if len(err) != len(tc.err) {
			t.Errorf("[%d] Expected error count %d but got %d", i, len(tc.err), len(err))
		}

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("[%d] Expected errors %s but got %s", i, tc.err, err)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		in     string
		config *Config
		err    string
	}{
		{
			"",
			nil,
			"File path is required to load config",
		},
		{
			filepath.Join("fixtures", "invalid_config.yml"),
			nil,
			"Error occurred parsing the config: yaml: did not find expected key",
		},
		{
			filepath.Join("fixtures", "valid_config.yml"),
			&Config{
				GeckoboardAPIKey: "1234dsfd21322",
				DatabaseConfig: &DatabaseConfig{
					Driver:   PostgresDriver,
					Username: "root",
					Password: "pass234",
					Host:     "/var/postgres/POSTGRES.5543",
					Protocol: "unix",
					Database: "someDB",
					TLSConfig: &TLSConfig{
						KeyFile:  "path/test.key",
						CertFile: "path/test.crt",
					},
					Params: map[string]string{
						"charset": "utf-8",
					},
				},
				RefreshTimeSec: 60,
				Datasets: []Dataset{
					{
						Name:       "active.users.by.org.plan",
						UpdateType: Replace,
						SQL:        "SELECT o.plan_type, count(*) user_count FROM users u, organisation o where o.user_id = u.id AND o.plan_type <> 'trial' order by user_count DESC limit 10",
						Fields: []Field{
							{Name: "count", Type: NumberType},
							{Name: "org", Type: StringType, Key: "custom_org"},
						},
					},
				},
			},
			"",
		},
		{
			filepath.Join("fixtures", "valid_config2.yml"),
			&Config{
				GeckoboardAPIKey: "1234dsfd21322",
				DatabaseConfig: &DatabaseConfig{
					Driver:   PostgresDriver,
					Host:     "fake-host",
					Port:     "5433",
					Database: "someDB",
					TLSConfig: &TLSConfig{
						CAFile:  "path/cert.pem",
						SSLMode: "verify-full",
					},
				},
				RefreshTimeSec: 60,
				Datasets: []Dataset{
					{
						Name:       "active.users.by.org.plan",
						UpdateType: Replace,
						SQL:        "SELECT o.plan_type, count(*) user_count FROM users u, organisation o where o.user_id = u.id AND o.plan_type <> 'trial' order by user_count DESC limit 10",
						Fields: []Field{
							{Name: "count", Type: NumberType},
							{Name: "org", Type: StringType},
							{Name: "Total Earnings", Type: MoneyType, CurrencyCode: "USD"},
						},
					},
				},
			},
			"",
		},
	}

	for i, tc := range testCases {
		c, err := LoadConfig(tc.in)

		if tc.err == "" && err != nil {
			t.Errorf("[%d] Expected no error but got %s", i, err)
			continue
		}

		if !reflect.DeepEqual(tc.config, c) {
			t.Errorf("[%d] Expected config %#v but got %#v", i, tc.config, c)
		}

		if err != nil && tc.err != err.Error() {
			t.Errorf("[%d] Expected error %s but got %s", i, tc.err, err.Error())
		}
	}
}

func TestFieldKeyValue(t *testing.T) {
	ds := Dataset{
		Fields: []Field{
			{
				Key:  "customKey",
				Name: "Percent Complete",
				Type: PercentageType,
			},
			{
				Name: "Total Cost",
				Type: MoneyType,
			},
		},
	}

	customKey := "customKey"
	normalKey := "total_cost"

	if key := ds.Fields[0].KeyValue(); key != customKey {
		t.Errorf("Expected keyvalue '%s' but got '%s'", customKey, key)
	}

	if key := ds.Fields[1].KeyValue(); key != normalKey {
		t.Errorf("Expected keyvalue '%s' but got '%s'", normalKey, key)
	}
}
