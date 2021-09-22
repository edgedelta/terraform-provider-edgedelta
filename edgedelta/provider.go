package edgedelta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"edgedelta_config": resourceConfig(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
