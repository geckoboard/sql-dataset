package models

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBuildSchemaFields(t *testing.T) {
	testCases := []struct {
		in  *Dataset
		out *Dataset
		err string
	}{
		{
			// Field name unaltered
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "count",
						Type: "datetime",
					},
				},
			},
			out: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "count",
						Type: "datetime",
					},
				},
				SchemaFields: map[string]Field{
					"count": Field{
						Name: "count",
						Type: "datetime",
					},
				},
			},
		},
		{
			// Field has key provided
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "Count All",
						Key:  "not_matching",
						Type: "datetime",
					},
				},
			},
			out: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "Count All",
						Key:  "not_matching",
						Type: "datetime",
					},
				},
				SchemaFields: map[string]Field{
					"not_matching": Field{
						Name: "Count All",
						Key:  "not_matching",
						Type: "datetime",
					},
				},
			},
		},
		{
			// Multiple fields
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "Count All",
						Type: "datetime",
					},
					{
						Name: "Service",
						Key:  "newkey",
						Type: "string",
					},
				},
			},
			out: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "Count All",
						Type: "datetime",
					},
					{
						Name: "Service",
						Key:  "newkey",
						Type: "string",
					},
				},
				SchemaFields: map[string]Field{
					"count_all": Field{
						Name: "Count All",
						Type: "datetime",
					},
					"newkey": {
						Name: "Service",
						Key:  "newkey",
						Type: "string",
					},
				},
			},
		},
		{
			// Unique not matching any fields
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "count",
						Type: "datetime",
					},
				},
				UniqueBy: []string{"count", "blah"},
			},
			out: &Dataset{},
			err: "Following unique by 'blah' for dataset 'users.count' has no matching field",
		},
		{
			// Unique by errors when the user doesn't use custom key supplied
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Key:  "name_of_appy",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"APP name", "Count All"},
			},
			out: &Dataset{},
			err: "Following unique by 'APP name' for dataset 'users.count' has no matching field",
		},
		{
			// Unique by errors with the users original input
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Key:  "name_of_appy",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"App name", "Count All"},
			},
			out: &Dataset{},
			err: "Following unique by 'App name' for dataset 'users.count' has no matching field",
		},
		{
			// Unique by converted correctly to match generated field keys
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"App name", "Count All"},
			},
			out: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"app_name", "count_all"},
				SchemaFields: map[string]Field{
					"app_name": Field{
						Name: "App name",
						Type: "string",
					},
					"count_all": Field{
						Name: "Count All",
						Type: "number",
					},
				},
			},
		},
		{
			// Unique by works with users supplied custom key
			in: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Key:  "name_of_appy",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"name_of_appy", "Count All"},
			},
			out: &Dataset{
				Name: "users.count",
				Fields: []Field{
					{
						Name: "App name",
						Key:  "name_of_appy",
						Type: "string",
					},
					{
						Name: "Count All",
						Type: "number",
					},
				},
				UniqueBy: []string{"name_of_appy", "count_all"},
				SchemaFields: map[string]Field{
					"name_of_appy": Field{
						Name: "App name",
						Key:  "name_of_appy",
						Type: "string",
					},
					"count_all": Field{
						Name: "Count All",
						Type: "number",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		err := tc.in.BuildSchemaFields()

		if tc.err == "" && err != nil {
			t.Errorf("[%d] Expected no error but got %s", i, err)
		}

		if tc.err != "" && err == nil {
			t.Errorf("[%d] Expected error %s but got none", i, tc.err)
		}

		if err != nil && tc.err != err.Error() {
			t.Errorf("[%d] Expected error %s but got %s", i, tc.err, err)
		}

		if tc.err == "" && !reflect.DeepEqual(tc.in, tc.out) {
			t.Errorf("[%d] Expected dataset %#v but got %#v", i, tc.in, tc.out)
		}
	}
}

func TestDatasetValidate(t *testing.T) {
	testCases := []struct {
		dataset Dataset
		err     []string
	}{
		{
			Dataset{},
			[]string{
				errMissingDatasetName,
				fmt.Sprintf(errInvalidDatasetUpdateType, ""),
				errMissingDatasetSQL,
				errMissingDatasetFields,
			},
		},
		{
			Dataset{Fields: []Field{{}}},
			[]string{
				errMissingDatasetName,
				fmt.Sprintf(errInvalidDatasetUpdateType, ""),
				errMissingDatasetSQL,
				fmt.Sprintf(errInvalidFieldType, "", fieldTypes),
				errMissingFieldName,
			},
		},
		{
			Dataset{
				Name:       "c",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{
				Name:       "cd",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{
				Name:       ".bbc",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{
				Name:       "ABCwat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{Name: "abc wat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{
				Name:       "-wat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{errInvalidDatasetName},
		},
		{
			Dataset{
				Name:       "users.count",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: "numbre"}},
			},
			[]string{fmt.Sprintf(errInvalidFieldType, "numbre", fieldTypes)},
		},
		{
			Dataset{
				Name:       "some.dataset",
				UpdateType: Replace,
				SQL:        "SELECT 1;",
				Fields: []Field{
					{
						Name: "Unique",
						Type: NumberType,
					},
					{
						Name: "counts",
						Type: NumberType,
					},
					{
						Name: "Count's",
						Type: NumberType,
					},
					{
						Name:         "Total Cost",
						Type:         MoneyType,
						CurrencyCode: "USD",
					},
				},
			},
			[]string{fmt.Sprintf(errDuplicateFieldNames, "Count's")},
		},
		{
			Dataset{
				Name:       "some.dataset",
				UpdateType: Replace,
				SQL:        "SELECT 1;",
				Fields: []Field{
					{
						Name: "Unique",
						Type: NumberType,
					},
					{
						Name: "counts",
						Type: NumberType,
					},
					{
						Name: "Count's",
						Type: NumberType,
					},
					{
						Name:         "Total Cost.",
						Type:         MoneyType,
						CurrencyCode: "USD",
					},
					{
						Name:         "Total C.o.S.t",
						Type:         MoneyType,
						CurrencyCode: "USD",
					},
				},
			},
			[]string{fmt.Sprintf(errDuplicateFieldNames, `Count's", "Total C.o.S.t`)},
		},
		{
			Dataset{
				Name:       "app",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType, CurrencyCode: "USD"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "a-abc",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType, CurrencyCode: "USD"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "abc_abc",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType, CurrencyCode: "USD"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "12abc",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType, CurrencyCode: "USD"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "app.build.cost",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType}},
			},
			[]string{errMissingCurrency},
		},
		{
			Dataset{
				Name:       "app.build.cost",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: MoneyType, CurrencyCode: "USD"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "users.count",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			nil,
		},
		{
			Dataset{
				Name:       "users.count",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: "datetime"}},
			},
			nil,
		},
	}

	for i, tc := range testCases {
		err := tc.dataset.Validate()

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
			t.Errorf("[%d] Expected errors %#v but got %#v", i, tc.err, err)
		}
	}
}

