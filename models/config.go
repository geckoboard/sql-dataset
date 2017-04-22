package models

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Driver string

const (
	PostgresDriver Driver = "postgresql"
	MysqlDriver    Driver = "mysql"
)

var supportedDrivers = []Driver{PostgresDriver, MysqlDriver}
var errParseConfigFile = "Error occurred parsing the config: %s"

type Config struct {
	GeckoboardAPIKey string          `yaml:"geckoboard_api_key"`
	DatabaseConfig   *DatabaseConfig `yaml:"database_config"`
	RefreshTimeSec   int32           `yaml:"refresh_time_sec"`
	Datasets         []Dataset       `yaml:"datasets"`
}

// DatabaseConfig holds the db type, url
// and other custom options such as tls config
type DatabaseConfig struct {
	Driver Driver `yaml:"driver"`
	URL    string `yaml:"url"`
}

func LoadConfig(filepath string) (config *Config, err error) {
	var b []byte

	if filepath == "" {
		return nil, errors.New("File path is required to load config")
	}

	if b, err = ioutil.ReadFile(filepath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf(errParseConfigFile, err)
	}

	return config, nil
}

func (c Config) Validate() (errors []string) {
	if c.GeckoboardAPIKey == "" {
		errors = append(errors, "Geckoboard api key is required")
	}

	if c.DatabaseConfig == nil {
		errors = append(errors, "Database config is required")
	} else {
		errors = append(errors, c.DatabaseConfig.Validate()...)
	}

	for _, ds := range c.Datasets {
		errors = append(errors, ds.Validate()...)
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
