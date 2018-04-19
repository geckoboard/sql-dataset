package models

import (
	"fmt"
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
