package eppoclient

import (
	"context"
	"fmt"
)

// IAssignmentLogger is be deprecated at next major version
// and replaced by IAssignmentLoggerContext
type IAssignmentLogger interface {
	LogAssignment(event AssignmentEvent)
}

type IAssignmentLoggerContext interface {
	LogAssignment(context.Context, AssignmentEvent)
}

// BanditActionLogger is going to be merged into IAssignmentLogger in
// the next major version.
type BanditActionLogger interface {
	LogBanditAction(event BanditEvent)
}

// TODO: in the next major release, upgrade Timestamp fields to time.Time.

type AssignmentEvent struct {
	Experiment        string            `json:"experiment"`
	FeatureFlag       string            `json:"featureFlag"`
	Allocation        string            `json:"allocation"`
	Variation         string            `json:"variation"`
	Subject           string            `json:"subject"`
	SubjectAttributes Attributes        `json:"subjectAttributes,omitempty"`
	Timestamp         string            `json:"timestamp"`
	MetaData          map[string]string `json:"metaData"`
	ExtraLogging      map[string]string `json:"extraLogging,omitempty"`
}
type BanditEvent struct {
	FlagKey                      string             `json:"flagKey"`
	BanditKey                    string             `json:"banditKey"`
	Subject                      string             `json:"subject"`
	Action                       string             `json:"action,omitempty"`
	ActionProbability            float64            `json:"actionProbability,omitempty"`
	OptimalityGap                float64            `json:"optimalityGap,omitempty"`
	ModelVersion                 string             `json:"modelVersion,omitempty"`
	Timestamp                    string             `json:"timestamp"`
	SubjectNumericAttributes     map[string]float64 `json:"subjectNumericAttributes,omitempty"`
	SubjectCategoricalAttributes map[string]string  `json:"subjectCategoricalAttributes,omitempty"`
	ActionNumericAttributes      map[string]float64 `json:"actionNumericAttributes,omitempty"`
	ActionCategoricalAttributes  map[string]string  `json:"actionCategoricalAttributes,omitempty"`
	MetaData                     map[string]string  `json:"metaData"`
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

func (al *AssignmentLogger) LogBanditAction(event BanditEvent) {
	fmt.Println("Bandit Action Logged")
	fmt.Println(event)
}
