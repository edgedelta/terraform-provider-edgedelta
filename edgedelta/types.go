package edgedelta

type EDConfig struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	ID          string `json:"id"`
	OrgID       string `json:"orgID"`
	Tag         string `json:"tag"`
}

type EDGetAllConfigsResponse []EDConfig
type EDGetConfigResponse EDConfig
type EDCreateConfigRequest EDConfig
