# Provider Eng Documentation

This document provides detailed information about the provider development and deployment processes.

## Developing the Provider

The Edge Delta provider consists of three main parts, namely the `Provider` and `Resource` definitions, and the `API Client`.

### Provider Struct

Provider struct is the entrypoint of the `terraform-provider-edgedelta`. The provider definition in [provider.go](../edgedelta/provider.go) defines the provider inputs, resources, data sources and context configuration function. At the time being, we don't have any data sources specified in the provider definition.

#### Schema

> Type: `map[string]*schema.Schema`

The provider schema defines the provider configuration parameters and parameter attributes. Commonly used parameter attributes are as follows:

| Name | Type | Description |
|------|------|-------------|
|Type|int|Parameter type|
|Description|string|Parameter description|
|Required|bool|When set to true, marks the parameter as required|
|Optional|bool|When set to true, marks the parameter as optional|
|Computed|bool|When set to true, marks the parameter as computed and the value of the parameter is dynamically set within the CRUD functions. Currently used only with `tag` of `edgedelta_config` [resource](#resource-struct).|
|Sensitive|bool|When set to true, marks the parameter as sensitive. Currently used only with `api_secret`|


The current provider params are as follows:

| Name         | Description                                                                                                        | Type               | Default                   | Required |
|--------------|--------------------------------------------------------------------------------------------------------------------|--------------------|---------------------------|----------|
| api_secret   | API token. User is  **highly encouraged**  to use terraform variables to pass the token value in resource schema | String,  Sensitive | n/a                       | yes      |
| org_id       | Unique organization ID                                                                                             | String             | n/a                       | yes      |
| api_endpoint | API base URL                                                                                                       | String             | https://api.edgedelta.com | no       |

#### ResourcesMap

> Type: `map[string]*schema.Resource`

The resources map provides a mapping of resource names and resource schemas. We currently have two resources, namely `edgedelta_config` and `edgedelta_monitor`, and the current resource maps are as follows:

```go
map[string]*schema.Resource{
    "edgedelta_config": {
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
    },
    "edgedelta_monitor": {
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
	}
}
```

#### ConfigureContextFunc

The context configuration function is used to initialize the metadata struct instance which is passed to the CRUD functions of the resources when the provider is invoked by Terraform CLI. The metadata struct holds the information used by every resource and CRUD function. The current metadata struct includes only an instance to `ConfigAPIClient` is as follows:

```go
type ProviderMetadata struct {
    client ConfigAPIClient
}
```

The metadata instance is initialized in the `providerConfigure` function in [edgedelta/provider.go](../edgedelta/provider.go)

### Resource Struct

The resource struct defines the parameters and CRUD functions of that particular resource. We currently have two resources, namely `edgedelta_config` and `edgedelta_monitor`, and are defined in [resources_config.go](../edgedelta/resources_config.go) and [resources_monitor.go](../edgedelta/resources_monitor.go) respectively.

Each resource instance has an additional `id` field, not defined explicitly in the resource struct, which is an unique identifier of the instance. In our implementation, we have used the `id` field to hold the configuration id and monitor id data, as well as the `conf_id` and `monitor_id` params. This design choice is made to prevent the `conf_id` and `monitor_id` to be set to `nil` every time `terraform apply` is used with a resource with no explicit `conf_id` or `monitor_id`. The `id` then holds the real resource ID after the creation process to later use in the API calls.

#### Schema

> Type: `map[string]*schema.Schema`

> Schema is the same data structure as provider's schema. Further information about this data structure can be found in [Provider Schema](#schema) section.

An example schema from the [edgedelta_config](../edgedelta/resources_config.go) resource can be found below:

```go
{
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
}
```

#### CRUD Functions

Every resource has 4 functions to control the resource state: create, read, update and delete. These functions are provided to the resource struct with the fields respectively `CreateContext`, `ReadContext`, `UpdateContext` and `DeleteContext`. Each function takes 3 arguments: `context` (`context.Context`), `data` (`*schema.ResourceData`) and `metadata` (`interface{}` in general, or `*ProviderMetadata` in our implementation). The create and read functions should set the parameter values in the `data` argument. These values are beging used to update the Terraform state. 

### API Client

The API client is a minimal SDK that provides the functionality to create and update the resources in Edge Delta's side. The API client is a struct defined in [api_client.go](../edgedelta/api_client.go) which definition can be seen below:

```go
type APIClient struct {
	OrgID      string
	APIBaseURL string
	apiSecret  string
	cl         *http.Client
}
```

The `OrgID`, `APIBaseURL` and `apiSecret` params should be passed in instantiation, then the `initializeHTTPClient` function should be called to initialize the http client `cl`.

The client has a number of functions, the detailed function information can be found in the table below:

|Name|API Resource Tag|Params|Return Value|
|-|-|-|-|
|GetConfigWithID|`confs`|**configID**: `string`|[*GetConfigResponse](../edgedelta/types.go)|
|CreateConfig|`confs`|**configObject**: [Config](../edgedelta/types.go)|[*CreateConfigResponse](../edgedelta/types.go)|
|UpdateConfigWithID|`confs`|**configID**: `string` <br><br>  **configObject**: [Config](../edgedelta/types.go)|[*UpdateConfigResponse](../edgedelta/types.go)|
|GetMonitorWithID|`monitors`|**monitorID**: `string`|[*GetMonitorResponse](../edgedelta/types.go)|
|CreateMonitor|`monitors`|**monitor**: [Monitor](../edgedelta/types.go)|[*CreateMonitorResponse](../edgedelta/types.go)|
|UpdateMonitorWithID|`monitors`|**monitorID**: `string` <br><br>  **monitor**: [Monitor](../edgedelta/types.go)|[*UpdateMonitorResponse](../edgedelta/types.go)|

## Running the Provider Locally

### Building

To build and use the provider on your local machine, use the Makefile with the environment variables below:

* `TERRAFORM_PROVIDER_ED_OS_ARCH`: Desired OS architecture. This has to be the arch of your local machine (Ex. darwin_amd64 for 64-bit Macs)
* `TERRAFORM_PROVIDER_ED_VERSION`: Desired provider version. This has no direct effect on the functionlity of the provider itself, you will use this version number in your `.tf` file later. (Ex. 0.4.2, 1.53.1 etc.)

To build the provider, run the command below:

```bash
TERRAFORM_PROVIDER_ED_VERSION=<desired-version> TERRAFORM_PROVIDER_ED_OS_ARCH=<desired-arch> make
```

Then the provider will be built and the binary will be put under `~/.terraform.d/plugins/edgedelta.com/local/edgedelta/${VERSION}/${OS_ARCH}`.

### Using

To test the provider on your local machine, use the name `edgedelta.com/local/edgedelta` for the provider name and the version you have provided in `TERRAFORM_PROVIDER_ED_VERSION` before for the version in your `.tf` file:

```hcl
terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta.com/local/edgedelta"
      version = "0.0.6"   # If you've set TERRAFORM_PROVIDER_ED_VERSION=0.0.6
    }
  }
}
```

## Publishing the Provider

You need to follow a number of simple steps to publish the provider. Publishing the provider to Terraform Registry is a pretty straightforward process. The steps described here are mainly taken from [the official Terraform documentation](https://www.terraform.io/docs/registry/providers/publishing.html).

### Preparing the Releaser

* Generate a signing key following the instructions [here](https://www.terraform.io/docs/registry/providers/publishing.html#preparing-and-adding-a-signing-key). You will later use this key to sign the provider releases.
* Create GoReleaser and GitHub Actions configuration files in `.goreleaser.yml` and `.github/workflows/release.yml` in the project repository, respectively. There is an exisitng configuration in the Edge Delta provider repo at the time being, however if you're planning to create one, [this GoReleaser config](https://github.com/hashicorp/terraform-provider-scaffolding/blob/main/.goreleaser.yml) and [this GitHub Actions config](https://github.com/hashicorp/terraform-provider-scaffolding/blob/main/.github/workflows/release.yml) are a good examples.
* Add the ASCII-armored GPG private key and key password to the repo secrets. Detailed instructions can be found in the 4th step of [this listing](https://www.terraform.io/docs/registry/providers/publishing.html#github-actions-preferred-).
* To test the functionality of the GitHub Actions releaser, push a new version tag.

### Publishing the Release to Registry

* Go to [Terraform Registry](https://registry.terraform.io/) and sign-in with your GitHub account.
* Authorize the registry for Edge Delta
* Add your ASCII-armored public key, which you have generated before, to the Terraform Registry
* Go to [publish page](https://registry.terraform.io/publish/provider), or [Edge Delta provider page](https://registry.terraform.io/providers/edgedelta/edgedelta) to publish the provider for the new versions.