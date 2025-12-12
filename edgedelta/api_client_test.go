package edgedelta

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	apiKey      = flag.String("api_key", "", "API auth key")
	apiEndpoint = flag.String("api_endpoint", "", "API endpoint")
	confID      = flag.String("conf_id", "", "Unique configuration ID")
	orgID       = flag.String("org_id", "", "Unique organization ID")
	confPath    = flag.String("conf_path", "", "Path to the new config file")
)

// Test constants
const (
	testOrgID       = "test-org-123"
	testAPISecret   = "test-api-secret"
	testDashboardID = "550e8400-e29b-41d4-a716-446655440000"
	testConfigID    = "660e8400-e29b-41d4-a716-446655440001"
)

// Helper function to create a mock server
func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// Helper function to create an API client pointing to mock server
func newTestClient(serverURL string) *APIClient {
	return &APIClient{
		OrgID:      testOrgID,
		APIBaseURL: serverURL,
		apiSecret:  testAPISecret,
	}
}

// =============================================================================
// Dashboard API Unit Tests (Mock Server)
// =============================================================================

func TestGetDashboard(t *testing.T) {
	expectedDashboard := Dashboard{
		DashboardID:   testDashboardID,
		DashboardName: "Test Dashboard",
		Description:   "A test dashboard",
		Tags:          []string{"test", "unit"},
		Creator:       "user-123",
		Updater:       "user-123",
		Created:       "2024-01-01T00:00:00Z",
		Updated:       "2024-01-02T00:00:00Z",
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/dashboards/" + testDashboardID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.Header.Get("X-ED-API-Token") != testAPISecret {
			t.Errorf("expected API token header, got %s", r.Header.Get("X-ED-API-Token"))
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDashboard); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetDashboard(testDashboardID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DashboardID != expectedDashboard.DashboardID {
		t.Errorf("expected dashboard ID %s, got %s", expectedDashboard.DashboardID, result.DashboardID)
	}
	if result.DashboardName != expectedDashboard.DashboardName {
		t.Errorf("expected dashboard name %s, got %s", expectedDashboard.DashboardName, result.DashboardName)
	}
	if result.Description != expectedDashboard.Description {
		t.Errorf("expected description %s, got %s", expectedDashboard.Description, result.Description)
	}
}

func TestGetDashboard_InvalidID(t *testing.T) {
	client := &APIClient{
		OrgID:      testOrgID,
		APIBaseURL: "http://localhost",
		apiSecret:  testAPISecret,
	}

	_, err := client.GetDashboard("invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid UUID, got nil")
	}
	if !strings.Contains(err.Error(), "failed to validate the dashboard ID") {
		t.Errorf("expected UUID validation error, got: %v", err)
	}
}

func TestGetDashboard_NotFound(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "dashboard not found"}`))
	})
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetDashboard(testDashboardID)

	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

func TestGetAllDashboards(t *testing.T) {
	expectedDashboards := []*Dashboard{
		{
			DashboardID:   "550e8400-e29b-41d4-a716-446655440001",
			DashboardName: "Dashboard 1",
		},
		{
			DashboardID:   "550e8400-e29b-41d4-a716-446655440002",
			DashboardName: "Dashboard 2",
		},
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/dashboards"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDashboards); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetAllDashboards()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 dashboards, got %d", len(result))
	}
	if result[0].DashboardName != "Dashboard 1" {
		t.Errorf("expected 'Dashboard 1', got %s", result[0].DashboardName)
	}
}

func TestGetAllDashboards_Empty(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetAllDashboards()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 dashboards, got %d", len(result))
	}
}

func TestCreateDashboard(t *testing.T) {
	inputDashboard := &Dashboard{
		DashboardName: "New Dashboard",
		Description:   "Created via API",
		Tags:          []string{"new", "test"},
	}

	expectedResponse := Dashboard{
		DashboardID:   testDashboardID,
		DashboardName: "New Dashboard",
		Description:   "Created via API",
		Tags:          []string{"new", "test"},
		Creator:       "user-123",
		Created:       "2024-01-01T00:00:00Z",
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/dashboards"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify request body
		var receivedDashboard Dashboard
		if err := json.NewDecoder(r.Body).Decode(&receivedDashboard); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if receivedDashboard.DashboardName != inputDashboard.DashboardName {
			t.Errorf("expected name %s, got %s", inputDashboard.DashboardName, receivedDashboard.DashboardName)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateDashboard(inputDashboard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DashboardID != testDashboardID {
		t.Errorf("expected dashboard ID %s, got %s", testDashboardID, result.DashboardID)
	}
	if result.DashboardName != inputDashboard.DashboardName {
		t.Errorf("expected name %s, got %s", inputDashboard.DashboardName, result.DashboardName)
	}
}

func TestUpdateDashboard(t *testing.T) {
	inputDashboard := &Dashboard{
		DashboardName: "Updated Dashboard",
		Description:   "Updated description",
	}

	expectedResponse := Dashboard{
		DashboardID:   testDashboardID,
		DashboardName: "Updated Dashboard",
		Description:   "Updated description",
		Updater:       "user-456",
		Updated:       "2024-01-03T00:00:00Z",
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/dashboards/" + testDashboardID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateDashboard(testDashboardID, inputDashboard)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DashboardName != inputDashboard.DashboardName {
		t.Errorf("expected name %s, got %s", inputDashboard.DashboardName, result.DashboardName)
	}
	if result.Description != inputDashboard.Description {
		t.Errorf("expected description %s, got %s", inputDashboard.Description, result.Description)
	}
}

func TestUpdateDashboard_InvalidID(t *testing.T) {
	client := &APIClient{
		OrgID:      testOrgID,
		APIBaseURL: "http://localhost",
		apiSecret:  testAPISecret,
	}

	_, err := client.UpdateDashboard("invalid-uuid", &Dashboard{})
	if err == nil {
		t.Error("expected error for invalid UUID, got nil")
	}
}

func TestDeleteDashboard(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/dashboards/" + testDashboardID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteDashboard(testDashboardID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteDashboard_InvalidID(t *testing.T) {
	client := &APIClient{
		OrgID:      testOrgID,
		APIBaseURL: "http://localhost",
		apiSecret:  testAPISecret,
	}

	err := client.DeleteDashboard("invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid UUID, got nil")
	}
}

func TestDeleteDashboard_NotFound(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "dashboard not found"}`))
	})
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteDashboard(testDashboardID)

	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

