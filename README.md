# Eppo SDK for Golang

EppoClient is a client sdk for the `eppo.cloud` randomization API.
It is used to retrieve the experiments data and put it to in-memory cache, and then get assignments information.

## Getting Started

Refer to our [SDK documentation](https://docs.geteppo.com/feature-flags/sdks/server-sdks/go) for how to install and use the SDK.

## Supported Go Versions
This version of the SDK is compatible with Go v1.19 and above.

## Example


```
	import (
		"github.com/Eppo-exp/golang-sdk/v2/eppoclient"
	)

	var eppoClient = &eppoclient.EppoClient{}

	func main() {
		eppoClient, stopFn, err := eppoclient.NewClient(eppoclient.Config{
			ApiKey:           "<your_api_key>",
			AssignmentLogger: eppoclient.AssignmentLogger{},
		})
        if err != nil {
            fmt.Printf("init eppo client: %v", err)
            return
        }
        defer stopFn()

        // call the client
	}

	func someBLFunc() {
		assignment, _ := eppoClient.GetStringAssignment("subject-1", "experiment_5", sbjAttrs)

		if assignment == "control" {
			// do something
		}
	}
```
