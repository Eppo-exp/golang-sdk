# Eppo Go SDK

[![Test SDK](https://github.com/Eppo-exp/golang-sdk/actions/workflows/test.yml/badge.svg)](https://github.com/Eppo-exp/golang-sdk/actions/workflows/test.yml)

EppoClient is a client sdk for the `eppo.cloud` randomization API.
It is used to retrieve the experiments data and put it to in-memory store, and then get assignments information.

[Eppo](https://www.geteppo.com/) is a modular flagging and experimentation analysis tool. Eppo's Go SDK is built to make assignments in multi-user server side contexts, compatible with Go v1.19 and above. Before proceeding you'll need an Eppo account.

## Features

- Feature gates
- Kill switches
- Progressive rollouts
- A/B/n experiments
- Mutually exclusive experiments (Layers)
- Dynamic configuration

## Installation

In your `go.mod`, add the SDK package as a dependency:

```
require (
    github.com/Eppo-exp/golang-sdk/v5
)
```

Or you can install the SDK from the command line with:

```
go get github.com/Eppo-exp/golang-sdk/v5
```

## Quick start

Begin by initializing a singleton instance of Eppo's client. Once initialized, the client can be used to make assignments anywhere in your app.

#### Initialize once

```go
import (
    "github.com/Eppo-exp/golang-sdk/v5/eppoclient"
)

var eppoClient *eppoclient.EppoClient

func main() {
    assignmentLogger := NewExampleAssignmentLogger()

    eppoClient, err = eppoclient.InitClient(eppoclient.Config{
        SdkKey:           "<your_sdk_key>",
        AssignmentLogger: assignmentLogger,
    })
    if err != nil {
        log.Fatalf("Failed to initialize Eppo client: %v", err)
    }
}
```

#### Assign anywhere

```go
import (
    "github.com/Eppo-exp/golang-sdk/v5/eppoclient"
)

var eppoClient *eppoclient.EppoClient

variation := eppoClient.GetStringAssignment(
   "new-user-onboarding",
   user.id,
   user.attributes,
   "control"
);
```

## Assignment functions

Every Eppo flag has a return type that is set once on creation in the dashboard. Once a flag is created, assignments in code should be made using the corresponding typed function:

```go
GetBooleanAssignment(...)
GetNumericAssignment(...)
GetIntegerAssignment(...)
GetStringAssignment(...)
GetJSONAssignment(...)
```

Each function has the same signature, but returns the type in the function name. For booleans use `getBooleanAssignment`, which has the following signature:

```go
func getBooleanAssignment(
	flagKey string,
	subjectKey string,
	subjectAttributes SubjectAttributes,
	defaultValue string
) bool
  ```

## Assignment logger

If you are using the Eppo SDK for experiment assignment (i.e randomization), pass in a callback logging function to the `InitClient` function on SDK initialization. The SDK invokes the callback to capture assignment data whenever a variation is assigned.

The code below illustrates an example implementation of a logging callback using Segment. You could also use your own logging system, the only requirement is that the SDK receives a `LogAssignment` function. Here we define an implementation of the Eppo `IAssignmentLogger` interface containing a single function named `LogAssignment`:


```go
import (
  "github.com/Eppo-exp/golang-sdk/v5/eppoclient"
  "gopkg.in/segmentio/analytics-go.v3"
)

func main() {
  // Connect to Segment (or your own event-tracking system)
  client := analytics.New("YOUR_WRITE_KEY")
  defer client.Close()

  type ExampleAssignmentLogger struct {
  }

  func NewExampleAssignmentLogger() *ExampleAssignmentLogger {
    return &ExampleAssignmentLogger{}
  }

  func (al *ExampleAssignmentLogger) LogAssignment(event eppoclient.AssignmentEvent) {
    client.Enqueue(analytics.Track{
      UserId: event.Subject,
      Event:  "Eppo Randomization Event",
      Properties: event
    })
  }
}
```

## Provide a custom logger

If you want to provide a logging implementation to the SDK to capture errors and other application logs, you can do so by passing in an implementation of the `ApplicationLogger` interface to the `InitClient` function.

You can use the `eppoclient.ScrubbingLogger` to scrub PII from the logs.

```go
import (
    "github.com/Eppo-exp/golang-sdk/v5/eppoclient"
    "github.com/sirupsen/logrus"
)

type LogrusApplicationLogger struct {
    logger *logrus.Logger
    logLevel logrus.Level
}

func NewLogrusApplicationLogger(logLevel logrus.Level) *LogrusApplicationLogger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    return &LogrusApplicationLogger{logger: logger, logLevel: logLevel}
}

func (l *LogrusApplicationLogger) Debug(args ...interface{}) {
    if l.logLevel <= logrus.DebugLevel {
        l.logger.Debug(args...)
    }
}

func (l *LogrusApplicationLogger) Info(args ...interface{}) {
    if l.logLevel <= logrus.InfoLevel {
        l.logger.Info(args...)
    }
}

func (l *LogrusApplicationLogger) Warn(args ...interface{}) {
    if l.logLevel <= logrus.WarnLevel {
        l.logger.Warn(args...)
    }
}

func (l *LogrusApplicationLogger) Error(args ...interface{}) {
    if l.logLevel <= logrus.ErrorLevel {
        l.logger.Error(args...)
    }
}

func main() {
  // Initialize a custom logger; example using logrus
  // Set log level to Info
  applicationLogger := NewLogrusApplicationLogger(logrus.InfoLevel)
  scrubbingLogger := eppoclient.NewScrubbingLogger(applicationLogger)

  // Initialize the Eppo client
  eppoClient, _ = eppoclient.InitClient(eppoclient.Config{
      SdkKey:            "<your_sdk_key>",
      AssignmentLogger:  assignmentLogger,
      ApplicationLogger: scrubbingLogger
  })
}
```

### De-duplication of assignments

The SDK may see many duplicate assignments in a short period of time, and if you
have configured a logging function, they will be transmitted to your downstream
event store. This increases the cost of storage as well as warehouse costs during experiment analysis.

To mitigate this, an in-memory assignment cache is optionally available with expiration based on the least recently accessed time.

It can be configured with a maximum size to fit your desired memory allocation.

```go
import (
  "github.com/Eppo-exp/golang-sdk/v5/eppoclient"
)

var eppoClient *eppoclient.EppoClient

func main() {
  assignmentLogger := NewExampleAssignmentLogger()

  eppoClient, _ = eppoclient.InitClient(eppoclient.Config{
    ApiKey:           "<your_sdk_key>",
    // 10000 is the maximum number of assignments to cache
    // Depending on the length of your flag and subject keys, taking a median
    // length of 32 characters, each assignment cache entry uses approximately 112 bytes.
    // Use this calculation to determine the maximum number of assignments to cache
    // for the memory you wish to allocate.
    AssignmentLogger: eppoclient.NewLruAssignmentLogger(assignmentLogger, 10000),
  })
}
```

## Philosophy

Eppo's SDKs are built for simplicity, speed and reliability. Flag configurations are compressed and distributed over a global CDN (Fastly), typically reaching your servers in under 15ms. Server SDKs continue polling Eppo’s API at 10-second intervals. Configurations are then cached locally, ensuring that each assignment is made instantly. Evaluation logic within each SDK consists of a few lines of simple numeric and string comparisons. The typed functions listed above are all developers need to understand, abstracting away the complexity of the Eppo's underlying (and expanding) feature set.
