package models

const (
	errParseConfigFile = "There are errors in your config: %s"

	errNoConfigFound = "No config file provided. Use -config path/to/file " +
		"to specify the location of your config"

	// Config
	errMissingDBConfig    = "No database config provided."
	errDriverNotSupported = `"%s" is not a supported driver. SQL-Dataset supports %s`
	errMissingDBDriver    = "No dataset driver provided."
	errMissingAPIKey      = "No Geckoboard API key provided."

	// SQL
	errFailedSQLQuery    = "Query failed. This is the error received: %s"
	errParseSQLResultSet = "Parsing query results failed. " +
		"This is the error received: %s"

	// Dataset validations
	errNoDatasets           = "At least one dataset is required to run"
	errMissingDatasetName   = "No dataset name provided."
	errMissingDatasetSQL    = "No SQL query provided."
	errMissingDatasetFields = "No dataset fields provided."

	errInvalidDatasetName = "Invalid dataset name. Dataset names must be at " +
		"least 3 characters in length, and use only lowercase letters, " +
		"numbers, dots, hyphens, and underscores."

	errInvalidDatasetUpdateType = `"%s" is not a valid update type. ` +
		`Update type must be either append or replace.`

	// Dataset field validations
	errMissingFieldName = "No field name provided."

	errInvalidFieldType = `"%s" is not a valid field type. ` +
		`Supported field types are %s.`

	errMissingCurrency = "No currency_code provided for the money field %s. " +
		"Please provide an ISO4217 currency code."

	errDuplicateFieldNames = `The field names "%s" will create duplicate keys. ` +
		`Please revise using a unique combination of letters and numbers.`
)
