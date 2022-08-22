package eppoclient

import (
	"testing"
)

var client = EppoClient{}
var mockLogger = NewAssignmentLogger()
var testConfig = Config{apiKey: "dummy", baseUrl: "http://127.0.0.1:4000", assignmentLogger: mockLogger}

var mockVariations = []Variation{
	{Name: "control", ShardRange: ShardRange{Start: 0, End: 34}},
	{Name: "variant-1", ShardRange: ShardRange{Start: 34, End: 67}},
	{Name: "variant-2", ShardRange: ShardRange{Start: 67, End: 100}},
}

func Test_AssignBlankExperiment(t *testing.T) {
	InitClient(testConfig)
	var mockConfigRequestor = NewExperimentConfigurationRequestor()

	client.New(&mockConfigRequestor, mockLogger)
	client.GetAssignment("subject-1", "", Dictionary{})

	// with pytest.raises(Exception) as exc_info:
	//
	//	client.get_assignment("subject-1", "")
	//
	// assert exc_info.value.args[0] == "Invalid value for experiment_key: cannot be blank"
}

type MockConfigRequestor struct {
}

func (mcr *MockConfigRequestor) New(httpClient HttpClient, configStore ConfigurationStore) {

}

func (mcr *MockConfigRequestor) GetConfiguration(experimentKey string) (ExperimentConfiguration, error) {
	overrides := make(Dictionary)
	overrides["a90ea45116d251a43da56e03d3dd7275"] = "variant-2"

	result := ExperimentConfiguration{
		Name:            "experiment_5",
		PercentExposure: 1,
		Enabled:         true,
		SubjectShards:   100,
		Overrides:       overrides,
		Variations:      mockVariations,
	}

	return result, nil
}

func (mcr *MockConfigRequestor) FetchAndStoreConfigurations() {
}
