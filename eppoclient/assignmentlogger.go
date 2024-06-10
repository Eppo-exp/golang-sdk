package eppoclient

import "fmt"

type IAssignmentLogger interface {
	LogAssignment(event AssignmentEvent)
}

type AssignmentEvent struct {
	Experiment        string            `json:"experiment"`
	FeatureFlag       string            `json:"featureFlag"`
	Allocation        string            `json:"allocation"`
	Variation         string            `json:"variation"`
	Subject           string            `json:"subject"`
	SubjectAttributes SubjectAttributes `json:"subjectAttributes,omitempty"`
	Timestamp         string            `json:"timestamp"`
	MetaData          map[string]string `json:"metaData"`
	ExtraLogging      map[string]string `json:"extraLogging,omitempty"`
}

type AssignmentLogger struct {
}

func NewAssignmentLogger() IAssignmentLogger {
	return &AssignmentLogger{}
}

func (al *AssignmentLogger) LogAssignment(event AssignmentEvent) {
	fmt.Println("Assignment Logged")
	fmt.Println(event)
}
