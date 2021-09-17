package edgedelta

import (
	"context"
	"strconv"
	"time"

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
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path of the config file",
			},

			// Optional params
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

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Implement here

	// Timestamp is just a placeholder, the resource data id MUST be set
	// within this function, otherwise terraform will throw an exception.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Implement here

	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Implement here

	return diags
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO: Implement here

	return diags
}
