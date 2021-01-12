package models

import (
	"fmt"
	"strconv"
)

const (
	intType     = "int"
	uintType    = "uint"
	float32Type = "float32"
	float64Type = "float64"
)

type Number struct {
	Int64   int64
	UInt64  uint64
	Float32 float32
	Float64 float64

	Type string
}

func (n *Number) Value(optional bool) interface{} {
	switch n.Type {
	case intType:
		return n.Int64
	case uintType:
		return n.UInt64
	case float32Type:
		return n.Float32
	case float64Type:
		return n.Float64
	default:
		if optional {
			return nil
		}

		return 0
	}
}

func (n *Number) Scan(value interface{}) error {
	switch val := value.(type) {
	default:
		return fmt.Errorf("can't convert type %T to number", value)
	case float64:
		n.Type = float64Type
		n.Float64 = val
	case float32:
		n.Type = float32Type
		n.Float32 = val
	case int:
		n.Type = intType
		n.Int64 = int64(val)
	case int8:
		n.Type = intType
		n.Int64 = int64(val)
	case int16:
		n.Type = intType
		n.Int64 = int64(val)
	case int32:
		n.Type = intType
		n.Int64 = int64(val)
	case int64:
		n.Type = intType
		n.Int64 = val
	case uint:
		n.Type = uintType
		n.UInt64 = uint64(val)
	case uint8:
		n.Type = uintType
		n.UInt64 = uint64(val)
	case uint16:
		n.Type = uintType
		n.UInt64 = uint64(val)
	case uint32:
		n.Type = uintType
		n.UInt64 = uint64(val)
	case uint64:
		n.Type = uintType
		n.UInt64 = val
	case []byte:
		return n.pruneBytes(value.([]byte))
	}
	return nil
}

func (n *Number) pruneBytes(value []byte) error {
	var floatPrecision uint8

	if len(value) == 0 {
		return nil
	}

	for i, b := range value {
		if b == 46 {
			floatPrecision = uint8(len(value[i:])) - 1
			break
		}
	}

	if floatPrecision == 0 {
		i, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			return err
		}

		n.Type = intType
		n.Int64 = i
	} else {
		f, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return err
		}

		n.Type = float64Type
		n.Float64 = f
	}

	return nil
}
