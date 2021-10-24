package edgedelta

type Config struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	ID          string `json:"id"`
	OrgID       string `json:"orgID"`
	Tag         string `json:"tag"`
}

type GetConfigResponse Config
type UpdateConfigResponse Config
type CreateConfigResponse Config

type Monitor struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
	ID      string `json:"id"`
	OrgID   string `json:"orgID"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
	Creator string `json:"creator"`
}

type GetMonitorResponse Monitor
type UpdateMonitorResponse Monitor
type CreateMonitorResponse Monitor
