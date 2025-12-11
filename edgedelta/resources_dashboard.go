package edgedelta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		CreateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		UpdateContext: resourceDashboardUpdate,
		DeleteContext: resourceDashboardDelete,
		Description:   "Manages an EdgeDelta dashboard resource.",
		Schema: map[string]*schema.Schema{
			// Required
			"dashboard_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the dashboard.",
			},

			// Optional
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the dashboard.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Searchable tags for the dashboard.",
			},
			"definition": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentJSON,
				ValidateFunc:     validateJSON,
				Description:      "Dashboard definition as a JSON string. Use file() or jsonencode() to provide the value.",
			},

			// Computed
			"dashboard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the dashboard.",
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID who created the dashboard.",
			},
			"updater": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID who last updated the dashboard.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp of dashboard creation.",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp of last update.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				meta := m.(*ProviderMetadata)
				dashboardID := d.Id()
				if dashboardID == "" {
					return nil, fmt.Errorf("could not determine the resource ID - possibly the ID was not set")
				}

				// Support importing all dashboards with "*"
				if dashboardID == "*" {
					dashboards, err := meta.client.GetAllDashboards()
					if err != nil {
						return nil, fmt.Errorf("could not get dashboards from API: %s", err)
					}
					results := make([]*schema.ResourceData, 0, len(dashboards))
					for _, dash := range dashboards {
						dd := resourceDashboard().Data(nil)
						dd.SetId(dash.DashboardID)
						if err := setDashboardState(dd, dash); err != nil {
							return nil, fmt.Errorf("failed to set dashboard state: %s", err)
						}
						results = append(results, dd)
					}
					return results, nil
				}

				// Support comma-separated IDs
				var dashboardIDs []string
				if strings.Contains(dashboardID, ",") {
					dashboardIDs = strings.Split(dashboardID, ",")
				} else {
					dashboardIDs = []string{dashboardID}
				}

				results := make([]*schema.ResourceData, 0, len(dashboardIDs))
				for _, id := range dashboardIDs {
					id = strings.TrimSpace(id)
					dd := resourceDashboard().Data(nil)
					resp, err := meta.client.GetDashboard(id)
					if err != nil {
						return nil, fmt.Errorf("could not get dashboard from API: %s (dashboard ID was: '%s')", err, id)
					}
					dd.SetId(resp.DashboardID)
					dashboard := Dashboard(*resp)
					if err := setDashboardState(dd, &dashboard); err != nil {
						return nil, fmt.Errorf("failed to set dashboard state: %s", err)
					}
					results = append(results, dd)
				}
				return results, nil
			},
		},
	}
}

type dashboardArgs struct {
	dashboardName string
	description   string
	tags          []string
	definition    map[string]interface{}
	diags         diag.Diagnostics
}

func parseDashboardArgs(d *schema.ResourceData) *dashboardArgs {
	args := &dashboardArgs{}

	// Required field
	if v, ok := d.GetOk("dashboard_name"); ok {
		args.dashboardName = v.(string)
	} else {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "dashboard_name is required",
		})
	}

	// Optional fields
	if v, ok := d.GetOk("description"); ok {
		args.description = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		args.tags = interfaceSliceToStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("definition"); ok {
		defStr := v.(string)
		if defStr != "" {
			def, err := stringToJSONMap(defStr)
			if err != nil {
				args.diags = append(args.diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid definition JSON",
					Detail:   err.Error(),
				})
			} else {
				args.definition = def
			}
		}
	}

	return args
}

func setDashboardState(d *schema.ResourceData, dash *Dashboard) error {
	if err := d.Set("dashboard_id", dash.DashboardID); err != nil {
		return err
	}
	if err := d.Set("dashboard_name", dash.DashboardName); err != nil {
		return err
	}
	if err := d.Set("description", dash.Description); err != nil {
		return err
	}
	if len(dash.Tags) > 0 {
		if err := d.Set("tags", stringSliceToInterface(dash.Tags)); err != nil {
			return err
		}
	}
	if dash.Definition != nil {
		defStr, err := jsonMapToString(dash.Definition)
		if err != nil {
			return err
		}
		if err := d.Set("definition", defStr); err != nil {
			return err
		}
	}
	if err := d.Set("creator", dash.Creator); err != nil {
		return err
	}
	if err := d.Set("updater", dash.Updater); err != nil {
		return err
	}
	if err := d.Set("created", dash.Created); err != nil {
		return err
	}
	if err := d.Set("updated", dash.Updated); err != nil {
		return err
	}
	return nil
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	args := parseDashboardArgs(d)
	if len(args.diags) > 0 {
		return args.diags
	}

	dashboard := &Dashboard{
		DashboardName: args.dashboardName,
		Description:   args.description,
		Tags:          args.tags,
		Definition:    args.definition,
	}

	resp, err := meta.client.CreateDashboard(dashboard)
	if err != nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not create the dashboard resource",
			Detail:   err.Error(),
		})
		return args.diags
	}

	d.SetId(resp.DashboardID)
	dashResp := Dashboard(*resp)
	if err := setDashboardState(d, &dashResp); err != nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to set dashboard state after create",
			Detail:   err.Error(),
		})
	}

	return args.diags
}

func resourceDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	var diags diag.Diagnostics

	dashboardID := d.Id()
	if dashboardID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot determine dashboard ID",
			Detail:   "Dashboard ID is required but not found in Terraform state",
		})
		return diags
	}

	resp, err := meta.client.GetDashboard(dashboardID)
	if err != nil {
		// Check if resource was deleted outside Terraform
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not read the dashboard resource",
			Detail:   err.Error(),
		})
		return diags
	}

	dashResp := Dashboard(*resp)
	if err := setDashboardState(d, &dashResp); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to set dashboard state after read",
			Detail:   err.Error(),
		})
	}

	return diags
}

func resourceDashboardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	args := parseDashboardArgs(d)
	if len(args.diags) > 0 {
		return args.diags
	}

	dashboardID := d.Id()
	if dashboardID == "" {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot determine dashboard ID for update",
			Detail:   "Dashboard ID is required but not found in Terraform state",
		})
		return args.diags
	}

	dashboard := &Dashboard{
		DashboardName: args.dashboardName,
		Description:   args.description,
		Tags:          args.tags,
		Definition:    args.definition,
	}

	resp, err := meta.client.UpdateDashboard(dashboardID, dashboard)
	if err != nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not update the dashboard resource",
			Detail:   err.Error(),
		})
		return args.diags
	}

	dashResp := Dashboard(*resp)
	if err := setDashboardState(d, &dashResp); err != nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to set dashboard state after update",
			Detail:   err.Error(),
		})
	}

	return args.diags
}

func resourceDashboardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	var diags diag.Diagnostics

	dashboardID := d.Id()
	if dashboardID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot determine dashboard ID for deletion",
			Detail:   "Dashboard ID is required but not found in Terraform state",
		})
		return diags
	}

	err := meta.client.DeleteDashboard(dashboardID)
	if err != nil {
		// If already deleted, just remove from state
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not delete the dashboard resource",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")
	return diags
}
