---
page_title: "edgedelta_config Resource - terraform-provider-edgedelta"
subcategory: ""
description: "The config resource allows you to manage your Edge Delta configurations"
  
---

# edgedelta_config (Resource)

The `config` resource allows you to manage your Edge Delta configurations

## Example Usage

You can define your configurations with or without specifying the configuration ID. If you have an existing configuration, using the configuration ID in the resource schema will be sufficient. If you don't specify the configuration ID in the resource schema, the provider will create a new Edge Delta configuration and save the new configuration's ID in Terraform state. 

```bash
export TF_VAR_ED_API_SECRET="<your-api-secret-goes-here>"
```

```hcl
variable "ED_API_SECRET" {
  type = string
}

resource "edgedelta_config" "bare_minimum" {
  org_id             = "22222222-2222-2222-2222-222222222222"
  config_content     = file("/path/to/ed-config/file.yml")
  api_secret         = var.ED_API_SECRET
}

resource "edgedelta_config" "config_with_id" {
  conf_id            = "00000000-0000-0000-0000-000000000000"
  org_id             = "11111111-1111-1111-1111-111111111111"
  config_content     = file("/path/to/ed-config/file.yml")
  api_endpoint       = "https://api.edgedelta.com"
  api_secret         = var.ED_API_SECRET
}

resource "edgedelta_config" "conf_without_id" {
  org_id             = "22222222-2222-2222-2222-222222222222"
  config_content     = file("/path/to/ed-config/file.yml")
  api_endpoint       = "https://api.edgedelta.com"
  api_secret         = var.ED_API_SECRET
}
```

## Schema

### Required

- **api_secret** (String, Sensitive) API secret. User is **highly encouraged** to use terraform variables to pass the secret value in resource schema.
- **config_content** (String) Configuration file data.
- **org_id** (String) Unique organization ID.

### Optional

- **api_endpoint** (String) API base URL, default is `https://api.edgedelta.com`
- **conf_id** (String) Unique configuration ID. When not specified in resource schema, a new Edge Delta config will be created on the first `terraform apply`.
- **id** (String) The ID of this resource. Only set on create, read-only.


