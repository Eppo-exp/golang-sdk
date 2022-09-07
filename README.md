# Eppo SDK for Golang

Eppoclient is a client sdk for `eppo.cloud` randomization API.
It is used to retrieve the experiments data and put it to in-memory cache, and then get assignments information.

## Getting Started

Refer to our [SDK documentation](https://docs.geteppo.com/prerequisites/feature-flagging/randomization-sdk/) for how to install and use the SDK.

## Supported Go Versions
This version of the SDK is compatible with Go v1.18 and above.

## Example


```
	import (
		"github.com/Eppo-exp/golang-sdk/eppoclient"
	)

	var eppoClient = &eppoclient.EppoClient{}

	func main() {
		eppoClient = eppoclient.InitClient(eppoclient.Config{
			ApiKey:           "<your_api_key>",
			BaseUrl:          "<base_url>", // optional, default https://eppo.cloud/api
			AssignmentLogger: eppoclient.AssignmentLogger{},
		})
	}

	func someBLFunc() {
		assignment, _ := eppoClient.GetAssignment("subject-1", "experiment_5", sbjAttrs)

		if assigment == "control" {
			// do something
		}
	}
```
