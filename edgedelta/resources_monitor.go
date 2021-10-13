package edgedelta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	monitorTypes = [...]string{
		"pattern-check",
		"pattern-skyline",
		"correlated-signal",
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
			// Optional params
			"monitor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique monitor ID",
			},
		},
	}
}

func resourceMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
