package eppoclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		assignments := []string{}
		numericAssignments := []float64{}

		for _, subject := range experiment.SubjectsWithAttributes {
			if experiment.ValueType == "numeric" {
				numericAssignment, err := client.GetNumericAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if err == nil {
					assert.Nil(t, err)
				}

				numericAssignments = append(numericAssignments, numericAssignment)
			} else {
				assignment, err := client.GetStringAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)
				if assignment != "" {
					assert.Nil(t, err)
				}

				assignments = append(assignments, assignment)
			}
		}

		for _, subject := range experiment.Subjects {
			if experiment.ValueType == "numeric" {
				numericAssignment, err := client.GetNumericAssignment(subject, expName, dictionary{})
				if err == nil {
					assert.Nil(t, err)
				}

				numericAssignments = append(numericAssignments, numericAssignment)
			} else {
				assignment, err := client.GetStringAssignment(subject, expName, dictionary{})

				if assignment != "" {
					assert.Nil(t, err)
				}

				assignments = append(assignments, assignment)
			}
		}

		if experiment.ValueType == "numeric" {
			expectedAssignments := []float64{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.numericValue)
			}
			assert.Equal(t, expectedAssignments, numericAssignments)
		} else {
			expectedAssignments := []string{}
			for _, assignment := range experiment.ExpectedAssignments {
				expectedAssignments = append(expectedAssignments, assignment.stringValue)
			}
			assert.Equal(t, expectedAssignments, assignments)
		}
	}
}

func initFixture() string {
	testResponse := getTestData() // this is here because we need to append to global testData

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/randomized_assignment/v3/config":
			json.NewEncoder(w).Encode(testResponse)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	return server.URL
}

func getTestData() dictionary {
	files, err := ioutil.ReadDir(TEST_DATA_DIR)

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
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &testCaseDict)
		tstData = append(tstData, testCaseDict)
	}

	var racResponseData map[string]interface{}
	racResponseJsonFile, _ := os.Open(MOCK_RAC_RESPONSE_FILE)
	byteValue, _ := ioutil.ReadAll(racResponseJsonFile)
	err = json.Unmarshal(byteValue, &racResponseData)
	if err != nil {
		fmt.Println("Error reading mock RAC response file")
	}
	return racResponseData
}
