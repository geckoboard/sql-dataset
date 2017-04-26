package main

import (
	"reflect"
	"testing"

	"github.com/geckoboard/sql-dataset/models"
)

func TestConfigureMySQLDSN(t *testing.T) {
	testCases := []struct {
		in    *models.Config
		dbOut *models.DatabaseConfig
		err   string
	}{
		{
			in: &models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.MysqlDriver,
					URL:    "root@/testdb",
				},
			},
			dbOut: &models.DatabaseConfig{
				Driver: models.MysqlDriver,
				URL:    "root@tcp(127.0.0.1:3306)/testdb?parseTime=true",
			},
			err: "",
		},
		{
			in: &models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.MysqlDriver,
					URL:    "root@/testdb",
					TLSConfig: &models.TLSConfig{
						CAFile: "models/fixtures/ca.cert.pem",
					},
				},
			},
			dbOut: &models.DatabaseConfig{
				Driver: models.MysqlDriver,
				URL:    "root@tcp(127.0.0.1:3306)/testdb?parseTime=true&tls=customCert",
				TLSConfig: &models.TLSConfig{
					CAFile: "models/fixtures/ca.cert.pem",
				},
			},
			err: "",
		},
		{
			in: &models.Config{
				DatabaseConfig: &models.DatabaseConfig{
					Driver: models.MysqlDriver,
					URL:    "root@/testdb",
					TLSConfig: &models.TLSConfig{
						CAFile: "models/fixtures/ca.key.pem",
					},
				},
			},
			err: "Failed to append PEM, is it a valid ca cert ?",
		},
	}

	for i, tc := range testCases {
		err := ConfigureMySQLDSN(tc.in)

		if err != nil && tc.err == "" {
			t.Errorf("[%d] Expected no error but got %s", i, err)
		}

		if err == nil && tc.err != "" {
			t.Errorf("[%d] Expected error %s but got none", i, tc.err)
		}

		if err != nil && tc.err != err.Error() {
			t.Errorf("[%d] Expected error %s but got %s", i, tc.err, err)
		}

		if tc.dbOut != nil && !reflect.DeepEqual(tc.in.DatabaseConfig, tc.dbOut) {
			t.Errorf("[%d] Expected db config %#v but got %#v", i, tc.dbOut, tc.in.DatabaseConfig)
		}
	}
}
