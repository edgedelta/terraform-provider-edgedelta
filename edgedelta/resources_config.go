package edgedelta

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Linux",
				Description: "Environment type (Kubernetes, Linux, Windows, MacOS, Docker, Helm)",
			},
			"fleet_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Edge",
				Description: "Fleet type (Edge, Cloud, Gateway)",
			},
			"fleet_subtype": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Fleet subtype (Edge, Coordinator, Gateway) - required when environment is Kubernetes",
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
					return nil, fmt.Errorf("could not determine the resource ID - possibly the ID was not set")
				}
				var confIDs []string
				if confID == "*" {
					confs, err := meta.client.GetAllConfigs()
					if err != nil {
						return nil, fmt.Errorf("could not get the configs from API: %s", err)
					}
					results := make([]*schema.ResourceData, 0, len(confs))
					for _, c := range confs {
						dd := resourceConfig().Data(nil)
						dd.SetId(c.ID)
						dd.Set("conf_id", c.ID)
						dd.Set("tag", c.Tag)
						dd.Set("org_id", c.OrgID)
						dd.Set("config_content", c.Content)
						results = append(results, dd)
					}
					return results, nil
				} else if strings.Contains(confID, ",") {
					confIDs = strings.Split(confID, ",")
				} else {
					confIDs = []string{confID}
				}
				results := make([]*schema.ResourceData, 0, len(confIDs))
				for _, id := range confIDs {
					dd := resourceConfig().Data(nil)
					resp, err := meta.client.GetConfigWithID(id)
					if err != nil {
						return nil, fmt.Errorf("could not get the resource data from API: %s (resource ID was: '%s')", err, id)
					}
					dd.SetId(resp.ID)
					dd.Set("conf_id", id)
					dd.Set("tag", resp.Tag)
					dd.Set("org_id", resp.OrgID)
					dd.Set("config_content", resp.Content)
					results = append(results, dd)
				}
				return results, nil
			},
		},
	}
}

// extractTagFromYAML attempts to extract the tag value from the YAML content
func extractTagFromYAML(content string) string {
	// Try to find tag in settings section (v3 format)
	re := regexp.MustCompile(`(?m)^\s*tag:\s*(.+)$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
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
	
	// Get the optional environment and fleet configuration
	environment := d.Get("environment").(string)
	fleetType := d.Get("fleet_type").(string)
	fleetSubtype := d.Get("fleet_subtype").(string)
	
	// Validate: fleet_subtype is required when environment is Kubernetes
	if environment == "Kubernetes" && fleetSubtype == "" {
		fleetSubtype = "Edge" // Default to Edge for Kubernetes
	}
	
	// Extract tag from YAML content
	tag := extractTagFromYAML(confData)
	
	confDataObj := Config{
		Content:      confData,
		Tag:          tag,
		Environment:  environment,
		FleetType:    fleetType,
		FleetSubtype: fleetSubtype,
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
		// Note: For updates, we don't send environment/fleet fields as they can't be changed
		updateObj := Config{
			Content: confData,
		}
		apiResp, err := meta.client.UpdateConfigWithID(confID, updateObj)
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
	confDataObj := Config{
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
