package eppoclient

import (
	"encoding/json"
	"fmt"
	"time"

	semver "github.com/Masterminds/semver/v3"
)

type configResponse struct {
	Flags   map[string]*flagConfiguration `json:"flags"`
	Bandits map[string][]banditVariation  `json:"bandits,omitempty"`
}

func (response *configResponse) precompute() {
	for i := range response.Flags {
		response.Flags[i].precompute()
	}
}

type flagConfiguration struct {
	Key           string               `json:"key"`
	Enabled       bool                 `json:"enabled"`
	VariationType variationType        `json:"variationType"`
	Variations    map[string]variation `json:"variations"`
	Allocations   []allocation         `json:"allocations"`
	TotalShards   int64                `json:"totalShards"`
	// Cached Variations parsed according to `VariationType`.
	//
	// Types are as follows:
	// - STRING -> string
	// - NUMERIC -> float64
	// - INTEGER -> int64
	// - BOOLEAN -> bool
	// - JSON -> interface{}
	ParsedVariations map[string]interface{} `json:"-"`
}

func (flag *flagConfiguration) precompute() {
	for i := range flag.Allocations {
		flag.Allocations[i].precompute()
	}

	flag.ParsedVariations = make(map[string]interface{}, len(flag.Variations))
	for i := range flag.Variations {
		value, err := flag.VariationType.parseVariationValue(flag.Variations[i].Value)
		if err == nil {
			flag.ParsedVariations[i] = value
		}
	}
}

type variation struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type variationType int

const (
	stringVariation variationType = iota
	integerVariation
	numericVariation
	booleanVariation
	jsonVariation
)

func (v variationType) MarshalJSON() ([]byte, error) {
	switch v {
	case stringVariation:
		return json.Marshal("STRING")
	case integerVariation:
		return json.Marshal("INTEGER")
	case numericVariation:
		return json.Marshal("NUMERIC")
	case booleanVariation:
		return json.Marshal("BOOLEAN")
	case jsonVariation:
		return json.Marshal("JSON")
	default:
		return nil, fmt.Errorf("unsupported variation type: %d", v)
	}
}

func (v *variationType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value {
	case "STRING":
		*v = stringVariation
	case "INTEGER":
		*v = integerVariation
	case "NUMERIC":
		*v = numericVariation
	case "BOOLEAN":
		*v = booleanVariation
	case "JSON":
		*v = jsonVariation
	default:
		return fmt.Errorf("unknown variation type: %s", value)
	}

	return nil
}

func (ty variationType) parseVariationValue(value json.RawMessage) (interface{}, error) {
	switch ty {
	case stringVariation:
		var s string
		err := json.Unmarshal(value, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	case integerVariation:
		var i int64
		err := json.Unmarshal(value, &i)
		if err != nil {
			return nil, err
		}
		return i, nil
	case numericVariation:
		var f float64
		err := json.Unmarshal(value, &f)
		if err != nil {
			return nil, err
		}
		return f, nil
	case booleanVariation:
		var b bool
		err := json.Unmarshal(value, &b)
		if err != nil {
			return nil, err
		}
		return b, nil
	case jsonVariation:
		var s string
		err := json.Unmarshal(value, &s)
		if err != nil {
			return nil, err
		}

		raw := []byte(s)

		var parsed interface{}
		err = json.Unmarshal(raw, &parsed)
		if err != nil {
			return nil, err
		}

		return parsed, nil
	default:
		return nil, fmt.Errorf("unexpected variation type: %v", ty)
	}
}

type allocation struct {
	Key     string    `json:"key"`
	Rules   []rule    `json:"rules"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
	Splits  []split   `json:"splits"`
	DoLog   *bool     `json:"doLog"`
}

func (a *allocation) precompute() {
	for i := range a.Rules {
		a.Rules[i].precompute()
	}
}

type rule struct {
	Conditions []condition `json:"conditions"`
}

func (r *rule) precompute() {
	for i := range r.Conditions {
		r.Conditions[i].precompute()
	}
}

type condition struct {
	Operator  string      `json:"operator"`
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`

	NumericValue      float64
	NumericValueValid bool
	SemVerValue       *semver.Version
	SemVerValueValid  bool
}

func (c *condition) precompute() {
	// Try to convert Value to a float64
	if num, err := toFloat64(c.Value); err == nil {
		c.NumericValue = num
		c.NumericValueValid = true
		return
	}

	// Try to convert Value to a string and then parse as semver
	if str, ok := c.Value.(string); ok {
		if semVer, err := semver.NewVersion(str); err == nil {
			c.SemVerValue = semVer
			c.SemVerValueValid = true
			return
		}
	}

	c.NumericValueValid = false
	c.SemVerValueValid = false
}

type split struct {
	Shards       []shard           `json:"shards"`
	VariationKey string            `json:"variationKey"`
	ExtraLogging map[string]string `json:"extraLogging"`
}

type shard struct {
	Salt   string       `json:"salt"`
	Ranges []shardRange `json:"ranges"`
}

type shardRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
