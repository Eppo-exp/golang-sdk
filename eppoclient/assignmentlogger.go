package eppoclient

import "fmt"

type IAssignmentLogger interface {
	LogAssignment(event string)
}

type AssignmentLogger struct {
}

func NewAssignmentLogger() IAssignmentLogger {
	return &AssignmentLogger{}
}

func (al *AssignmentLogger) LogAssignment(event string) {
	fmt.Println("Assignment Logged: " + event)
}
