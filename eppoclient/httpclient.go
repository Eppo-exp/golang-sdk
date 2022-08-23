package eppoclient

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

const REQUEST_TIMEOUT_SECONDS = time.Duration(1 * time.Second)

type HttpClient struct {
	baseUrl        string
	sdkParams      SDKParams
	isUnauthorized bool
	client         resty.Client
}

type SDKParams struct {
	apiKey     string
	sdkName    string
	sdkVersion string
}

// todo move this to requestor
type Experiment struct {
	Name   string
	Latest string
}

type Experiments struct {
	Results []*Experiment
}

func NewHttpClient(baseUrl string, sdkParams SDKParams) *HttpClient {
	var hc = &HttpClient{}
	hc.baseUrl = baseUrl
	hc.sdkParams = sdkParams
	hc.isUnauthorized = false
	hc.client = *resty.New()

	hc.client.SetTimeout(REQUEST_TIMEOUT_SECONDS)
	return hc
}

func (hc *HttpClient) Get(resource string) string {
	url := hc.baseUrl + resource

	resp, err := hc.client.R().
		SetQueryParams(map[string]string{
			"apiKey":     hc.sdkParams.apiKey,
			"sdkName":    hc.sdkParams.sdkName,
			"sdkVersion": hc.sdkParams.sdkVersion,
		}).
		Get(url)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode() == 401 {
		hc.isUnauthorized = true
	}

	return resp.String()
}
