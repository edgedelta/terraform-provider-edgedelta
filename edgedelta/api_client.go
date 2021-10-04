package edgedelta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ConfigAPIClient struct {
	OrgID      string
	APIBaseURL string
	apiSecret  string
	cl         *http.Client
}

func (cli *ConfigAPIClient) initializeHTTPClient() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 1
	t.MaxConnsPerHost = 1
	t.MaxIdleConnsPerHost = 1

	cli.cl = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

func (cli *ConfigAPIClient) getConfigWithID(configID string) (*GetConfigResponse, error) {
	cli.initializeHTTPClient()

	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/orgs/%s/confs/%s", cli.APIBaseURL, cli.OrgID, configID))
	if err != nil {
		return nil, fmt.Errorf("url parsing error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http request wrapper error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-ED-API-Token", cli.apiSecret)

	resp, err := cli.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET %s. error: %v", req.URL.RequestURI(), err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from: %s. err: %v", req.URL.RequestURI(), err)
	}

	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return nil, fmt.Errorf("got non OK http status from: %s, status: %v, response: %q", req.URL.RequestURI(), resp.StatusCode, string(b))
	}

	var responseData GetConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}

	return &responseData, nil
}

func (cli *ConfigAPIClient) createConfig(configObject Config) (*CreateConfigResponse, error) {
	cli.initializeHTTPClient()

	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/orgs/%s/confs", cli.APIBaseURL, cli.OrgID))
	if err != nil {
		return nil, fmt.Errorf("url parsing error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	d, err := json.Marshal(configObject)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the config object: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewBuffer(d))
	if err != nil {
		return nil, fmt.Errorf("http request wrapper error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-ED-API-Token", cli.apiSecret)

	resp, err := cli.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do POST %s. error: %v", req.URL.RequestURI(), err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from: %s. err: %v", req.URL.RequestURI(), err)
	}

	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return nil, fmt.Errorf("got non OK http status from: %s, status: %v, response: %q", req.URL.RequestURI(), resp.StatusCode, string(b))
	}

	var responseData CreateConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}

	return &responseData, nil
}

func (cli *ConfigAPIClient) updateConfigWithID(configID string, configObject Config) (*UpdateConfigResponse, error) {
	cli.initializeHTTPClient()

	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/orgs/%s/confs/%s", cli.APIBaseURL, cli.OrgID, configID))
	if err != nil {
		return nil, fmt.Errorf("url parsing error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	d, err := json.Marshal(configObject)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the config object: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewBuffer(d))
	if err != nil {
		return nil, fmt.Errorf("http request wrapper error: %v (base url was '%s')", err, cli.APIBaseURL)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-ED-API-Token", cli.apiSecret)

	resp, err := cli.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do PUT %s. error: %v", req.URL.RequestURI(), err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from: %s. err: %v", req.URL.RequestURI(), err)
	}

	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return nil, fmt.Errorf("got non OK http status from: %s, status: %v, response: %q", req.URL.RequestURI(), resp.StatusCode, string(b))
	}

	var responseData UpdateConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}

	return &responseData, nil
}
