package eppoclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const TEST_DATA_DIR = "test-data/assignment"
const BUCKET_NAME = "sdk-test-data"

var testData = []TestData{}

func Test_e2e(t *testing.T) {
	serverUrl := initFixture()

	client := InitClient(Config{BaseUrl: serverUrl, ApiKey: "dummy", AssignmentLogger: AssignmentLogger{}})

	time.Sleep(2 * time.Second)

	for _, experiment := range testData {
		expName := experiment.Experiment

		assignments := []string{}

		for _, subject := range experiment.Subjects {
			assignment, err := client.GetAssignment(subject, expName, Dictionary{})

			if assignment != "" {
				assert.Nil(t, err)
			}

			assignments = append(assignments, assignment)
		}

		assert.Equal(t, experiment.ExpectedAssignments, assignments)
	}
}

func initFixture() string {
	downloadTestData()
	testResponse := getTestData() // this is here because we need to append to global testData

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/randomized_assignment/config":
			json.NewEncoder(w).Encode(testResponse)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	return server.URL
}

func getTestData() Dictionary {
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

		testCaseDict := TestData{}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &testCaseDict)
		testData = append(testData, testCaseDict)
	}

	expConfigs := Dictionary{}

	for _, experimentTest := range testData {
		experimentName := experimentTest.Experiment
		expMap := Dictionary{}
		expMap["subjectShards"] = 10000
		expMap["enabled"] = true
		expMap["variations"] = experimentTest.Variations
		expMap["name"] = experimentName
		expMap["percentExposure"] = experimentTest.PercentExposure

		expConfigs[experimentName] = expMap
	}

	response := Dictionary{}
	response["experiments"] = expConfigs

	return response
}

func downloadTestData() {
	if _, err := os.Stat(TEST_DATA_DIR); os.IsNotExist(err) {
		if err := os.MkdirAll(TEST_DATA_DIR, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else {
		return //data is already downloaded, skip this step
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		fmt.Println(err)
	}

	query := &storage.Query{Prefix: "assignment/test-case"}
	bkt := client.Bucket(BUCKET_NAME)
	it := bkt.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		obj := bkt.Object(attrs.Name)
		rdr, err := obj.NewReader(ctx)

		if err != nil {
			log.Fatal(err)
		}
		defer rdr.Close()

		out, err := os.Create("test-data/" + obj.ObjectName())
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, rdr)
	}
}
