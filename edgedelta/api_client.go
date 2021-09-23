package edgedelta

import (
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
	apiKey     string
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
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", cli.apiKey)

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
