---
page_title: "edgedelta_config Resource - terraform-provider-edgedelta"
subcategory: ""
description: "The edgedelta_config resource allows you to manage your Edge Delta configurations"
  
---

# edgedelta_config (Resource)

The `edgedelta_config` resource allows you to manage your Edge Delta configurations

## Example Usage

You can define your configurations with or without specifying the configuration ID. If you have an existing configuration, using the configuration ID in the resource schema will be sufficient. If you don't specify the configuration ID in the resource schema, the provider will create a new Edge Delta configuration and save the new configuration's ID in Terraform state. 

```bash
export TF_VAR_ED_API_TOKEN="<your-api-token-goes-here>"
```

```hcl
variable "ED_API_TOKEN" {
  type = string
}

provider "edgedelta" {
  org_id             = "22222222-2222-2222-2222-222222222222"
  api_secret         = var.ED_API_TOKEN
}

resource "edgedelta_config" "conf_without_id" {
  config_content     = file("/path/to/ed-config/file.yml")
}

resource "edgedelta_config" "config_with_id" {
  conf_id            = "00000000-0000-0000-0000-000000000000"
  config_content     = file("/path/to/ed-config/file.yml")
}
```

## Schema

| Name           | Description                                                                                                                             | Type   | Default | Required |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------|--------|---------|----------|
| conf_id        | Unique configuration ID. When not specified in resource schema, a new Edge Delta config will be created on the first  `terraform apply` | String | ""      | no       |
| config_content | Configuration file data                                                                                                                 | String | n/a     | yes      |