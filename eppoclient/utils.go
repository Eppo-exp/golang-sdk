package eppoclient

type dictionary map[string]interface{}

type testData struct {
	Experiment             string                  `json:"experiment"`
	PercentExposure        float32                 `json:"percentExposure"`
	Variations             []testDataVariations    `json:"variations"`
	Subjects               []string                `json:"subjects"`
	SubjectsWithAttributes []subjectWithAttributes `json:"subjectsWithAttributes"`
	ExpectedAssignments    []string                `json:"expectedAssignments"`
}

type subjectWithAttributes struct {
	SubjectKey        string     `json:"subjectKey"`
	SubjectAttributes dictionary `json:"subjectAttributes"`
}

type testDataVariations struct {
	Name       string             `json:"name"`
	ShardRange testDataShardRange `json:"shardRange"`
}

type testDataShardRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
