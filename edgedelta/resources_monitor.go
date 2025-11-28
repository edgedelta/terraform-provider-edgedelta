package edgedelta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	monitorTypes = [...]string{
		"pattern-check",
		"pattern-skyline",
		"correlated-signal",
		"metric-alert",
	}
)

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		CreateContext: resourceMonitorCreate,
		ReadContext:   resourceMonitorRead,
		UpdateContext: resourceMonitorUpdate,
		DeleteContext: resourceMonitorDelete,
		Schema: map[string]*schema.Schema{
			// Required params
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor name",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor type",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					for _, t := range monitorTypes {
						if v == t {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of the values from %v, got: %s", key, monitorTypes, v))
					return
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Monitor enabled flag",
			},
			"payload": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor payload",
			},
			"creator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor creator (email)",
			},
			// Optional params
			"monitor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique monitor ID",
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				meta := m.(*ProviderMetadata)
				monitorID := d.Id()
				if monitorID == "" { // monitorID DNE
					return nil, fmt.Errorf("could not determine the resource ID - possibly the ID was not set")
				}
				var monitorIDs []string
				if monitorID == "*" {
					monitors, err := meta.client.GetAllMonitors()
					if err != nil {
						return nil, fmt.Errorf("could not get the monitors from API: %s", err)
					}
					results := make([]*schema.ResourceData, 0, len(monitors))
					for _, m := range monitors {
						dd := resourceMonitor().Data(nil)
						dd.SetId(m.ID)
						if err := dd.Set("name", m.Name); err != nil {
							return nil, fmt.Errorf("failed to set name: %s", err)
						}
						if err := dd.Set("type", m.Type); err != nil {
							return nil, fmt.Errorf("failed to set type: %s", err)
						}
						if err := dd.Set("monitor_id", m.ID); err != nil {
							return nil, fmt.Errorf("failed to set monitor_id: %s", err)
						}
						if err := dd.Set("enabled", m.Enabled); err != nil {
							return nil, fmt.Errorf("failed to set enabled: %s", err)
						}
						if err := dd.Set("payload", m.Payload); err != nil {
							return nil, fmt.Errorf("failed to set payload: %s", err)
						}
						if err := dd.Set("creator", m.Creator); err != nil {
							return nil, fmt.Errorf("failed to set creator: %s", err)
						}
						results = append(results, dd)
					}
					return results, nil
				} else if strings.Contains(monitorID, ",") {
					monitorIDs = strings.Split(monitorID, ",")
				} else {
					monitorIDs = []string{monitorID}
				}
				results := make([]*schema.ResourceData, 0, len(monitorIDs))
				for _, id := range monitorIDs {
					dd := resourceMonitor().Data(nil)
					resp, err := meta.client.GetMonitorWithID(id)
					if err != nil {
						return nil, fmt.Errorf("could not get the resource data from API: %s (resource ID was: '%s')", err, id)
					}
					dd.SetId(resp.ID)
					if err := dd.Set("monitor_id", id); err != nil {
						return nil, fmt.Errorf("failed to set monitor_id: %s", err)
					}
					if err := dd.Set("name", resp.Name); err != nil {
						return nil, fmt.Errorf("failed to set name: %s", err)
					}
					if err := dd.Set("type", resp.Type); err != nil {
						return nil, fmt.Errorf("failed to set type: %s", err)
					}
					if err := dd.Set("enabled", resp.Enabled); err != nil {
						return nil, fmt.Errorf("failed to set enabled: %s", err)
					}
					if err := dd.Set("payload", resp.Payload); err != nil {
						return nil, fmt.Errorf("failed to set payload: %s", err)
					}
					if err := dd.Set("creator", resp.Creator); err != nil {
						return nil, fmt.Errorf("failed to set creator: %s", err)
					}
					results = append(results, dd)
				}
				return results, nil
			},
		},
	}
}

