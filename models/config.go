package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	ClickHouseDriver = "clickhouse"
	MSSQLDriver      = "mssql"
	MySQLDriver      = "mysql"
	PostgresDriver   = "postgres"
	SQLiteDriver     = "sqlite3"
)

var (
	SupportedDrivers = []string{ClickHouseDriver, MSSQLDriver, MySQLDriver, PostgresDriver, SQLiteDriver}
	interpolateRegex = regexp.MustCompile(`{{\s*([a-zA-Z0-9_]+)\s*}}`)
)

type Config struct {
	GeckoboardAPIKey string          `yaml:"geckoboard_api_key"`
	DatabaseConfig   *DatabaseConfig `yaml:"database"`
	RefreshTimeSec   uint16          `yaml:"refresh_time_sec"`
	Datasets         []Dataset       `yaml:"datasets"`
}

// DatabaseConfig holds the db type, url
// and other custom options such as tls config
type DatabaseConfig struct {
	Driver    string            `yaml:"driver"`
	URL       string            `yaml:"-"`
	Host      string            `yaml:"host"`
	Port      string            `yaml:"port"`
	Protocol  string            `yaml:"protocol"`
	Database  string            `yaml:"name"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	TLSConfig *TLSConfig        `yaml:"tls_config"`
	Params    map[string]string `yaml:"params"`
}

type TLSConfig struct {
	KeyFile  string `yaml:"key_file"`
	CertFile string `yaml:"cert_file"`
	CAFile   string `yaml:"ca_file"`
	SSLMode  string `yaml:"ssl_mode"`
}

func LoadConfig(filepath string) (config *Config, err error) {
	var b []byte

	if filepath == "" {
		return nil, errors.New(errNoConfigFound)
	}

	if b, err = ioutil.ReadFile(filepath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf(errParseConfigFile, err)
	}

	config.replaceSupportedInterpolatedValues()

	return config, nil
}

func (c Config) Validate() (errors []string) {
	if c.GeckoboardAPIKey == "" {
		errors = append(errors, errMissingAPIKey)
	}

	if c.DatabaseConfig == nil {
		errors = append(errors, errMissingDBConfig)
	} else {
		errors = append(errors, c.DatabaseConfig.Validate()...)
	}

	if len(c.Datasets) == 0 {
		errors = append(errors, errNoDatasets)
	}

	for _, ds := range c.Datasets {
		errors = append(errors, ds.Validate()...)
	}

	return errors
}

func (dc DatabaseConfig) Validate() (errors []string) {
	if dc.Driver == "" {
		errors = append(errors, errMissingDBDriver)
	} else {
		var matched bool

		for _, d := range SupportedDrivers {
			if d == dc.Driver {
				matched = true
				break
			}
		}

		if !matched {
			errors = append(errors, fmt.Sprintf(errDriverNotSupported, dc.Driver, SupportedDrivers))
		}
	}

	return errors
}

func (c *Config) replaceSupportedInterpolatedValues() {
	c.GeckoboardAPIKey = convertEnvToValue(c.GeckoboardAPIKey)

	if c.DatabaseConfig != nil {
		dc := c.DatabaseConfig

		dc.Username = convertEnvToValue(dc.Username)
		dc.Password = convertEnvToValue(dc.Password)
		dc.Host = convertEnvToValue(dc.Host)
		dc.Database = convertEnvToValue(dc.Database)
		dc.Port = convertEnvToValue(dc.Port)
	}
}

func convertEnvToValue(value string) string {
	if value == "" {
		return ""
	}

	keys := interpolateRegex.FindStringSubmatch(value)

	if len(keys) != 2 {
		return value
	}

	return os.Getenv(keys[1])
}
