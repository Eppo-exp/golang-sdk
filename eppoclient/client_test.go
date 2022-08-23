package eppoclient

import (
	"testing"
)

var mockLogger = NewAssignmentLogger()
var testConfig = Config{apiKey: "dummy", baseUrl: "http://127.0.0.1:4000", assignmentLogger: mockLogger}

var mockVariations = []Variation{
	{Name: "control", ShardRange: ShardRange{Start: 0, End: 34}},
	{Name: "variant-1", ShardRange: ShardRange{Start: 34, End: 67}},
	{Name: "variant-2", ShardRange: ShardRange{Start: 67, End: 100}},
}

func Test_AssignBlankExperiment(t *testing.T) {
	// No need to check whether `recover()` is nil. Just turn off the panic.
	defer func() { _ = recover() }()

	InitClient(testConfig)
	var mockConfigRequestor = &MockConfigRequestor{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	client.GetAssignment("subject-1", "", Dictionary{})
	// Never reaches here if `GetAssignment` panics.
	t.Errorf("did not panic")
}

func Test_AssignBlankSubject(t *testing.T) {
	// No need to check whether `recover()` is nil. Just turn off the panic.
	defer func() { _ = recover() }()

	InitClient(testConfig)
	var mockConfigRequestor = &MockConfigRequestor{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	client.GetAssignment("", "experiment-1", Dictionary{})
	// Never reaches here if `GetAssignment` panics.
	t.Errorf("did not panic")
}

// func Test_SubjectNotInSample(t *testing.T) {
// 	InitClient(testConfig)
// 	var mockConfigRequestor = &MockConfigRequestor{}

// 	client := NewEppoClient(mockConfigRequestor, mockLogger)

// 	client.GetAssignment("user-1", "experiment-key-1", Dictionary{})
// 	// Never reaches here if `GetAssignment` panics.
// 	t.Errorf("did not panic")
// }

// def test_assign_subject_not_in_sample(mock_config_requestor):
//     mock_config_requestor.get_configuration.return_value = ExperimentConfigurationDto(
//         subjectShards=10000,
//         percentExposure=0,
//         enabled=True,
//         variations=[
//             VariationDto(name="control", shardRange=ShardRange(start=0, end=10000))
//         ],
//         name="recommendation_algo",
//         overrides=dict(),
//     )
//     client = EppoClient(
//         config_requestor=mock_config_requestor, assignment_logger=AssignmentLogger()
//     )
//     assert client.get_assignment("user-1", "experiment-key-1") is None

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
