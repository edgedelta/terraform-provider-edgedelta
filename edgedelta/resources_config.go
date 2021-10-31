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
			"config_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configuration file data",
			},
			// Optional params
			"conf_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique configuration ID",
			},
			// Computed
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				meta := m.(*ProviderMetadata)
				confID := d.Id()
				if confID == "" { // confID DNE
					return nil, fmt.Errorf("Could not determine the resource ID - possibly the ID was not set")
				}
				resp, err := meta.client.GetConfigWithID(confID)
				if err != nil {
					return nil, fmt.Errorf("Could not get the resource data from API: %s (resource ID was: '%s')", err, confID)
				}
				d.SetId(resp.ID)
				d.Set("conf_id", confID)
				d.Set("org_id", resp.OrgID)
				d.Set("tag", resp.Tag)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func parseArgs(d *schema.ResourceData) (confID string, confData string, diags diag.Diagnostics) {
	confIDRaw := d.Get("conf_id")
	configDataRaw := d.Get("config_content")
	if confIDRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "conf_id is nil",
		})
	} else {
		confID = confIDRaw.(string)
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
	meta := m.(*ProviderMetadata)
	confID, confData, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
	}
	var confDataObj Config
	confDataObj = Config{
		Content: confData,
	}
	if confID == "" {
		// Create a new config
		apiResp, err := meta.client.CreateConfig(confDataObj)
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
		d.Set("tag", apiResp.Tag)

	} else {
		// First run of the terraform config, just update the existing ed-config
		apiResp, err := meta.client.UpdateConfigWithID(confID, confDataObj)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not update the config resource (create=>update)",
				Detail:   fmt.Sprintf("%s", err),
			})
			return diags
		}
		d.SetId(apiResp.ID)
		d.Set("conf_id", apiResp.ID)
		d.Set("org_id", apiResp.OrgID)
		d.Set("tag", apiResp.Tag)
	}

	return diags
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	confID, _, diags := parseArgs(d)
	if len(diags) > 0 {
		return diags
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
	apiResp, err := meta.client.GetConfigWithID(activeConfID)
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
	d.Set("tag", apiResp.Tag)

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	confID, confData, diags := parseArgs(d)
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
	_, err := meta.client.UpdateConfigWithID(confID, confDataObj)
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
	return diag.Diagnostics{}
}
