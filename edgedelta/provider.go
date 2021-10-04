package edgedelta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Required params
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique organization ID",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API secret",
				Sensitive:   true,
			},
			// Optional params
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.edgedelta.com",
				Description: "API base URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"edgedelta_config": resourceConfig(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

type ProviderMetadata struct {
	client ConfigAPIClient
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	return &ProviderMetadata{
		client: ConfigAPIClient{
			APIBaseURL: d.Get("api_endpoint").(string),
			OrgID:      d.Get("org_id").(string),
			apiSecret:  d.Get("api_secret").(string),
		},
	}, nil
}
