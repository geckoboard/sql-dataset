package models

import (
	"fmt"
	"regexp"
	"strings"
)

type DatasetType string
type FieldType string

const (
	Append  DatasetType = "append"
	Replace DatasetType = "replace"

	NumberType     FieldType = "number"
	DateType       FieldType = "date"
	DatetimeType   FieldType = "datetime"
	MoneyType      FieldType = "money"
	PercentageType FieldType = "percentage"
	StringType     FieldType = "string"
)

var (
	datasetNameRegexp = regexp.MustCompile(`(?)^[0-9a-z][0-9a-z._\-]{1,}[0-9a-z]$`)
	fieldIdRegexp     = regexp.MustCompile(`[^a-z0-9_ ]+|[_\W]+$|^[_\W]+`)
)

var fieldTypes = []FieldType{
	NumberType,
	DateType,
	DatetimeType,
	MoneyType,
	PercentageType,
	StringType,
}

type Dataset struct {
	Name         string           `json:"id"                   yaml:"name"`
	UpdateType   DatasetType      `json:"-"                    yaml:"update_type"`
	UniqueBy     []string         `json:"unique_by,omitempty"  yaml:"unique_by,omitempty"`
	SQL          string           `json:"-"                    yaml:"sql"`
	Fields       []Field          `json:"-"                    yaml:"fields"`
	SchemaFields map[string]Field `json:"fields"               yaml:"-"`
}

type Field struct {
	Type         FieldType `json:"type"                     yaml:"type"`
	Key          string    `json:"-"                        yaml:"key"`
	Name         string    `json:"name"                     yaml:"name"`
	CurrencyCode string    `json:"currency_code,omitempty"  yaml:"currency_code"`
	Optional     bool      `json:"optional,omitempty"       yaml:"optional,omitempty"`
}

// KeyValue returns the field key if present
// otherwise by default returns the field name underscored
func (f Field) KeyValue() string {
	if f.Key != "" {
		return f.Key
	}

	// Remove any characters not alphanumeric or space and replace with nothing
	key := fieldIdRegexp.ReplaceAllString(strings.ToLower(f.Name), "")

	return strings.Replace(key, " ", "_", -1)
}

// BuildSchemaFields creates a map[string]Field of
// the dataset fields for sending over to Geckoboard
func (ds *Dataset) BuildSchemaFields() error {
	if ds.SchemaFields != nil {
		return nil
	}

	fields := make(map[string]Field)

	for _, f := range ds.Fields {
		fields[f.KeyValue()] = f
	}

	ds.SchemaFields = fields
	return ds.updateUniqueByKeys(fields)
}

func (ds *Dataset) updateUniqueByKeys(newFields map[string]Field) error {
	for i, ub := range ds.UniqueBy {
		matched := false
		for k := range newFields {
			if ub == k {
				matched = true
				break
			}
		}

		// Do the same the under the hood key update to match
		// the new field key that may have been generated
		if !matched {
			key := fieldIdRegexp.ReplaceAllString(strings.ToLower(ds.UniqueBy[i]), "")
			ds.UniqueBy[i] = strings.Replace(key, " ", "_", -1)
		}
	}

	// Now double check all unique by have matching field
	for _, ub := range ds.UniqueBy {
		matched := false
		for k := range newFields {
			if ub == k {
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf(
				"Following unique by %s for dataset %s has no matching field",
				ub,
				ds.Name,
			)
		}
	}

	return nil
}

func (ds Dataset) Validate() (errors []string) {
	if ds.Name == "" {
		errors = append(errors, errMissingDatasetName)
	}

	if ds.Name != "" && !datasetNameRegexp.MatchString(ds.Name) {
		errors = append(errors, errInvalidDatasetName)
	}

	if ds.UpdateType != Append && ds.UpdateType != Replace {
		errors = append(errors,
			fmt.Sprintf(errInvalidDatasetUpdateType, ds.UpdateType))
	}

	if ds.SQL == "" {
		errors = append(errors, errMissingDatasetSQL)
	}

	if len(ds.Fields) == 0 {
		errors = append(errors, errMissingDatasetFields)
	}

	for _, f := range ds.Fields {
		errors = append(errors, f.Validate()...)
	}

	if err := ds.validateGeneratedFieldKeysUnique(); err != "" {
		errors = append(errors, err)
	}

	return errors
}

func (f Field) Validate() (errors []string) {
	validType := false

	for _, t := range fieldTypes {
		if t == f.Type {
			validType = true
			break
		}
	}

	if !validType {
		errors = append(errors, fmt.Sprintf(errInvalidFieldType, f.Type, fieldTypes))
	}

	if f.Name == "" {
		errors = append(errors, errMissingFieldName)
	}

	if f.Type == MoneyType && f.CurrencyCode == "" {
		errors = append(errors, errMissingCurrency)
	}

	return errors
}

func (ds Dataset) validateGeneratedFieldKeysUnique() string {
	uniqueNameMap := make(map[string]interface{})
	var names []string

	for _, f := range ds.Fields {
		k := f.KeyValue()
		if uniqueNameMap[k] == nil {
			uniqueNameMap[k] = struct{}{}
			continue
		}

		names = append(names, f.Name)
	}

	if len(names) > 0 {
		return fmt.Sprintf(errDuplicateFieldNames, strings.Join(names, `", "`))
	}

	return ""
}
