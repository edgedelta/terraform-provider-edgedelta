package edgedelta

import (
	"context"
	"fmt"
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
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment where the pipeline will be deployed (Kubernetes, Linux, Windows, MacOS, Docker)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validEnvs := []EnvironmentType{
						KubernetesEnvironmentType,
						HelmEnvironmentType,
						DockerEnvironmentType,
						MacOSEnvironmentType,
						LinuxEnvironmentType,
						WindowsEnvironmentType,
					}
					validEnvStrings := make([]string, len(validEnvs))
					for i, env := range validEnvs {
						if v == string(env) {
							return
						}
						validEnvStrings[i] = string(env)
					}
					errs = append(errs, fmt.Errorf("%q must be one of: %v, got: %s", key, validEnvStrings, v))
					return
				},
			},
			"fleet_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(EdgeFleetType),
				Description: "Fleet type (Edge, Cloud). Defaults to Edge.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validFleetTypes := []FleetType{
						EdgeFleetType,
						CloudFleetType,
					}
					validFleetTypeStrings := make([]string, len(validFleetTypes))
					for i, ft := range validFleetTypes {
						if v == string(ft) {
							return
						}
						validFleetTypeStrings[i] = string(ft)
					}
					errs = append(errs, fmt.Errorf("%q must be one of: %v, got: %s", key, validFleetTypeStrings, v))
					return
				},
			},
			// Optional params
			"conf_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique configuration ID",
			},
			"fleet_subtype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fleet subtype (Edge, Coordinator, Gateway). Required when environment is Kubernetes and fleet_type is Edge.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						return // Optional field, empty is allowed
					}
					validFleetSubtypes := []FleetSubtype{
						EdgeFleetSubtype,
						CoordinatorFleetSubtype,
						GatewayFleetSubtype,
					}
					validFleetSubtypeStrings := make([]string, len(validFleetSubtypes))
					for i, fst := range validFleetSubtypes {
						if v == string(fst) {
							return
						}
						validFleetSubtypeStrings[i] = string(fst)
					}
					errs = append(errs, fmt.Errorf("%q must be one of: %v, got: %s", key, validFleetSubtypeStrings, v))
					return
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster name for grouping pipelines",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the pipeline",
			},
			"auto_deploy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Automatically deploy the config after saving. If false, only saves the config.",
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
						if err := dd.Set("conf_id", c.ID); err != nil {
							return nil, fmt.Errorf("failed to set conf_id: %s", err)
						}
						if err := dd.Set("tag", c.Tag); err != nil {
							return nil, fmt.Errorf("failed to set tag: %s", err)
						}
						if err := dd.Set("config_content", c.Content); err != nil {
							return nil, fmt.Errorf("failed to set config_content: %s", err)
						}
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
					if err := dd.Set("conf_id", id); err != nil {
						return nil, fmt.Errorf("failed to set conf_id: %s", err)
					}
					if err := dd.Set("tag", resp.Tag); err != nil {
						return nil, fmt.Errorf("failed to set tag: %s", err)
					}
					if err := dd.Set("config_content", resp.Content); err != nil {
						return nil, fmt.Errorf("failed to set config_content: %s", err)
					}
					results = append(results, dd)
				}
				return results, nil
			},
		},
	}
}

type configArgs struct {
	confID       string
	confData     string
	environment  EnvironmentType
	fleetType    FleetType
	fleetSubtype FleetSubtype
	clusterName  string
	description  string
	autoDeploy   bool
	diags        diag.Diagnostics
}

func parseArgs(d *schema.ResourceData) *configArgs {
	args := &configArgs{}
	confIDRaw := d.Get("conf_id")
	configDataRaw := d.Get("config_content")
	environmentRaw := d.Get("environment")
	fleetTypeRaw := d.Get("fleet_type")
	fleetSubtypeRaw := d.Get("fleet_subtype")
	clusterNameRaw := d.Get("cluster_name")
	descriptionRaw := d.Get("description")
	autoDeployRaw := d.Get("auto_deploy")

	if confIDRaw != nil {
		args.confID = confIDRaw.(string)
	}
	if configDataRaw == nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "config_content is required",
		})
	} else {
		args.confData = configDataRaw.(string)
	}
	if environmentRaw == nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "environment is required",
		})
	} else {
		args.environment = EnvironmentType(environmentRaw.(string))
	}
	if fleetTypeRaw != nil {
		args.fleetType = FleetType(fleetTypeRaw.(string))
	} else {
		args.fleetType = EdgeFleetType // Default
	}
	if fleetSubtypeRaw != nil && fleetSubtypeRaw.(string) != "" {
		args.fleetSubtype = FleetSubtype(fleetSubtypeRaw.(string))
	}
	if clusterNameRaw != nil {
		args.clusterName = clusterNameRaw.(string)
	}
	if descriptionRaw != nil {
		args.description = descriptionRaw.(string)
	}
	if autoDeployRaw != nil {
		args.autoDeploy = autoDeployRaw.(bool)
	}
	return args
}

