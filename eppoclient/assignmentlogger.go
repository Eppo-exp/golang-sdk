package eppoclient

import "fmt"

type IAssignmentLogger interface {
	LogAssignment(event AssignmentEvent)
}

type AssignmentEvent struct {
	Experiment        string            `json:"experiment"`
	FeatureFlag       string            `json:"featureFlag"`
	Allocation        string            `json:"allocation"`
	Variation         Value             `json:"variation"`
	Subject           string            `json:"subject"`
	Timestamp         string            `json:"timestamp"`
	SubjectAttributes SubjectAttributes `json:"subjectAttributes,omitempty"`
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