func parseMonitorArgs(d *schema.ResourceData) (monitorID string, name string, mType string, enabled bool, payload string, creator string, diags diag.Diagnostics) {
	monitorIDRaw := d.Get("monitor_id")
	nameRaw := d.Get("name")
	typeRaw := d.Get("type")
	enabledRaw := d.Get("enabled")
	payloadRaw := d.Get("payload")
	creatorRaw := d.Get("creator")
	if monitorIDRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "monitor_id is nil",
		})
	} else {
		monitorID = monitorIDRaw.(string)
	}
	if nameRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "name is nil",
		})
	} else {
		name = nameRaw.(string)
	}
	if typeRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "type is nil",
		})
	} else {
		mType = typeRaw.(string)
	}
	if enabledRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "enabled is nil",
		})
	} else {
		enabled = enabledRaw.(bool)
	}
	if payloadRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "payload is nil",
		})
	} else {
		payload = payloadRaw.(string)
	}
	if creatorRaw == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "creator is nil",
		})
	} else {
		creator = creatorRaw.(string)
	}
	return
}

func resourceMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	monitorID, name, mType, enabled, payload, creator, diags := parseMonitorArgs(d)
	if len(diags) > 0 {
		return diags
	}
	mon := Monitor{
		Enabled: enabled,
		Name:    name,
		Payload: payload,
		Type:    mType,
		Creator: creator,
	}
	if monitorID == "" {
		// Create a new monitor
		resp, err := meta.client.CreateMonitor(mon)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not create the monitor resource",
				Detail:   fmt.Sprintf("%s", err),
			})
			return diags
		}
		d.SetId(resp.ID)
		diags = setWithError(d, "monitor_id", monitorID, diags)
		diags = setWithError(d, "name", resp.Name, diags)
		diags = setWithError(d, "type", resp.Type, diags)
		diags = setWithError(d, "enabled", resp.Enabled, diags)
		diags = setWithError(d, "payload", resp.Payload, diags)
		diags = setWithError(d, "creator", resp.Creator, diags)
	} else {
		// First run of the terraform apply, just update the existing monitor
		resp, err := meta.client.UpdateMonitorWithID(monitorID, mon)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not update the monitor resource (create=>update)",
				Detail:   fmt.Sprintf("%s", err),
			})
			return diags
		}
		d.SetId(resp.ID)
		diags = setWithError(d, "monitor_id", resp.ID, diags)
		diags = setWithError(d, "name", resp.Name, diags)
		diags = setWithError(d, "type", resp.Type, diags)
		diags = setWithError(d, "enabled", resp.Enabled, diags)
		diags = setWithError(d, "payload", resp.Payload, diags)
	}
	return diags
}

func resourceMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	monitorID, _, _, _, _, _, diags := parseMonitorArgs(d)
	if len(diags) > 0 {
		return diags
	}
	activeMonID := monitorID
	if activeMonID == "" {
		if d.Id() == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot determine monitor ID",
				Detail:   "Possibly tried to read the resource with monitor_id=nil and id=nil",
			})
			return diags
		}
		activeMonID = d.Id()
	}
	resp, err := meta.client.GetMonitorWithID(activeMonID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not get the resource data from API",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}
	d.SetId(resp.ID)
	diags = setWithError(d, "monitor_id", monitorID, diags)
	diags = setWithError(d, "name", resp.Name, diags)
	diags = setWithError(d, "type", resp.Type, diags)
	diags = setWithError(d, "enabled", resp.Enabled, diags)
	diags = setWithError(d, "payload", resp.Payload, diags)
	diags = setWithError(d, "creator", resp.Creator, diags)
	return diags
}

func resourceMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	monitorID, name, mType, enabled, payload, creator, diags := parseMonitorArgs(d)
	if len(diags) > 0 {
		return diags
	}
	if monitorID == "" {
		// Just get the monitor id from the tf state
		monitorID = d.Id()
	}
	mon := Monitor{
		Enabled: enabled,
		Name:    name,
		Payload: payload,
		Type:    mType,
		Creator: creator,
	}
	_, err := meta.client.UpdateMonitorWithID(monitorID, mon)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not update the monitor resource",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}
	return diags
}

func resourceMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	monitorID, _, _, _, _, _, diags := parseMonitorArgs(d)
	if len(diags) > 0 {
		return diags
	}
	if monitorID == "" {
		// Just get the monitor id from the tf state
		monitorID = d.Id()
	}
	err := meta.client.DeleteMonitorWithID(monitorID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not delete the monitor resource",
			Detail:   fmt.Sprintf("%s", err),
		})
		return diags
	}
	return diags
}
