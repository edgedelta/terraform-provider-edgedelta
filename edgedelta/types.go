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