// =============================================================================
// Config API Unit Tests (Mock Server)
// =============================================================================

func TestGetConfigWithID_Mock(t *testing.T) {
	expectedConfig := GetConfigResponse{
		ID:          testConfigID,
		Content:     "test: config",
		Description: "Test config",
		Tag:         "v1.0.0",
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		expectedPath := "/v1/orgs/" + testOrgID + "/confs/" + testConfigID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedConfig); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetConfigWithID(testConfigID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != expectedConfig.ID {
		t.Errorf("expected config ID %s, got %s", expectedConfig.ID, result.ID)
	}
}

func TestGetConfigWithID_InvalidID(t *testing.T) {
	client := &APIClient{
		OrgID:      testOrgID,
		APIBaseURL: "http://localhost",
		apiSecret:  testAPISecret,
	}

	_, err := client.GetConfigWithID("invalid-uuid")
	if err == nil {
		t.Error("expected error for invalid UUID, got nil")
	}
}

func TestGetAllConfigs_Mock(t *testing.T) {
	expectedConfigs := []*Config{
		{ID: "550e8400-e29b-41d4-a716-446655440001", Content: "config1"},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Content: "config2"},
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedConfigs); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetAllConfigs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 configs, got %d", len(result))
	}
}

func TestCreateConfig_Mock(t *testing.T) {
	inputConfig := Config{
		Content:     "new: config",
		Environment: KubernetesEnvironmentType,
	}

	expectedResponse := CreateConfigResponse{
		ID:      testConfigID,
		Content: "new: config",
		Tag:     "v1.0.0",
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateConfig(inputConfig)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != testConfigID {
		t.Errorf("expected config ID %s, got %s", testConfigID, result.ID)
	}
}

func TestDeleteConfigWithID_Mock(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})
	defer server.Close()

	client := newTestClient(server.URL)
	if err := client.DeleteConfigWithID(testConfigID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Integration Tests (Require Real API - Skip if no credentials)
// =============================================================================

func TestGetConfigWithID_Integration(t *testing.T) {
	if *orgID == "" || *apiKey == "" || *confID == "" || *apiEndpoint == "" {
		t.Skip("Skipping integration test: missing required flags (org_id, api_key, conf_id, api_endpoint)")
	}

	cli := APIClient{
		OrgID:      *orgID,
		apiSecret:  *apiKey,
		APIBaseURL: *apiEndpoint,
	}

	confObject, err := cli.GetConfigWithID(*confID)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", confObject)
}

func TestUpdateConfigWithID_Integration(t *testing.T) {
	if *orgID == "" || *confID == "" || *confPath == "" || *apiEndpoint == "" || *apiKey == "" {
		t.Skip("Skipping integration test: missing required flags")
	}

	cli := APIClient{
		OrgID:      *orgID,
		APIBaseURL: *apiEndpoint,
		apiSecret:  *apiKey,
	}

	confDataRaw, err := os.ReadFile(*confPath)
	if err != nil {
		t.Error(err)
		return
	}

	confData := Config{
		Content: string(confDataRaw[:]),
	}

	confObject, err := cli.UpdateConfigWithID(*confID, confData)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", confObject)
}

func TestCreateConfig_Integration(t *testing.T) {
	if *orgID == "" || *confPath == "" || *apiKey == "" || *apiEndpoint == "" {
		t.Skip("Skipping integration test: missing required flags")
	}

	cli := APIClient{
		OrgID:      *orgID,
		APIBaseURL: *apiEndpoint,
		apiSecret:  *apiKey,
	}

	confDataRaw, err := os.ReadFile(*confPath)
	if err != nil {
		t.Error(err)
		return
	}

	confData := Config{
		Content: string(confDataRaw[:]),
	}

	confObject, err := cli.CreateConfig(confData)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", confObject)
}
