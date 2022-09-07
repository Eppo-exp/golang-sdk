package eppoclient

import (
	"io"
	"log"
	"time"

	"net/http"
)

const REQUEST_TIMEOUT_SECONDS = time.Duration(1 * time.Second)

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

func (hc *httpClient) get(resource string) string {
	url := hc.baseUrl + resource

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	q := req.URL.Query()
	q.Add("apiKey", hc.sdkParams.apiKey)
	q.Add("sdkName", hc.sdkParams.sdkName)
	q.Add("sdkVersion", hc.sdkParams.sdkVersion)
	req.URL.RawQuery = q.Encode()

	resp, err := hc.client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 401 {
		hc.isUnauthorized = true
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}
