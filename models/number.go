package models

import (
	"fmt"
	"strconv"
)

const (
	intType     = "int"
	float32Type = "float32"
	float64Type = "float64"
)

type Number struct {
	Int64   int64
	Float32 float32
	Float64 float64

	Type string
}

//TODO: Allow null for optional fields that support it in Geckoboard
func (n *Number) Value() interface{} {
	switch n.Type {
	case intType:
		return n.Int64
	case float32Type:
		return n.Float32
	case float64Type:
		return n.Float64
	default:
		return 0
	}
}

func (n *Number) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		return fmt.Errorf("can't convert string %#v to number", value.(string))
	case float64:
		n.Type = float64Type
		n.Float64 = value.(float64)
	case float32:
		n.Type = float32Type
		n.Float32 = value.(float32)
	case int, int32, int64:
		n.Type = intType
		n.Int64 = value.(int64)
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
