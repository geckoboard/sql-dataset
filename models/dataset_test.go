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
