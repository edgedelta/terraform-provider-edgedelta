package edgedelta

type Config struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	ID          string `json:"id"`
	OrgID       string `json:"orgID"`
	Tag         string `json:"tag"`
}

type GetAllConfigsResponse []Config
type GetConfigResponse Config
type CreateConfigRequest Config
