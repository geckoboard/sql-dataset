package models

import (
	"fmt"
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

var fieldTypes = []FieldType{
	NumberType,
	DateType,
	DatetimeType,
	MoneyType,
	PercentageType,
	StringType,
}

type Dataset struct {
	Name       string      `json:"id"                   yaml:"name"`
	UpdateType DatasetType `json:"-"                    yaml:"update_type"`
	UniqueBy   []string    `json:"unique_by,omitempty"  yaml:"unique_by,omitempty"`
	SQL        string      `json:"-"                    yaml:"sql"`
	Fields     []Field     `json:"-"                    yaml:"fields"`
	GBFields   *GBFields   `json:"fields"               yaml:"-"`
}

type Field struct {
	Type         FieldType `json:"type"                     yaml:"type"`
	Key          string    `json:"-"                        yaml:"key"`
	Name         string    `json:"name"                     yaml:"name"`
	CurrencyCode string    `json:"currency_code,omitempty"  yaml:"currency_code"`
	Optional     bool      `json:"optional,omitempty"       yaml:"optional,omitempty"`
}

type GBFields map[string]Field

// KeyValue returns the field key if present
// otherwise by default returns the field name underscored
func (f Field) KeyValue() string {
	if f.Key != "" {
		return f.Key
	}

	return strings.ToLower(strings.Replace(f.Name, " ", "_", -1))
}

// BuildGBFields creates a map[string]Field of
// the dataset fields for sending over to Geckoboard
func (ds *Dataset) BuildGBFields() {
	if ds.GBFields != nil {
		return
	}

	fields := make(GBFields)

	for _, f := range ds.Fields {
		fields[f.KeyValue()] = f
	}

	ds.GBFields = &fields
}

func (ds Dataset) Validate() (errors []string) {
	if ds.Name == "" {
		errors = append(errors, "Dataset name is required")
	}

	if ds.UpdateType != Append && ds.UpdateType != Replace {
		errors = append(errors, "Dataset update type must be append or replace")
	}

	if ds.SQL == "" {
		errors = append(errors, "Dataset sql is required")
	}

	if len(ds.Fields) == 0 {
		errors = append(errors, "At least one field is required for a dataset")
	}

	for _, f := range ds.Fields {
		errors = append(errors, f.Validate()...)
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
		errors = append(errors, fmt.Sprintf("Unknown field type '%s' supported field types %s", f.Type, fieldTypes))
	}

	if f.Name == "" {
		errors = append(errors, "Field name is required")
	}

	if f.Type == MoneyType && f.CurrencyCode == "" {
		errors = append(errors, "Money type field requires an ISO 4217 currency code")
	}

	return errors
}
