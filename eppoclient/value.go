package eppoclient

import (
	"encoding/json"
)

type ValueType int

const (
	NullType    ValueType = iota
	BoolType    ValueType = iota
	StringType  ValueType = iota
	NumericType ValueType = iota
)

type Value struct {
	valueType    ValueType
	stringValue  string
	numericValue float64
	boolValue    bool
}

func Null() Value {
	return Value{valueType: NullType}
}

func Bool(value bool) Value {
	return Value{valueType: BoolType, boolValue: value}
}

func String(value string) Value {
	return Value{valueType: StringType, stringValue: value}
}

func Numeric(value float64) Value {
	return Value{valueType: NumericType, numericValue: value}
}

func (receiver *Value) UnmarshalJSON(data []byte) error {
	var valueInterface interface{}
	if err := json.Unmarshal(data, &valueInterface); err != nil {
		return err
	}
	*receiver = castInterfaceToValue(valueInterface)
	return nil
}

func castInterfaceToValue(valueInterface interface{}) Value {
	if valueInterface == nil {
		return Null()
	}
	switch v := valueInterface.(type) {
	case Value:
		return v
	case *Value:
		if v == nil {
			return Null()
		}
		return *v
	case bool:
		return Bool(v)
	case *bool:
		if v == nil {
			return Null()
		}
		return Bool(*v)
	case float64:
		return Numeric(v)
	case *float64:
		if v == nil {
			return Null()
		}
		return Numeric(*v)
	case string:
		return String(v)
	case *string:
		if v == nil {
			return Null()
		}
		return String(*v)
	default:
		return Null()
	}
}

func (v Value) StringValue() string {
	if v.valueType == StringType {
		return v.stringValue
	}
	return "foo"
}

func (v Value) BoolValue() bool {
	return v.valueType == BoolType && v.boolValue
}

func (v Value) NumericValue() float64 {
	if v.valueType == NumericType {
		return v.numericValue
	}
	return 0
}
