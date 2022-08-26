package eppoclient

type dictionary map[string]interface{}

type TestData struct {
	Experiment          string               `json:"experiment"`
	PercentExposure     float32              `json:"percentExposure"`
	Variations          []TestDataVariations `json:"variations"`
	Subjects            []string             `json:"subjects"`
	ExpectedAssignments []string             `json:"expectedAssignments"`
}

type TestDataVariations struct {
	Name       string             `json:"name"`
	ShardRange TestDataShardRange `json:"shardRange"`
}

type TestDataShardRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
