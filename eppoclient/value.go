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
	type Alias Value // Alias to avoid recursion
	return json.Marshal(&struct {
		ValueType ValueType `json:"valueType"`
		*Alias
	}{
		ValueType: v.ValueType,
		Alias:     (*Alias)(&v),
	})
}

func (v *Value) UnmarshalJSON(data []byte) error {
	// Temporary struct to capture the JSON structure
	var temp struct {
		ValueType    ValueType `json:"valueType"`
		BoolValue    *bool     `json:"boolValue,omitempty"`
		NumericValue *float64  `json:"numericValue,omitempty"`
		StringValue  *string   `json:"stringValue,omitempty"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	v.ValueType = temp.ValueType

	switch temp.ValueType {
	case BoolType:
		if temp.BoolValue != nil {
			v.BoolValue = *temp.BoolValue
		} else {
			return fmt.Errorf("bool value missing for BoolType")
		}
	case NumericType:
		if temp.NumericValue != nil {
			v.NumericValue = *temp.NumericValue
		} else {
			return fmt.Errorf("numeric value missing for NumericType")
		}
	case StringType:
		if temp.StringValue != nil {
			v.StringValue = *temp.StringValue
		} else {
			return fmt.Errorf("string value missing for StringType")
		}
	case NullType:
		// Handle NullType if necessary
	default:
		return fmt.Errorf("unsupported value type")
	}

	return nil
}
