package drivers

import (
	"path/filepath"
	"testing"

	"github.com/geckoboard/sql-dataset/models"
)

func TestNewDSNBuilder(t *testing.T) {
	testCases := []struct {
		in  models.DatabaseConfig
		out string
		err string
	}{
		//SQLite Driver
		{
			in: models.DatabaseConfig{
				Driver: models.SQLiteDriver,
			},
			err: "Database is required for a connection",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.SQLiteDriver,
				Database: "models/fixtures/db.sqlite",
			},
			out: "models/fixtures/db.sqlite",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.SQLiteDriver,
				Database: "dir/db.sqlite",
				Password: "blah123",
				Params: map[string]string{
					"cache": "shared",
					"mode":  "rwc",
				},
			},
			out: "file:dir/db.sqlite?password=blah123&cache=shared&mode=rwc",
		},
		//Mysql Driver
		{
			in: models.DatabaseConfig{
				Driver: models.MysqlDriver,
			},
			err: ErrDatabaseRequired.Error(),
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Database: "some_name",
			},
			err: ErrUsernameRequired.Error(),
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Database: "someDB",
			},
			out: "root@tcp(localhost:3306)/someDB?parseTime=true",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
			},
			out: "root:fp123@tcp(localhost:3306)/someDB?parseTime=true",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Host:     "fake-host",
			},
			out: "root:fp123@tcp(fake-host:3306)/someDB?parseTime=true",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Host:     "fake-host",
				Port:     "3366",
			},
			out: "root:fp123@tcp(fake-host:3366)/someDB?parseTime=true",
		},
		{
			//Unix socket connection
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Host:     "/tmp/mysql",
				Protocol: "unix",
			},
			out: "root:fp123@unix(/tmp/mysql)/someDB?parseTime=true",
		},
		{
			//IPv6 needs to be in square brackets
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Host:     "de:ad:be:ef::ca:fe",
			},
			out: "root:fp123@tcp([de:ad:be:ef::ca:fe]:3306)/someDB?parseTime=true",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Host:     "project-id:region:instance",
				Protocol: "cloudsql",
			},
			out: "root:fp123@cloudsql(project-id:region:instance)/someDB?parseTime=true",
		},
		{
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				Params: map[string]string{
					"charset": "utf8mb4,utf8",
					"loc":     "US/Pacific",
				},
			},
			out: "root:fp123@tcp(localhost:3306)/someDB?charset=utf8mb4,utf8&loc=US/Pacific&parseTime=true",
		},
		{
			// ca cert file path
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				TLSConfig: &models.TLSConfig{
					CAFile: filepath.Join("..", "models", "fixtures", "ca.cert.pem"),
				},
			},
			out: "root:fp123@tcp(localhost:3306)/someDB?parseTime=true&tls=customCert",
		},
		{
			// invalid ca cert file
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				TLSConfig: &models.TLSConfig{
					CAFile: filepath.Join("..", "models", "fixtures", "ca.key.pem"),
				},
			},
			err: "Failed to append PEM, is it a valid ca cert ?",
		},
		{
			// ssl only
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				TLSConfig: &models.TLSConfig{
					SSLMode: "true",
				},
			},
			out: "root:fp123@tcp(localhost:3306)/someDB?parseTime=true&tls=true",
		},
		{
			// key and cert file path
			in: models.DatabaseConfig{
				Driver:   models.MysqlDriver,
				Username: "root",
				Password: "fp123",
				Database: "someDB",
				TLSConfig: &models.TLSConfig{
					KeyFile:  filepath.Join("..", "models", "fixtures", "test.key"),
					CertFile: filepath.Join("..", "models", "fixtures", "test.crt"),
				},
			},
			out: "root:fp123@tcp(localhost:3306)/someDB?parseTime=true&tls=customCert",
		},
	}

	for i, tc := range testCases {
		n := NewDSNBuilder(tc.in.Driver)
		dsn, err := n.Build(&tc.in)

		if tc.err == "" && err != nil {
			t.Errorf("[%d] Expected no error but got %s", i, err)
		}

		if tc.err != "" && err == nil {
			t.Errorf("[%d] Expected error %s but got nothing", i, tc.err)
		}

		if err != nil && tc.err != err.Error() {
			t.Errorf("[%d] Expected error %s but got %s", i, tc.err, err)
		}

		if dsn != tc.out {
			t.Errorf("[%d] Expected dsn connection string '%s' but got '%s'", i, tc.out, dsn)
		}
	}
}