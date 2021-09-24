package edgedelta

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfig() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
		Schema: map[string]*schema.Schema{
			// Required params
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique organization ID",
			},
			"api_key_envvar": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API base URL",
			},
			"config_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configuration file data",
			},
			// Optional params
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.edgedelta.com",
				Description: "API base URL",
			},
			"conf_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique configuration ID",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Debug flag",
			},
		},
	}
}

func parseArgs(d *schema.ResourceData) (orgID string, confID string, apiEndpoint string, apiKey string, confData string, diags diag.Diagnostics) {
	orgIDRaw := d.Get("org_id")
	confIDRaw := d.Get("conf_id")
	apiEndpointRaw := d.Get("api_endpoint")
	apiKeyRaw := d.Get("api_key_envvar")
	configDataRaw := d.Get("config_content")

	if orgIDRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "org_id is nil",
		})
	} else {
		orgID = orgIDRaw.(string)
	}
	if confIDRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "conf_id is nil",
		})
	} else {
		confID = confIDRaw.(string)
	}
	if apiEndpointRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "api_endpoint is nil",
		})
	} else {
		apiEndpoint = apiEndpointRaw.(string)
	}
	if apiKeyRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "api_key_envvar is nil",
		})
	} else {
		apiKey = os.Getenv(apiKeyRaw.(string))
	}
	if configDataRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "conf_data is nil",
		})
	} else {
		confData = configDataRaw.(string)
	}

	return
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	orgID, confID, apiEndpoint, apiKey, confData, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	var confDataObj Config
	confDataObj = Config{
		Content: confData,
	}
	apiClient := ConfigAPIClient{
		OrgID:      orgID,
		APIBaseURL: apiEndpoint,
		apiKey:     apiKey,
	}
	if confID == "" {
		// Create a new config
		apiResp, err := apiClient.createConfig(confDataObj)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not create the config resource",
				Detail:   fmt.Sprintf("%s", err),
			})
			return diags
		}
		d.SetId(apiResp.ID)
		d.Set("conf_id", apiResp.ID)
		d.Set("org_id", apiResp.OrgID)
	} else {
		// First run of the terraform config, just update the existing ed-config
		apiResp, err := apiClient.updateConfigWithID(confID, confDataObj)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not update the config resource",
				Detail:   fmt.Sprintf("%s", err),
			})
			return diags
		}
		d.SetId(apiResp.ID)
		d.Set("conf_id", apiResp.ID)
		d.Set("org_id", apiResp.OrgID)
	}

	return diags
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	orgID, confID, apiEndpoint, apiKey, _, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	apiClient := ConfigAPIClient{
		OrgID:      orgID,
		APIBaseURL: apiEndpoint,
		apiKey:     apiKey,
	}
	apiResp, err := apiClient.getConfigWithID(confID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not get the resource data from API",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}
	d.SetId(apiResp.ID)
	d.Set("conf_id", apiResp.ID)
	d.Set("org_id", apiResp.OrgID)

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	orgID, confID, apiEndpoint, apiKey, confData, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	if confID == "" {
		// Just get the config id from the tf state
		d.Set("conf_id", d.Id())
		confID = d.Id()
	}
	var confDataObj Config
	confDataObj = Config{
		Content: confData,
	}
	apiClient := ConfigAPIClient{
		OrgID:      orgID,
		APIBaseURL: apiEndpoint,
		apiKey:     apiKey,
	}
	_, err := apiClient.updateConfigWithID(confID, confDataObj)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not update the config resource",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}

	return diags
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
