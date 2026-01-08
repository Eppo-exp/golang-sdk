package eppoclient

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"net/http"
)

const REQUEST_TIMEOUT_SECONDS = time.Duration(10 * time.Second)

type httpClient struct {
	baseUrl        string
	sdkParams      SDKParams
	isUnauthorized bool
	client         *http.Client
}

type SDKParams struct {
	sdkKey     string
	sdkName    string
	sdkVersion string
}

func newHttpClient(baseUrl string, client *http.Client, sdkParams SDKParams) *httpClient {
	var hc = &httpClient{
		baseUrl:        baseUrl,
		sdkParams:      sdkParams,
		isUnauthorized: false,
		client:         client,
	}
	return hc
}

// scrubError removes sensitive information (like apiKey or sdkKey) from error messages
// to prevent exposure of these keys when errors are returned or logged.
func scrubError(err error) error {
	if err == nil {
		return nil
	}
	re := regexp.MustCompile(`(apiKey|sdkKey)=[^&\s"]*`)
	scrubbedMsg := re.ReplaceAllString(err.Error(), "$1=XXXXXX")
	return fmt.Errorf("%s", scrubbedMsg)
}

func (hc *httpClient) get(resource string) ([]byte, error) {
	url := hc.baseUrl + resource

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	q := req.URL.Query()
	// todo: migrate to bearer token authorization header
	q.Add("apiKey", hc.sdkParams.sdkKey) // origin server uses apiKey
	q.Add("sdkName", hc.sdkParams.sdkName)
	q.Add("sdkVersion", hc.sdkParams.sdkVersion)
	req.URL.RawQuery = q.Encode()

	resp, err := hc.client.Do(req)
	if err != nil {
		// from https://golang.org/pkg/net/http/#Client.Do
		//
		// An error is returned if caused by client policy (such as
		// CheckRedirect), or failure to speak HTTP (such as a network
		// connectivity problem). A non-2xx status code doesn't cause an
		// error.
		//
		// We should almost never expect to see this condition be executed.
		// Scrub the error to prevent SDK key exposure in error messages.
		return nil, scrubError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		hc.isUnauthorized = true
		return nil, fmt.Errorf("unauthorized access")
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("server error: unreadable body")
	}
	return b, nil
}
