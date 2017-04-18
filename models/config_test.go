package models

import (
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
				"Database url is required",
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
				"Unsupported driver 'pear' only [postgresql mysql] are supported",
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
					Driver: MysqlDriver,
					URL:    "mysql://localhost/testdb",
				},
			},
			nil,
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
