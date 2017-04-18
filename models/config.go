package models

import "fmt"

type Config struct {
	GeckoboardAPIKey string          `json:"geckoboard_api_key"`
	DatabaseConfig   *DatabaseConfig `json:"database"`
	RefreshTimeSec   int32           `json:"refresh_time_sec"`
}

// DatabaseConfig holds the db type, url
// and other custom options such as tls config
type DatabaseConfig struct {
	Driver Driver
	URL    string
}

type Driver string

const (
	PostgresDriver Driver = "postgresql"
	MysqlDriver    Driver = "mysql"
)

var supportedDrivers = []Driver{PostgresDriver, MysqlDriver}

func (c Config) Validate() (errors []string) {
	if c.GeckoboardAPIKey == "" {
		errors = append(errors, "Geckoboard api key is required")
	}

	if c.DatabaseConfig == nil {
		errors = append(errors, "Database config is required")
	} else {
		errors = append(errors, c.DatabaseConfig.Validate()...)
	}

	return errors
}

func (dc DatabaseConfig) Validate() (errors []string) {
	if dc.Driver == "" {
		errors = append(errors, "Database driver is required")
	} else {
		var matched bool

		for _, d := range supportedDrivers {
			if d == dc.Driver {
				matched = true
				break
			}
		}

		if !matched {
			errors = append(errors, fmt.Sprintf("Unsupported driver '%s' only %s are supported", dc.Driver, supportedDrivers))
		}
	}

	if dc.URL == "" {
		errors = append(errors, "Database url is required")
	}

	return errors
}
