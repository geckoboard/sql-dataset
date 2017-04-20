package models

import "fmt"

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
	Name       string      `json:"name"`
	UpdateType DatasetType `json:"update_type"`
	SQL        string      `json:"sql"`
	Fields     []Field     `json:"fields"`
}

type Field struct {
	Type FieldType `json:"type"`
	Name string    `json:"name"`
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

	return errors
}
