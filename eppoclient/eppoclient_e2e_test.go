package eppoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TEST_DATA_DIR = "test-data/assignment-v2"
const MOCK_RAC_RESPONSE_FILE = "test-data/rac-experiments-v3.json"

var tstData = []testData{}

func Test_e2e(t *testing.T) {
	serverUrl := initFixture()

	asmntLogger := &AssignmentLogger{}
	client := InitClient(Config{BaseUrl: serverUrl, ApiKey: "dummy", AssignmentLogger: asmntLogger})

	time.Sleep(2 * time.Second)

	for _, experiment := range tstData {
		expName := experiment.Experiment

		booleanAssignments := []bool{}
		jsonAssignments := []string{}
		numericAssignments := []float64{}
		stringAssignments := []string{}

		for _, subject := range experiment.SubjectsWithAttributes {
			switch experiment.ValueType {
			case "boolean":
				booleanAssignment, err := client.GetBoolAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if err == nil {
					assert.Nil(t, err)
				}

				booleanAssignments = append(booleanAssignments, booleanAssignment)
			case "numeric":
				numericAssignment, err := client.GetNumericAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if err == nil {
					assert.Nil(t, err)
				}

				numericAssignments = append(numericAssignments, numericAssignment)
			case "json":
				jsonAssignment, err := client.GetJSONStringAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if err == nil {
					assert.Nil(t, err)
				}

				jsonAssignments = append(jsonAssignments, jsonAssignment)
			case "string":
				stringAssignment, err := client.GetStringAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if err == nil {
					assert.Nil(t, err)
				}

				stringAssignments = append(stringAssignments, stringAssignment)
			}
		}

		for _, subject := range experiment.Subjects {
			switch experiment.ValueType {
			case "boolean":
				booleanAssignment, err := client.GetBoolAssignment(subject, expName, dictionary{})
				if err == nil {
					assert.Nil(t, err)
				}

				booleanAssignments = append(booleanAssignments, booleanAssignment)
			case "json":
				jsonAssignment, err := client.GetJSONStringAssignment(subject, expName, dictionary{})
				if err == nil {
					assert.Nil(t, err)
				}

				jsonAssignments = append(jsonAssignments, jsonAssignment)
			case "numeric":
				numericAssignment, err := client.GetNumericAssignment(subject, expName, dictionary{})
				if err == nil {
					assert.Nil(t, err)
				}

				numericAssignments = append(numericAssignments, numericAssignment)
			case "string":
				stringAssignment, err := client.GetStringAssignment(subject, expName, dictionary{})

				if err == nil {
					assert.Nil(t, err)
				}

				stringAssignments = append(stringAssignments, stringAssignment)
			}
		}

		switch experiment.ValueType {
		case "boolean":
			expectedAssignments := []bool{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.BoolValue)
			}
			assert.Equal(t, expectedAssignments, booleanAssignments)
		case "json":
			expectedAssignments := []string{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.StringValue)
			}
			assert.Equal(t, expectedAssignments, jsonAssignments)
		case "numeric":
			expectedAssignments := []float64{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.NumericValue)
			}
			assert.Equal(t, expectedAssignments, numericAssignments)
		case "string":
			expectedAssignments := []string{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.StringValue)
			}
			assert.Equal(t, expectedAssignments, stringAssignments)
		}
	}
}

func initFixture() string {
	testResponse := getTestData() // this is here because we need to append to global testData

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/randomized_assignment/v3/config":
			err := json.NewEncoder(w).Encode(testResponse)
			if err != nil {
				fmt.Println("Error encoding test response")
			}
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	return server.URL
}

func getTestData() dictionary {
	files, err := os.ReadDir(TEST_DATA_DIR)

	if err != nil {
		panic("test cases files read error")
	}

	for _, file := range files {
		jsonFile, _ := os.Open(TEST_DATA_DIR + "/" + file.Name())

		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		testCaseDict := testData{}
		byteValue, _ := io.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, &testCaseDict)
		if err != nil {
			fmt.Println("Error reading test case file")
		}
		tstData = append(tstData, testCaseDict)
	}

	var racResponseData map[string]interface{}
	racResponseJsonFile, _ := os.Open(MOCK_RAC_RESPONSE_FILE)
	byteValue, _ := io.ReadAll(racResponseJsonFile)
	err = json.Unmarshal(byteValue, &racResponseData)
	if err != nil {
		fmt.Println("Error reading mock RAC response file")
	}

	if err != nil {
		fmt.Println("Error reading mock RAC response file")
	}
	return racResponseData
}
