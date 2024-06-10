# Eppo Go SDK

[![Test and lint SDK](https://github.com/Eppo-exp/java-server-sdk/actions/workflows/lint-test-sdk.yml/badge.svg)](https://github.com/Eppo-exp/java-server-sdk/actions/workflows/lint-test-sdk.yml)

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
    github.com/Eppo-exp/golang-sdk/v3
)
```

Or you can install the SDK from the command line with:

```
go get github.com/Eppo-exp/golang-sdk/v3
```

## Quick start

Begin by initializing a singleton instance of Eppo's client. Once initialized, the client can be used to make assignments anywhere in your app.

#### Initialize once

```go
import (
    "github.com/Eppo-exp/golang-sdk/v3/eppoclient"
)

var eppoClient = &eppoclient.EppoClient{}

func main() {
    assignmentLogger := NewExampleAssignmentLogger()

    eppoClient = eppoclient.InitClient(eppoclient.Config{
        SdkKey:           "<your_sdk_key>",
        AssignmentLogger: assignmentLogger,
    })
}
```


#### Assign anywhere

```go
import (
    "github.com/Eppo-exp/golang-sdk/v3/eppoclient"
)

var eppoClient = &eppoclient.EppoClient{}

variation := eppoClient.GetStringAssignment(
   'new-user-onboarding', 
   user.id, 
   user.attributes, 
   'control'
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
	subjectAttributes map[string]interface{}, 
	defaultValue string
) bool
  ```

## Assignment logger 

If you are using the Eppo SDK for experiment assignment (i.e randomization), pass in a callback logging function to the `InitClient` function on SDK initialization. The SDK invokes the callback to capture assignment data whenever a variation is assigned.

The code below illustrates an example implementation of a logging callback using Segment. You could also use your own logging system, the only requirement is that the SDK receives a `LogAssignment` function. Here we define an implementation of the Eppo `IAssignmentLogger` interface containing a single function named `LogAssignment`:


```go
import (
  "github.com/Eppo-exp/golang-sdk/v2/eppoclient"
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

## Philosophy

Eppo's SDKs are built for simplicity, speed and reliability. Flag configurations are compressed and distributed over a global CDN (Fastly), typically reaching your servers in under 15ms. Server SDKs continue polling Eppo’s API at 30-second intervals. Configurations are then cached locally, ensuring that each assignment is made instantly. Evaluation logic within each SDK consists of a few lines of simple numeric and string comparisons. The typed functions listed above are all developers need to understand, abstracting away the complexity of the Eppo's underlying (and expanding) feature set.
