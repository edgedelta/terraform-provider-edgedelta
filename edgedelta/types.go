package edgedelta

type EnvironmentType string

const (
	KubernetesEnvironmentType EnvironmentType = "Kubernetes"
	HelmEnvironmentType       EnvironmentType = "Helm"
	DockerEnvironmentType     EnvironmentType = "Docker"
	MacOSEnvironmentType      EnvironmentType = "MacOS"
	LinuxEnvironmentType      EnvironmentType = "Linux"
	WindowsEnvironmentType    EnvironmentType = "Windows"
)

type FleetType string

const (
	EdgeFleetType  FleetType = "Edge"
	CloudFleetType FleetType = "Cloud"
)

type FleetSubtype string

const (
	EdgeFleetSubtype        FleetSubtype = "Edge"
	CoordinatorFleetSubtype FleetSubtype = "Coordinator"
	GatewayFleetSubtype     FleetSubtype = "Gateway"
)

type Config struct {
	Content      string          `json:"content"`
	Description  string          `json:"description"`
	ID           string          `json:"id"`
	OrgID        string          `json:"orgID"`
	Tag          string          `json:"tag"`
	Environment  EnvironmentType `json:"environment"`
	FleetType    FleetType       `json:"fleet_type"`
	FleetSubtype FleetSubtype    `json:"fleet_subtype,omitempty"`
	ClusterName  string          `json:"cluster_name,omitempty"`
}

type GetConfigResponse Config
type UpdateConfigResponse Config
type CreateConfigResponse Config

type SaveRequest struct {
	Content     *string `json:"content,omitempty"`
	Description string  `json:"description"`
}

type SaveConfigResponse struct {
	ID          string `json:"id"`
	Content     string `json:"content,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

type DeployConfigResponse struct {
	ID          string `json:"id"`
	Content     string `json:"content,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

type ConfigHistory struct {
	ConfigID  string `json:"config_id"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
	Status    string `json:"status"`
}

// Dashboard represents an EdgeDelta dashboard
type Dashboard struct {
	OrgID                   string                   `json:"org_id,omitempty"`
	DashboardID             string                   `json:"dashboard_id,omitempty"`
	DashboardName           string                   `json:"dashboard_name"`
	Description             string                   `json:"description,omitempty"`
	Tags                    []string                 `json:"tags,omitempty"`
	Creator                 string                   `json:"creator,omitempty"`
	Updater                 string                   `json:"updater,omitempty"`
	Created                 string                   `json:"created,omitempty"`
	Updated                 string                   `json:"updated,omitempty"`
	Definition              map[string]interface{}   `json:"definition,omitempty"`
	ResourceAccesses        []map[string]interface{} `json:"resource_accesses,omitempty"`
	SharingSecuritySettings map[string]interface{}   `json:"sharing_security_settings,omitempty"`
}

// Dashboard API response types
type GetDashboardResponse Dashboard
type CreateDashboardResponse Dashboard
type UpdateDashboardResponse Dashboard
