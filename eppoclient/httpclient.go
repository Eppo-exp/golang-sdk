package eppoclient

import (
	"fmt"
	"io"
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
	apiKey     string
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

func (hc *httpClient) get(resource string) (string, error) {
	url := hc.baseUrl + resource

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err // Return an empty string and the error
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	q := req.URL.Query()
	q.Add("apiKey", hc.sdkParams.apiKey)
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
		return "", err // Return an empty string and the error
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode == 401 {
		hc.isUnauthorized = true
		return "", fmt.Errorf("unauthorized access") // Return an error indicating unauthorized access
	}

	if resp.StatusCode >= 500 {
		return "", fmt.Errorf("server error: %d", resp.StatusCode) // Handle server errors (status code > 500)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("server error: unreadable body") // Return an empty string and the error
	}
	return string(b), nil // Return the response body as a string and nil for the error
}
