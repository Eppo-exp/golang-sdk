package eppoclient

import (
	"context"
	"testing"
	"time"

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
	client := newEppoClient(newConfigurationStore(), nil, nil, mockLogger, nil, applicationLogger)

	_, err := client.GetStringAssignment("", "subject-1", Attributes{}, "")
	assert.Error(t, err)
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockLogger = new(mockLogger)
	client := newEppoClient(newConfigurationStore(), nil, nil, mockLogger, nil, applicationLogger)

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
				Flags: map[string]*flagConfiguration{
					"experiment-key-1": &flagConfiguration{
						Key:           "experiment-key-1",
						Enabled:       true,
						TotalShards:   10000,
						VariationType: stringVariation,
						Variations: map[string]variation{
							"control": variation{
								Key:   "control",
								Value: []byte("\"control\""),
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

			client := newEppoClient(newConfigurationStoreWithConfig(configuration{flags: config}), nil, nil, mockLogger, nil, applicationLogger)

			assignment, err := client.GetStringAssignment("experiment-key-1", "user-1", Attributes{}, "")
			expected := "control"

			assert.Nil(t, err)
			assert.Equal(t, expected, assignment)
			mockLogger.AssertNumberOfCalls(t, "LogAssignment", tt.expectedCalls)
		})
	}
}

func Test_LogAssignmentContext(t *testing.T) {
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
			var (
				mockLoggerContext = new(mockLoggerContext)
				mockLogger        = new(mockLogger)
			)
			mockLoggerContext.Mock.
				On("LogAssignment", mock.Anything, mock.Anything).
				Return()

			config := configResponse{
				Flags: map[string]*flagConfiguration{
					"experiment-key-1": &flagConfiguration{
						Key:           "experiment-key-1",
						Enabled:       true,
						TotalShards:   10000,
						VariationType: stringVariation,
						Variations: map[string]variation{
							"control": variation{
								Key:   "control",
								Value: []byte("\"control\""),
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

			client := newEppoClient(
				newConfigurationStoreWithConfig(configuration{flags: config}),
				nil,
				nil,
				mockLogger,
				mockLoggerContext,
				applicationLogger,
			)

			ctx := context.TODO()
			assignment, err := client.GetStringAssignmentContext(ctx, "experiment-key-1", "user-1", Attributes{}, "")
			expected := "control"

			assert.Nil(t, err)
			assert.Equal(t, expected, assignment)
			mockLoggerContext.AssertNumberOfCalls(t, "LogAssignment", tt.expectedCalls)
		})
	}
}

func Test_ContextMethodsPassContext(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"

	tests := []struct {
		name          string
		variationType variationType
		variationValue []byte
		invoke        func(ctx context.Context, client *EppoClient) error
	}{
		{
			name:           "GetBoolAssignmentContext",
			variationType:  booleanVariation,
			variationValue: []byte("true"),
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetBoolAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, false)
				return err
			},
		},
		{
			name:           "GetNumericAssignmentContext",
			variationType:  numericVariation,
			variationValue: []byte("3.14"),
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetNumericAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, 0.0)
				return err
			},
		},
		{
			name:           "GetIntegerAssignmentContext",
			variationType:  integerVariation,
			variationValue: []byte("123"),
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetIntegerAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, 0)
				return err
			},
		},
		{
			name:           "GetStringAssignmentContext",
			variationType:  stringVariation,
			variationValue: []byte("\"test-string\""),
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetStringAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, "")
				return err
			},
		},
		{
			name:           "GetJSONAssignmentContext",
			variationType:  jsonVariation,
			variationValue: []byte(`"{\"key\": \"value\"}"`), // JSON string containing JSON
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetJSONAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, nil)
				return err
			},
		},
		{
			name:           "GetJSONBytesAssignmentContext",
			variationType:  jsonVariation,
			variationValue: []byte(`"{\"key\": \"value\"}"`), // JSON string containing JSON
			invoke: func(ctx context.Context, c *EppoClient) error {
				_, err := c.GetJSONBytesAssignmentContext(ctx, "flag-key", "user-1", Attributes{}, nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(mockLoggerContext)
			mockLogger.Mock.
				On("LogAssignment", mock.Anything, mock.Anything).
				Return()

			config := configResponse{
				Flags: map[string]*flagConfiguration{
					"flag-key": {
						Key:           "flag-key",
						Enabled:       true,
						TotalShards:   10000,
						VariationType: tt.variationType,
						Variations: map[string]variation{
							"control": {
								Key:   "control",
								Value: tt.variationValue,
							},
						},
						Allocations: []allocation{
							{
								Key:   "allocation-key",
								DoLog: &[]bool{true}[0],
								Splits: []split{
									{
										VariationKey: "control",
										Shards: []shard{
											{
												Salt: "",
												Ranges: []shardRange{
													{Start: 0, End: 10000},
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

			client := newEppoClient(
				newConfigurationStoreWithConfig(configuration{flags: config}),
				nil,
				nil,
				nil,
				mockLogger,
				applicationLogger,
			)

			ctx := context.WithValue(context.Background(), testKey, "test-value")
			err := tt.invoke(ctx, client)

			assert.NoError(t, err)
			mockLogger.AssertNumberOfCalls(t, "LogAssignment", 1)

			loggedCtx := mockLogger.Calls[0].Arguments[0].(context.Context)
			assert.Equal(t, "test-value", loggedCtx.Value(testKey), "context should be passed through to logger")
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

	client := newEppoClient(newConfigurationStoreWithConfig(configuration{flags: flags, bandits: bandits}), nil, nil, logger, nil, applicationLogger)
	actions := map[string]ContextAttributes{
		"action1": {},
	}
	client.GetBanditAction("testFlag", "subject", ContextAttributes{}, actions, "bandit")

	t.Log(len(logger.Calls))

	event := logger.Calls[0].Arguments[0].(BanditEvent)
	assert.Equal(t, "testFlag", event.FlagKey)
	assert.Equal(t, "bandit", event.BanditKey)
	assert.Equal(t, "subject", event.Subject)
	assert.Equal(t, "action1", event.Action)
}

func Test_client_loggerContextIsCalledWithProperBanditEvent(t *testing.T) {
	var (
		logger        = new(mockLogger)
		loggerContext = new(mockLoggerContext)
	)
	loggerContext.Mock.On("LogAssignment", mock.Anything, mock.Anything).Return()
	loggerContext.Mock.On("LogBanditAction", mock.Anything, mock.Anything).Return()

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

	client := newEppoClient(
		newConfigurationStoreWithConfig(configuration{flags: flags, bandits: bandits}),
		nil,
		nil,
		logger,
		loggerContext,
		applicationLogger,
	)
	actions := map[string]ContextAttributes{
		"action1": {},
	}
	client.GetBanditAction("testFlag", "subject", ContextAttributes{}, actions, "bandit")

	event := loggerContext.Calls[0].Arguments[1].(BanditEvent)
	assert.Equal(t, "testFlag", event.FlagKey)
	assert.Equal(t, "bandit", event.BanditKey)
	assert.Equal(t, "subject", event.Subject)
	assert.Equal(t, "action1", event.Action)
}

func Test_GetStringAssignmentHandlesLoggingPanic(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Panic("logging panic")

	config := configResponse{Flags: map[string]*flagConfiguration{
		"experiment-key-1": &flagConfiguration{
			Key:           "experiment-key-1",
			Enabled:       true,
			TotalShards:   10000,
			VariationType: stringVariation,
			Variations: map[string]variation{
				"control": variation{
					Key:   "control",
					Value: []byte("\"control\""),
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

	client := newEppoClient(newConfigurationStoreWithConfig(configuration{flags: config}), nil, nil, mockLogger, nil, applicationLogger)

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

	client := newEppoClient(newConfigurationStoreWithConfig(configuration{flags: flags, bandits: bandits}), nil, nil, logger, nil, applicationLogger)
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

	client := newEppoClient(newConfigurationStoreWithConfig(configuration{flags: flags, bandits: bandits}), nil, nil, logger, nil, applicationLogger)
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

func Test_Initialized_timeout(t *testing.T) {
	var mockLogger = new(mockLogger)
	client := newEppoClient(newConfigurationStore(), nil, nil, mockLogger, nil, applicationLogger)

	timedOut := false
	select {
	case <-client.Initialized():
		timedOut = false
	case <-time.After(1 * time.Millisecond):
		timedOut = true
	}

	assert.True(t, timedOut)
}

func Test_Initialized_success(t *testing.T) {
	var mockLogger = new(mockLogger)
	configurationStore := newConfigurationStore()
	client := newEppoClient(configurationStore, nil, nil, mockLogger, nil, applicationLogger)

	go func() {
		<-time.After(1 * time.Microsecond)
		configurationStore.setConfiguration(configuration{})
	}()

	timedOut := false
	select {
	case <-client.Initialized():
		timedOut = false
	case <-time.After(1 * time.Millisecond):
		timedOut = true
	}

	assert.False(t, timedOut)
}
