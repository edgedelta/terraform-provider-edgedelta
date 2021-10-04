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
