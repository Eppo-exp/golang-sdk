package eppoclient

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_banditResponse_parseTestFile(t *testing.T) {
	banditResponse := banditResponse{}

	file, err := os.Open("test-data/ufc/bandit-models-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&banditResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, banditResponse.Bandits)
}
