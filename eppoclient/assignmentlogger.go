package eppoclient

import "fmt"

type IAssignmentLogger interface {
	LogAssignment(event AssignmentEvent)
}

type AssignmentEvent struct {
	Experiment        string
	Variation         Value
	Subject           string
	Timestamp         string
	SubjectAttributes dictionary
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