func TestFieldKeyValue(t *testing.T) {
	testCases := []struct {
		field Field
		out   string
	}{
		{
			Field{
				Key:  "customKey",
				Name: "Percent Complete",
				Type: PercentageType,
			},
			"customKey",
		},
		{
			Field{
				Name: "Total Cost",
				Type: MoneyType,
			},
			"total_cost",
		},
		{
			Field{
				Name: "Total's",
				Type: MoneyType,
			},
			"totals",
		},
		{
			Field{
				Name: "MRR. Tot",
				Type: MoneyType,
			},
			"mrr_tot",
		},
		{
			Field{
				Name: "_MRR. T-",
				Type: MoneyType,
			},
			"mrr_t",
		},
		{
			Field{
				Name: "_MRR. Tot_",
				Type: MoneyType,
			},
			"mrr_tot",
		},
		{
			Field{
				Name: "Random Names'",
				Type: MoneyType,
			},
			"random_names",
		},
		{
			Field{
				Name: "2nd stage",
				Type: MoneyType,
			},
			"2nd_stage",
		},
		{
			Field{
				Name: " extra whitespace ",
				Type: MoneyType,
			},
			"extra_whitespace",
		},
		{
			Field{
				Name: "  extra  whitespace   ",
				Type: MoneyType,
			},
			"extra__whitespace",
		},
		{
			// Let the server validate length
			Field{
				Name: "mr",
				Type: MoneyType,
			},
			"mr",
		},
	}

	for _, tc := range testCases {
		if key := tc.field.KeyValue(); key != tc.out {
			t.Errorf("Expected keyvalue '%s' but got '%s'", tc.out, key)
		}
	}
}
