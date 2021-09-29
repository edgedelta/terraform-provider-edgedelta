package edgedelta

import (
	"context"
	"fmt"

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
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API secret",
				Sensitive:   true,
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

func parseArgs(d *schema.ResourceData) (orgID string, confID string, apiEndpoint string, apiSecret string, confData string, diags diag.Diagnostics) {
	orgIDRaw := d.Get("org_id")
	confIDRaw := d.Get("conf_id")
	apiEndpointRaw := d.Get("api_endpoint")
	apiSecretRaw := d.Get("api_secret")
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
	if apiSecretRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "api_secret is nil",
		})
	} else {
		apiSecret = apiSecretRaw.(string)
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
	orgID, confID, apiEndpoint, apiSecret, confData, diags := parseArgs(d)
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
		apiSecret:  apiSecret,
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
		d.Set("conf_id", confID)
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
	orgID, confID, apiEndpoint, apiSecret, _, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	apiClient := ConfigAPIClient{
		OrgID:      orgID,
		APIBaseURL: apiEndpoint,
		apiSecret:  apiSecret,
	}
	activeConfID := confID
	if activeConfID == "" {
		if d.Id() == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot determine config ID",
				Detail:   "Possibly tried to read the resource with conf_id=nil and id=nil",
			})
			return diags
		}

		activeConfID = d.Id()
	}
	apiResp, err := apiClient.getConfigWithID(activeConfID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not get the resource data from API",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}
	d.SetId(apiResp.ID)
	d.Set("conf_id", confID)
	d.Set("org_id", apiResp.OrgID)

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	orgID, confID, apiEndpoint, apiSecret, confData, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	if confID == "" {
		// Just get the config id from the tf state
		confID = d.Id()
	}
	var confDataObj Config
	confDataObj = Config{
		Content: confData,
	}
	apiClient := ConfigAPIClient{
		OrgID:      orgID,
		APIBaseURL: apiEndpoint,
		apiSecret:  apiSecret,
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
