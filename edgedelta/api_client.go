package edgedelta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func validateUUID(val string) bool {
	_, err := uuid.Parse(val)
	return err == nil
}

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
		Timeout:   60 * time.Second,
		Transport: t,
	}
}

func (cli *APIClient) doRequest(entityName string, entityID string, method string, checkOKResp bool, checkNilBody bool, bodyObj interface{}) ([]byte, int, error) {
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
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body from '%s'. err: %v", req.URL.RequestURI(), err)
	}
	if checkOKResp && (200 > resp.StatusCode || resp.StatusCode > 299) {
		return nil, 0, fmt.Errorf("got non OK http status from: %s %s, status: %v, response: %q", req.Method, req.URL.RequestURI(), resp.StatusCode, string(body))
	}
	if checkNilBody && (strings.TrimSpace(string(body)) == "null") {
		return nil, 0, fmt.Errorf("API returned null response body from: %s, status: %v, response: %q", req.URL.RequestURI(), resp.StatusCode, string(body))

	}
	return body, resp.StatusCode, nil
}

func (cli *APIClient) GetConfigWithID(configID string) (*GetConfigResponse, error) {
	if ok := validateUUID(configID); !ok {
		return nil, fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", configID, http.MethodGet, true, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData GetConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) GetAllConfigs() ([]*Config, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", "", http.MethodGet, true, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData []*Config
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return responseData, nil
}

func (cli *APIClient) CreateConfig(configObject Config) (*CreateConfigResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", "", http.MethodPost, true, true, configObject)
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
	if ok := validateUUID(configID); !ok {
		return nil, fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("confs", configID, http.MethodPut, true, true, configObject)
	if err != nil {
		return nil, err
	}
	var responseData UpdateConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) SaveConfig(configID string, saveReq SaveRequest) (*SaveConfigResponse, error) {
	if ok := validateUUID(configID); !ok {
		return nil, fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("pipelines", fmt.Sprintf("%s/save", configID), http.MethodPost, true, true, saveReq)
	if err != nil {
		return nil, err
	}
	var responseData SaveConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) DeployConfig(configID string, version int64) (*DeployConfigResponse, error) {
	if ok := validateUUID(configID); !ok {
		return nil, fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("pipelines", fmt.Sprintf("%s/deploy/%d", configID, version), http.MethodPost, true, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData DeployConfigResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

func (cli *APIClient) GetLatestConfigHistoryVersion(configID string) (int64, error) {
	if ok := validateUUID(configID); !ok {
		return 0, fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("pipelines", fmt.Sprintf("%s/history", configID), http.MethodGet, true, true, nil)
	if err != nil {
		return 0, err
	}
	var histories []ConfigHistory
	if err := json.Unmarshal(b, &histories); err != nil {
		return 0, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	if len(histories) == 0 {
		return 0, fmt.Errorf("no config history found for config ID: '%s'", configID)
	}
	// The backend returns histories sorted by timestamp descending, so the first entry is the latest
	return histories[0].Timestamp, nil
}

func (cli *APIClient) DeleteConfigWithID(configID string) error {
	if ok := validateUUID(configID); !ok {
		return fmt.Errorf("failed to validate the config ID: '%s'", configID)
	}
	cli.initializeHTTPClient()
	_, _, err := cli.doRequest("confs", configID, http.MethodDelete, true, true, nil)
	if err != nil {
		return err
	}
	return nil
}

// Dashboard API methods

// GetDashboard retrieves a single dashboard by ID
func (cli *APIClient) GetDashboard(dashboardID string) (*GetDashboardResponse, error) {
	if ok := validateUUID(dashboardID); !ok {
		return nil, fmt.Errorf("failed to validate the dashboard ID: '%s'", dashboardID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("dashboards", dashboardID, http.MethodGet, true, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData GetDashboardResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

// GetAllDashboards retrieves all dashboards for the organization (used for import)
func (cli *APIClient) GetAllDashboards() ([]*Dashboard, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("dashboards", "", http.MethodGet, true, true, nil)
	if err != nil {
		return nil, err
	}
	var responseData []*Dashboard
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return responseData, nil
}

// CreateDashboard creates a new dashboard
func (cli *APIClient) CreateDashboard(dashboard *Dashboard) (*CreateDashboardResponse, error) {
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("dashboards", "", http.MethodPost, true, true, dashboard)
	if err != nil {
		return nil, err
	}
	var responseData CreateDashboardResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

// UpdateDashboard updates an existing dashboard
func (cli *APIClient) UpdateDashboard(dashboardID string, dashboard *Dashboard) (*UpdateDashboardResponse, error) {
	if ok := validateUUID(dashboardID); !ok {
		return nil, fmt.Errorf("failed to validate the dashboard ID: '%s'", dashboardID)
	}
	cli.initializeHTTPClient()
	b, _, err := cli.doRequest("dashboards", dashboardID, http.MethodPut, true, true, dashboard)
	if err != nil {
		return nil, err
	}
	var responseData UpdateDashboardResponse
	if err := json.Unmarshal(b, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response body: %s", err)
	}
	return &responseData, nil
}

// DeleteDashboard deletes a dashboard by ID
func (cli *APIClient) DeleteDashboard(dashboardID string) error {
	if ok := validateUUID(dashboardID); !ok {
		return fmt.Errorf("failed to validate the dashboard ID: '%s'", dashboardID)
	}
	cli.initializeHTTPClient()
	_, _, err := cli.doRequest("dashboards", dashboardID, http.MethodDelete, true, false, nil)
	if err != nil {
		return err
	}
	return nil
}
