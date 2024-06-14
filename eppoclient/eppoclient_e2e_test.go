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
	"github.com/stretchr/testify/mock"
)

const TEST_DATA_DIR = "test-data/ufc/tests"
const MOCK_UFC_RESPONSE_FILE = "test-data/ufc/flags-v1.json"

type testData struct {
	Flag          string        `json:"flag"`
	VariationType variationType `json:"variationType"`
	DefaultValue  interface{}   `json:"defaultValue"`
	Subjects      []struct {
		SubjectKey        string            `json:"subjectKey"`
		SubjectAttributes SubjectAttributes `json:"subjectAttributes"`
		Assignment        interface{}       `json:"assignment"`
	} `json:"subjects"`
}

var tstData = map[string]testData{}

func Test_e2e(t *testing.T) {
	serverUrl := initFixture()

	mockLogger := new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()
	client := InitClient(Config{BaseUrl: serverUrl, SdkKey: "dummy", AssignmentLogger: mockLogger})

	// give client the time to "fetch" the mock config
	time.Sleep(2 * time.Second)

	for name, test := range tstData {
		t.Run(name, func(t *testing.T) {
			for _, subject := range test.Subjects {
				t.Run(subject.SubjectKey, func(t *testing.T) {
					switch test.VariationType {
					case booleanVariation:
						value, _ := client.GetBoolAssignment(test.Flag, subject.SubjectKey, subject.SubjectAttributes, test.DefaultValue.(bool))
						assert.Equal(t, subject.Assignment, value)
					case numericVariation:
						value, _ := client.GetNumericAssignment(test.Flag, subject.SubjectKey, subject.SubjectAttributes, test.DefaultValue.(float64))
						assert.Equal(t, subject.Assignment, value)
					case integerVariation:
						value, _ := client.GetIntegerAssignment(test.Flag, subject.SubjectKey, subject.SubjectAttributes, int64(test.DefaultValue.(float64)))
						assert.Equal(t, int64(subject.Assignment.(float64)), value)

					case jsonVariation:
						value, _ := client.GetJSONAssignment(test.Flag, subject.SubjectKey, subject.SubjectAttributes, test.DefaultValue)
						assert.Equal(t, subject.Assignment, value)
					case stringVariation:
						value, _ := client.GetStringAssignment(test.Flag, subject.SubjectKey, subject.SubjectAttributes, test.DefaultValue.(string))
						assert.Equal(t, subject.Assignment, value)
					default:
						panic(fmt.Sprintf("unknown variation: %v", test.VariationType))
					}
				})
			}
		})
	}
}

func initFixture() string {
	testResponse := getTestData() // this is here because we need to append to global testData

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/flag-config/v1/config":
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

func getTestData() ufcResponse {
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
		tstData[file.Name()] = testCaseDict
	}

	var ufcResponse ufcResponse
	ufcResponseJsonFile, _ := os.Open(MOCK_UFC_RESPONSE_FILE)
	byteValue, _ := io.ReadAll(ufcResponseJsonFile)
	err = json.Unmarshal(byteValue, &ufcResponse)
	if err != nil {
		fmt.Println("Error reading mock UFC response file")
	}

	if err != nil {
		fmt.Println("Error reading mock UFC response file")
	}

	return ufcResponse
}
