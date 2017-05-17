package models

import (
	"reflect"
	"testing"
)

func TestDatasetValidate(t *testing.T) {
	testCases := []struct {
		dataset Dataset
		err     []string
	}{
		{
			Dataset{},
			[]string{
				"Dataset name is required",
				"Dataset update type must be append or replace",
				"Dataset sql is required",
				"At least one field is required for a dataset",
			},
		},
		{
			Dataset{Fields: []Field{{}}},
			[]string{
				"Dataset name is required",
				"Dataset update type must be append or replace",
				"Dataset sql is required",
				"Unknown field type '' supported field types [number date datetime money percentage string]",
				"Field name is required",
			},
		},
		{
			Dataset{
				Name:       "c",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{
				Name:       "cd",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{
				Name:       ".bbc",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{
				Name:       "ABCwat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{Name: "abc wat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{
				Name:       "-wat",
				UpdateType: Replace,
				SQL:        "SELECT 1",
				Fields:     []Field{{Name: "count", Type: "number"}},
			},
			[]string{"Dataset name is invalid, should be 3 or more characters with only lowercase alphanumeric characters, dots, hyphens, and underscores"},
		},
		{
			Dataset{
				Name:       "users.count",
				UpdateType: Replace,
				SQL:        "SELECT * FROM some_funky_table;",
				Fields:     []Field{{Name: "count", Type: "numbre"}},
			},
			[]string{
				"Unknown field type 'numbre' supported field types [number date datetime money percentage string]",
			},
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
			[]string{`The field names "Count's" will create duplicate keys. Please revise using a unique combination of letters and numbers.`},
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
			[]string{`The field names "Count's", "Total C.o.S.t" will create duplicate keys. Please revise using a unique combination of letters and numbers.`},
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
			[]string{
				"Money type field requires an ISO 4217 currency code",
			},
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
