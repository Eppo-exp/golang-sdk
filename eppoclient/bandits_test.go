package eppoclient

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type banditTest struct {
	Flag         string
	DefaultValue string
	Subjects     []struct {
		SubjectKey        string
		SubjectAttributes struct {
			Numeric     map[string]float64 `json:"numericAttributes"`
			Categorical map[string]string  `json:"categoricalAttributes"`
		}
		Actions []struct {
			ActionKey             string
			NumericAttributes     map[string]float64
			CategoricalAttributes map[string]string
		}
		Assignment BanditResult
	}
}

func Test_InferContextAttributes(t *testing.T) {
	attributes := Attributes{
		"string": "blah",
		"int":    42,
		"bool":   true,
	}
	contextAttributes := InferContextAttributes(attributes)

	expected := ContextAttributes{
		Numeric: map[string]float64{
			"int": 42.0,
		},
		Categorical: map[string]string{
			"string": "blah",
			"bool":   "true",
		},
	}

	assert.Equal(t, expected, contextAttributes)
}

func Test_bandits_sdkTestData(t *testing.T) {
	flags := readJsonFile[configResponse]("test-data/ufc/bandit-flags-v1.json")
	bandits := readJsonFile[banditResponse]("test-data/ufc/bandit-models-v1.json")
	configStore := newConfigurationStoreWithConfig(configuration{
		flags:   flags,
		bandits: bandits,
	})
	logger := new(mockLogger)
	logger.Mock.On("LogAssignment", mock.Anything).Return()
	logger.Mock.On("LogBanditAction", mock.Anything).Return()
	client := newEppoClient(configStore, nil, nil, logger, applicationLogger)

	tests := readJsonDirectory[banditTest]("test-data/ufc/bandit-tests/")
	for file, test := range tests {
		t.Run(file, func(t *testing.T) {
			for _, subject := range test.Subjects {
				t.Run(subject.SubjectKey, func(t *testing.T) {
					actions := make(map[string]ContextAttributes)
					for _, a := range subject.Actions {
						actions[a.ActionKey] = ContextAttributes{
							Numeric:     a.NumericAttributes,
							Categorical: a.CategoricalAttributes,
						}
					}

					result := client.GetBanditAction(
						test.Flag,
						subject.SubjectKey,
						ContextAttributes{
							Numeric:     subject.SubjectAttributes.Numeric,
							Categorical: subject.SubjectAttributes.Categorical,
						},
						actions,
						test.DefaultValue)

					assert.Equal(t, subject.Assignment, result)
				})
			}
		})
	}
}

func readJsonFile[T any](filePath string) T {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var target T
	err = json.Unmarshal(byteValue, &target)
	if err != nil {
		panic(err)
	}

	return target
}

func readJsonDirectory[T any](dirPath string) map[string]T {
	results := make(map[string]T)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" && !strings.Contains(filepath.Base(path), ".dynamic-typing.") {
			results[filepath.Base(path)] = readJsonFile[T](path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return results
}
