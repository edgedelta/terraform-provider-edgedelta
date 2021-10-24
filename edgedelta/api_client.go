package edgedelta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	OrgID      string
	APIBaseURL string
	apiSecret  string
	cl         *http.Client
}

func (cli *APIClient) initializeHTTPClient() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 1
	t.MaxConnsPerHost = 1
	t.MaxIdleConnsPerHost = 1

	cli.cl = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

func (cli *APIClient) doRequest(entityName string, entityID string, method string, checkOKResp bool, bodyObj interface{}) ([]byte, int, error) {
	var baseURL *url.URL
	var err error

	if entityID == "" {
		baseURL, err = url.Parse(fmt.Sprintf("%s/v1/orgs/%s/%s", cli.APIBaseURL, cli.OrgID, entityName))
	} else {
		baseURL, err = url.Parse(fmt.Sprintf("%s/v1/orgs/%s/%s/%s", cli.APIBaseURL, cli.OrgID, entityName, entityID))
	}
	if err != nil {
		return nil, 0, fmt.Errorf("url parsing error: %v (base url was '%s')", err, cli.APIBaseURL)
	}
	var d io.Reader = nil
	if bodyObj != nil {
		db, err := json.Marshal(bodyObj)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal the config object: %v", err)
		}
		d = bytes.NewBuffer(db)
	}
	req, err := http.NewRequest(method, baseURL.String(), d)
	if err != nil {
		return nil, 0, fmt.Errorf("http request wrapper error: %v (base url was '%s')", err, cli.APIBaseURL)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-ED-API-Token", cli.apiSecret)
	resp, err := cli.cl.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to do '%s %s'. error: %v", req.Method, req.URL.RequestURI(), err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body from '%s'. err: %v", req.URL.RequestURI(), err)
	}
	if checkOKResp && (200 > resp.StatusCode || resp.StatusCode > 299) {
		return nil, 0, fmt.Errorf("got non OK http status from: %s, status: %v, response: %q", req.URL.RequestURI(), resp.StatusCode, string(body))
	}
	return body, resp.StatusCode, nil
}

func (cli *APIClient) GetConfigWithID(configID string) (*GetConfigResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", configID, http.MethodGet, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData GetConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) CreateConfig(configObject Config) (*CreateConfigResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", "", http.MethodPost, true, configObject)
	if err != nil {
		return nil, err
	}
	var responseData CreateConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) UpdateConfigWithID(configID string, configObject Config) (*UpdateConfigResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", configID, http.MethodPut, true, configObject)
	if err != nil {
		return nil, err
	}
	var responseData UpdateConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) GetMonitorWithID(monitorID string) (*GetMonitorResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("alert_definitions", monitorID, http.MethodGet, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData GetMonitorResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) CreateMonitor(monitor Monitor) (*CreateMonitorResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("alert_definitions", "", http.MethodPost, true, monitor)
	if err != nil {
		return nil, err
	}
	var responseData CreateMonitorResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) UpdateMonitorWithID(monitorID string, monitor Monitor) (*UpdateMonitorResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("alert_definitions", monitorID, http.MethodPut, true, monitor)
	if err != nil {
		return nil, err
	}
	var responseData UpdateMonitorResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) DeleteMonitorWithID(monitorID string) error {
	cli.initializeHTTPClient()
	_, _, err := cli.doRequest("alert_definitions", monitorID, http.MethodDelete, true, nil)
	if err != nil {
		return err
	}
	return nil
}
