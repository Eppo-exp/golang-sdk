package eppoclient

import "fmt"

type AssignmentLogger struct {
}

func NewAssignmentLogger() AssignmentLogger {
	return AssignmentLogger{}
}

func (al *AssignmentLogger) LogAssignment(event string) {
	fmt.Println("Assignment Logged: " + event)
}
