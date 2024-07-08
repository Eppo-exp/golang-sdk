package eppoclient

import "time"

type banditResponse struct {
	Bandits   map[string]banditConfiguration `json:"bandits"`
	UpdatedAt time.Time                      `json:"updatedAt"`
}

type banditConfiguration struct {
	BanditKey    string          `json:"banditKey"`
	ModelName    string          `json:"modelName"`
	ModelVersion string          `json:"modelVersion"`
	ModelData    banditModelData `json:"modelData"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type banditModelData struct {
	Gamma                  float64                       `json:"gamma"`
	DefaultActionScore     float64                       `json:"defaultActionScore"`
	ActionProbabilityFloor float64                       `json:"actionProbabilityFloor"`
	Coefficients           map[string]banditCoefficients `json:"coefficients"`
}

type banditCoefficients struct {
	ActionKey                      string                                  `json:"actionKey"`
	Intercept                      float64                                 `json:"intercept"`
	SubjectNumericCoefficients     []banditNumericAttributeCoefficient     `json:"subjectNumericCoefficients"`
	SubjectCategoricalCoefficients []banditCategoricalAttributeCoefficient `json:"subjectCategoricalCoefficients"`
	ActionNumericCoefficients      []banditNumericAttributeCoefficient     `json:"actionNumericCoefficients"`
	ActionCategoricalCoefficients  []banditCategoricalAttributeCoefficient `json:"actionCategoricalCoefficients"`
}

type banditCategoricalAttributeCoefficient struct {
	AttributeKey            string             `json:"attributeKey"`
	MissingValueCoefficient float64            `json:"missingValueCoefficient"`
	ValueCoefficients       map[string]float64 `json:"valueCoefficients"`
}

type banditNumericAttributeCoefficient struct {
	AttributeKey            string  `json:"attributeKey"`
	Coefficient             float64 `json:"coefficient"`
	MissingValueCoefficient float64 `json:"missingValueCoefficient"`
}

type banditVariation struct {
	Key            string `json:"key"`
	FlagKey        string `json:"flagKey"`
	VariationKey   string `json:"variationKey"`
	VariationValue string `json:"variationValue"`
}
