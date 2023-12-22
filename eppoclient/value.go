package eppoclient

import (
	"encoding/json"
	"fmt"
)

type ValueType int

const (
	NullType    ValueType = iota
	BoolType    ValueType = iota
	NumericType ValueType = iota
	StringType  ValueType = iota
)

type Value struct {
	ValueType    ValueType `json:"valueType"`
	BoolValue    bool      `json:"boolValue,omitempty"`
	NumericValue float64   `json:"numericValue,omitempty"`
	StringValue  string    `json:"stringValue,omitempty"`
}

func Null() Value {
	return Value{ValueType: NullType}
}

func Bool(value bool) Value {
	return Value{ValueType: BoolType, BoolValue: value}
}

func Numeric(value float64) Value {
	return Value{ValueType: NumericType, NumericValue: value}
}

func String(value string) Value {
	return Value{ValueType: StringType, StringValue: value}
}

func (v Value) GetBoolValue() bool {
	return v.ValueType == BoolType && v.BoolValue
}

func (v Value) GetStringValue() string {
	if v.ValueType == StringType {
		return v.StringValue
	}
	return ""
}

func (v Value) GetNumericValue() float64 {
	if v.ValueType == NumericType {
		return v.NumericValue
	}

	return 0
}

func (v Value) MarshalJSON() ([]byte, error) {
	switch v.ValueType {
	case BoolType:
		return json.Marshal(v.BoolValue)
	case NumericType:
		return json.Marshal(v.NumericValue)
	case StringType:
		return json.Marshal(v.StringValue)
	case NullType:
		return json.Marshal(nil)
	default:
		return nil, fmt.Errorf("unsupported value type")
	}
}

func (v *Value) UnmarshalJSON(data []byte) error {
	// Unmarshal the data into an interface{} to check its type
	var typedValue interface{}
	if err := json.Unmarshal(data, &typedValue); err != nil {
		return err
	}

	// Determine the type of typedValue and set the Value struct accordingly
	switch value := typedValue.(type) {
	case Value:
		*v = value
	case *Value:
		if v == nil {
			*v = Null()
		}
	case string:
		*v = String(value)
	case *string:
		if v == nil {
			*v = Null()
		}
	case float64:
		// JSON numbers are float64 by default
		*v = Numeric(value)
	case bool:
		*v = Bool(value)
	case *bool:
		if v == nil {
			*v = Null()
		}
		*v = Bool(*value)
	case map[string]interface{}:
		out, _ := json.Marshal(typedValue)
		*v = String(string(out))
	case nil:
		*v = Null()
	default:
		*v = Null()
	}

	return nil
}
