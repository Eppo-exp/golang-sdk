/*
	EppoClient is a client sdk for the `eppo.cloud` randomization API.
	It is used to retrieve the experiments data and put it to in-memory store, and then get assignments information.

	Usage:

	import (
		"github.com/Eppo-exp/golang-sdk/v2/eppoclient"
	)

	var eppoClient = &eppoclient.EppoClient{}

	func main() {
		eppoClient = eppoclient.InitClient(eppoclient.Config{
			ApiKey:           "<your_api_key>",
			AssignmentLogger: eppoclient.AssignmentLogger{},
		})
	}

	func apiEndpoint() {
		assignment, _ := eppoClient.GetStringAssignment("subject-1", "experiment_5", sbjAttrs)

		if assignment == "control" {
			// do something
		}
	}
*/

package eppoclient