func setWithError(d *schema.ResourceData, key string, value any, diags diag.Diagnostics) diag.Diagnostics {
	if err := d.Set(key, value); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Failed to set %s in state", key),
			Detail:   fmt.Sprintf("%s", err),
		})
	}
	return diags
}

func saveAndDeployConfig(client APIClient, confID string, saveReq SaveRequest, autoDeploy bool, errorContext string) (*SaveConfigResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Step 1: Save the config
	saveResp, err := client.SaveConfig(confID, saveReq)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Could not save the config resource%s", errorContext),
			Detail:   fmt.Sprintf("%s", err),
		})
		return nil, diags
	}

	// Step 2: Conditionally deploy if auto_deploy is true
	if autoDeploy {
		// Get the latest config history version (timestamp) after save
		version, err := client.GetLatestConfigHistoryVersion(confID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not get latest config history version for deployment",
				Detail:   fmt.Sprintf("%s", err),
			})
			return nil, diags
		}

		// Deploy the saved version
		_, err = client.DeployConfig(confID, version)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not deploy the config resource",
				Detail:   fmt.Sprintf("Config was saved but deployment failed: %s", err),
			})
			return nil, diags
		}
	}

	return saveResp, diags
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	args := parseArgs(d)
	if len(args.diags) > 0 {
		return args.diags
	}
	confDataObj := Config{
		Content:      args.confData,
		Environment:  args.environment,
		FleetType:    args.fleetType,
		FleetSubtype: args.fleetSubtype,
		ClusterName:  args.clusterName,
		Description:  args.description,
	}
	if args.confID == "" {
		// Create a new config
		apiResp, err := meta.client.CreateConfig(confDataObj)
		if err != nil {
			args.diags = append(args.diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Could not create the config resource",
				Detail:   fmt.Sprintf("%s", err),
			})
			return args.diags
		}
		d.SetId(apiResp.ID)
		args.diags = setWithError(d, "conf_id", args.confID, args.diags)
		args.diags = setWithError(d, "tag", apiResp.Tag, args.diags)

	} else {
		// First run of the terraform config, save the existing ed-config
		saveReq := SaveRequest{
			Content:     &args.confData,
			Description: args.description,
		}
		saveResp, saveDiags := saveAndDeployConfig(meta.client, args.confID, saveReq, args.autoDeploy, " (create=>save)")
		if len(saveDiags) > 0 {
			args.diags = append(args.diags, saveDiags...)
			return args.diags
		}

		d.SetId(saveResp.ID)
		args.diags = setWithError(d, "conf_id", saveResp.ID, args.diags)
		// Get the full config to get tag
		configResp, err := meta.client.GetConfigWithID(args.confID)
		if err == nil {
			args.diags = setWithError(d, "tag", configResp.Tag, args.diags)
		}
	}

	return args.diags
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	args := parseArgs(d)
	if len(args.diags) > 0 {
		return args.diags
	}
	activeConfID := args.confID
	if activeConfID == "" {
		if d.Id() == "" {
			args.diags = append(args.diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot determine config ID",
				Detail:   "Possibly tried to read the resource with conf_id=nil and id=nil",
			})
			return args.diags
		}

		activeConfID = d.Id()
	}
	apiResp, err := meta.client.GetConfigWithID(activeConfID)
	if err != nil {
		args.diags = append(args.diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not get the resource data from API",
			Detail:   fmt.Sprintf("%s", err),
		})
		return args.diags
	}
	d.SetId(apiResp.ID)
	args.diags = setWithError(d, "conf_id", args.confID, args.diags)
	args.diags = setWithError(d, "tag", apiResp.Tag, args.diags)

	return args.diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := m.(*ProviderMetadata)
	args := parseArgs(d)
	if len(args.diags) > 0 {
		return args.diags
	}
	confID := args.confID
	if confID == "" {
		// Just get the config id from the tf state
		confID = d.Id()
	}

	// Save and optionally deploy the config
	saveReq := SaveRequest{
		Content:     &args.confData,
		Description: args.description,
	}
	_, saveDiags := saveAndDeployConfig(meta.client, confID, saveReq, args.autoDeploy, "")
	if len(saveDiags) > 0 {
		args.diags = append(args.diags, saveDiags...)
		return args.diags
	}

	return args.diags
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
