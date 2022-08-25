package eppoclient

import "fmt"

type IAssignmentLogger interface {
	LogAssignment(event map[string]string)
}

type AssignmentLogger struct {
}

func NewAssignmentLogger() IAssignmentLogger {
	return &AssignmentLogger{}
}

func (al *AssignmentLogger) LogAssignment(event map[string]string) {
	fmt.Println("Assignment Logged")
	fmt.Println(event)
}
