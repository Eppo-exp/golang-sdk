package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var (
	zapLogger, _      = zap.NewDevelopment()
	applicationLogger = NewZapLogger(zapLogger)
)

func Test_AssignBlankExperiment(t *testing.T) {
	var mockLogger = new(mockLogger)
	client := newEppoClient(newConfigurationStore(configuration{}), nil, nil, mockLogger, applicationLogger)

	_, err := client.GetStringAssignment("", "subject-1", Attributes{}, "")
	assert.Error(t, err)
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockLogger = new(mockLogger)
	client := newEppoClient(newConfigurationStore(configuration{}), nil, nil, mockLogger, applicationLogger)

	_, err := client.GetStringAssignment("experiment-1", "", Attributes{}, "")
	assert.Error(t, err)
}

func Test_LogAssignment(t *testing.T) {
	tests := []struct {
		name          string
		doLog         *bool
		expectedCalls int
	}{
		{
			name:          "DoLog key is absent",
			doLog:         nil,
			expectedCalls: 1,
		},
		{
			name:          "DoLog key is present but false",
			doLog:         &[]bool{false}[0],
			expectedCalls: 0,
		},
		{
			name:          "DoLog key is present and true",
			doLog:         &[]bool{true}[0],
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockLogger = new(mockLogger)
			mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

			config := configResponse{
				Flags: map[string]flagConfiguration{
					"experiment-key-1": flagConfiguration{
						Key:           "experiment-key-1",
						Enabled:       true,
						TotalShards:   10000,
						VariationType: stringVariation,
						Variations: map[string]variation{
							"control": variation{
								Key:   "control",
								Value: "control",
							},
						},
						Allocations: []allocation{
							{
								Key:   "allocation-key",
								DoLog: tt.doLog,
								Splits: []split{
									{
										VariationKey: "control",
										Shards: []shard{
											{
												Salt: "",
												Ranges: []shardRange{
													{
														Start: 0,
														End:   10000,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			client := newEppoClient(newConfigurationStore(configuration{flags: config}), nil, nil, mockLogger, applicationLogger)

			assignment, err := client.GetStringAssignment("experiment-key-1", "user-1", Attributes{}, "")
			expected := "control"

			assert.Nil(t, err)
			assert.Equal(t, expected, assignment)
			mockLogger.AssertNumberOfCalls(t, "LogAssignment", tt.expectedCalls)
		})
	}
}

func Test_client_loggerIsCalledWithProperBanditEvent(t *testing.T) {
	var logger = new(mockLogger)
	logger.Mock.On("LogAssignment", mock.Anything).Return()
	logger.Mock.On("LogBanditAction", mock.Anything).Return()

	flags := configResponse{
		Bandits: map[string][]banditVariation{
			"bandit": []banditVariation{
				banditVariation{
					Key:            "bandit",
					FlagKey:        "testFlag",
					VariationKey:   "bandit",
					VariationValue: "bandit",
				},
			},
		},
	}
	bandits := banditResponse{
		Bandits: map[string]banditConfiguration{
			"bandit": {
				BanditKey:    "bandit",
				ModelName:    "falcon",
				ModelVersion: "v123",
				ModelData: banditModelData{
					Gamma:                  0,
					DefaultActionScore:     0,
					ActionProbabilityFloor: 0,
					Coefficients:           map[string]banditCoefficients{},
				},
			},
		},
	}

	client := newEppoClient(newConfigurationStore(configuration{flags: flags, bandits: bandits}), nil, nil, logger, applicationLogger)
	actions := map[string]ContextAttributes{
		"action1": {},
	}
	client.GetBanditAction("testFlag", "subject", ContextAttributes{}, actions, "bandit")

	event := logger.Calls[0].Arguments[0].(BanditEvent)
	assert.Equal(t, "testFlag", event.FlagKey)
	assert.Equal(t, "bandit", event.BanditKey)
	assert.Equal(t, "subject", event.Subject)
	assert.Equal(t, "action1", event.Action)
}

func Test_GetStringAssignmentHandlesLoggingPanic(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Panic("logging panic")

	config := configResponse{Flags: map[string]flagConfiguration{
		"experiment-key-1": flagConfiguration{
			Key:           "experiment-key-1",
			Enabled:       true,
			TotalShards:   10000,
			VariationType: stringVariation,
			Variations: map[string]variation{
				"control": variation{
					Key:   "control",
					Value: "control",
				},
			},
			Allocations: []allocation{
				{
					Key: "allocation-key",
					Splits: []split{
						{
							VariationKey: "control",
							Shards: []shard{
								{
									Salt: "",
									Ranges: []shardRange{
										{
											Start: 0,
											End:   10000,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}}

	client := newEppoClient(newConfigurationStore(configuration{flags: config}), nil, nil, mockLogger, applicationLogger)

	assignment, err := client.GetStringAssignment("experiment-key-1", "user-1", Attributes{}, "")
	expected := "control"

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}

func Test_client_handlesBanditLoggerPanic(t *testing.T) {
	var logger = new(mockLogger)
	logger.Mock.On("LogAssignment", mock.Anything).Return()
	logger.Mock.On("LogBanditAction", mock.Anything).Panic("logging panic")

	flags := configResponse{
		Bandits: map[string][]banditVariation{
			"bandit": []banditVariation{
				banditVariation{
					Key:            "bandit",
					FlagKey:        "testFlag",
					VariationKey:   "bandit",
					VariationValue: "bandit",
				},
			},
		},
	}
	bandits := banditResponse{
		Bandits: map[string]banditConfiguration{
			"bandit": {
				BanditKey:    "bandit",
				ModelName:    "falcon",
				ModelVersion: "v123",
				ModelData: banditModelData{
					Gamma:                  0,
					DefaultActionScore:     0,
					ActionProbabilityFloor: 0,
					Coefficients:           map[string]banditCoefficients{},
				},
			},
		},
	}

	client := newEppoClient(newConfigurationStore(configuration{flags: flags, bandits: bandits}), nil, nil, logger, applicationLogger)
	actions := map[string]ContextAttributes{
		"action1": {},
	}
	client.GetBanditAction("testFlag", "subject", ContextAttributes{}, actions, "bandit")

	logger.AssertNumberOfCalls(t, "LogBanditAction", 1)
}

func Test_client_correctActionIsReturnedIfBanditLoggerPanics(t *testing.T) {
	var logger = new(mockLogger)
	logger.Mock.On("LogAssignment", mock.Anything).Return()
	logger.Mock.On("LogBanditAction", mock.Anything).Panic("logging panic")

	flags := configResponse{
		Bandits: map[string][]banditVariation{
			"bandit": []banditVariation{
				banditVariation{
					Key:            "bandit",
					FlagKey:        "testFlag",
					VariationKey:   "bandit",
					VariationValue: "bandit",
				},
			},
		},
	}
	bandits := banditResponse{
		Bandits: map[string]banditConfiguration{
			"bandit": {
				BanditKey:    "bandit",
				ModelName:    "falcon",
				ModelVersion: "v123",
				ModelData: banditModelData{
					Gamma:                  0,
					DefaultActionScore:     0,
					ActionProbabilityFloor: 0,
					Coefficients:           map[string]banditCoefficients{},
				},
			},
		},
	}

	client := newEppoClient(newConfigurationStore(configuration{flags: flags, bandits: bandits}), nil, nil, logger, applicationLogger)
	actions := map[string]ContextAttributes{
		"action1": {},
	}
	result := client.GetBanditAction("testFlag", "subject", ContextAttributes{}, actions, "bandit")

	expectedAction := "action1"
	assert.Equal(t, BanditResult{
		Variation: "bandit",
		Action:    &expectedAction,
	}, result)
}
