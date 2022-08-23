package eppoclient

var mockVariations = []Variation{
	{Name: "control", ShardRange: ShardRange{Start: 0, End: 34}},
	{Name: "variant-1", ShardRange: ShardRange{Start: 34, End: 67}},
	{Name: "variant-2", ShardRange: ShardRange{Start: 67, End: 100}},
}

var mockLogger = NewAssignmentLogger()
var testConfig = Config{apiKey: "dummy", baseUrl: "http://127.0.0.1:4000", assignmentLogger: mockLogger}

type MockConfigRequestor struct {
}

func (mcr *MockConfigRequestor) GetConfiguration(experimentKey string) (ExperimentConfiguration, error) {
	overrides := make(Dictionary)

	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 10000}},
	}

	result := ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 0,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}

	return result, nil
}

func (mcr *MockConfigRequestor) FetchAndStoreConfigurations() {
}

type MockConfigRequestor100PercentExposure struct {
}

func (mcr *MockConfigRequestor100PercentExposure) GetConfiguration(experimentKey string) (ExperimentConfiguration, error) {
	overrides := make(Dictionary)

	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 10000}},
	}

	result := ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}

	return result, nil
}

func (mcr *MockConfigRequestor100PercentExposure) FetchAndStoreConfigurations() {
}
