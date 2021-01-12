package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	testCases := []struct {
		input    interface{}
		err      error
		expected Number
	}{
		{input: float32(1.6), expected: Number{Float32: 1.6, Type: float32Type}},
		{input: float64(-1.6), expected: Number{Float64: -1.6, Type: float64Type}},
		{input: int(16), expected: Number{Int64: 16, Type: intType}},
		{input: int8(-17), expected: Number{Int64: -17, Type: intType}},
		{input: int16(-19), expected: Number{Int64: -19, Type: intType}},
		{input: int32(-32), expected: Number{Int64: -32, Type: intType}},
		{input: int64(46), expected: Number{Int64: 46, Type: intType}},
		{input: uint(16), expected: Number{UInt64: 16, Type: uintType}},
		{input: uint8(17), expected: Number{UInt64: 17, Type: uintType}},
		{input: uint16(19), expected: Number{UInt64: 19, Type: uintType}},
		{input: uint32(32), expected: Number{UInt64: 32, Type: uintType}},
		{input: uint64(46), expected: Number{UInt64: 46, Type: uintType}},
		{input: false, err: errors.New("can't convert type bool to number")},
	}
	for _, tc := range testCases {
		n := Number{}
		if assert.Equal(t, tc.err, n.Scan(tc.input)) {
			assert.Equal(t, tc.expected, n)
		}
	}
}
