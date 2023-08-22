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
const MOCK_RAC_RESPONSE_FILE = "test-data/rac-experiments-v2.json"

var tstData = []testData{}

func Test_e2e(t *testing.T) {
	serverUrl := initFixture()

	asmntLogger := &AssignmentLogger{}
	client := InitClient(Config{BaseUrl: serverUrl, ApiKey: "dummy", AssignmentLogger: asmntLogger})

	time.Sleep(2 * time.Second)

	for _, experiment := range tstData {
		expName := experiment.Experiment

		assignments := []string{}

		for _, subject := range experiment.SubjectsWithAttributes {
			assignment, err := client.GetAssignment(subject.SubjectKey, expName, subject.SubjectAttributes)

			if assignment != "" {
				assert.Nil(t, err)
			}

			assignments = append(assignments, assignment)
		}

		for _, subject := range experiment.Subjects {
			assignment, err := client.GetAssignment(subject, expName, dictionary{})

			if assignment != "" {
				assert.Nil(t, err)
			}

			assignments = append(assignments, assignment)
		}

		assert.Equal(t, experiment.ExpectedAssignments, assignments)
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
