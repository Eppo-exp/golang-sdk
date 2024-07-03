package eppoclient

import (
	"fmt"
	"io"
	"time"

	"net/http"
)

const REQUEST_TIMEOUT_SECONDS = time.Duration(10 * time.Second)

type HttpClientInterface interface {
	get(resource string) (HttpResponse, error)
}

type httpClient struct {
	baseUrl   string
	sdkParams SDKParams
	client    *http.Client
}

type SDKParams struct {
	sdkKey     string
	sdkName    string
	sdkVersion string
}

type HttpResponse struct {
	Body string
	ETag string
}

func newHttpClient(baseUrl string, client *http.Client, sdkParams SDKParams) HttpClientInterface {
	var hc = &httpClient{
		baseUrl:   baseUrl,
		sdkParams: sdkParams,
		client:    client,
	}
	return hc
}

func (hc *httpClient) get(resource string) (HttpResponse, error) {
	url := hc.baseUrl + resource

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return HttpResponse{}, err // Return empty strings and the error
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
		return HttpResponse{}, err // Return empty strings and the error
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode == 401 {
		return HttpResponse{}, fmt.Errorf("unauthorized access") // Return an error indicating unauthorized access
	}

	if resp.StatusCode >= 500 {
		return HttpResponse{}, fmt.Errorf("server error: %d", resp.StatusCode) // Handle server errors (status code > 500)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return HttpResponse{}, fmt.Errorf("server error: unreadable body") // Return empty strings and the error
	}

	return HttpResponse{
		Body: string(b),
		ETag: resp.Header.Get("ETag"), // Capture the ETag header
	}, nil
}
