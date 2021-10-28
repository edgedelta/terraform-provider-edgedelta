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

output "instance1_conf_id" {
  value              = edgedelta_config.conf_without_id.id
  description        = "The config ID created of the edgedelta_config instance"
}

output "instance1_tag" {
  value              = edgedelta_config.conf_without_id.tag
  description        = "The tag created of the edgedelta_config instance"
}

output "instance2_conf_id" {
  value              = edgedelta_config.config_with_id.id
  description        = "The config ID of the edgedelta_config instance"
}

output "instance2_tag" {
  value              = edgedelta_config.config_with_id.tag
  description        = "The tag of the edgedelta_config instance"
}

resource "edgedelta_config" "conf_without_id" {
  config_content     = file("/path/to/ed-config/file.yml")
}

resource "edgedelta_config" "config_with_id" {
  conf_id            = "00000000-0000-0000-0000-000000000000"
  config_content     = file("/path/to/ed-config/file.yml")
}
```

## Importing Existing Configs

You can import your existing config resources to the terraform state using `terraform import` command with a specific config ID. 

- Define a config resource
```hcl
resource "edgedelta_config" "imported_config" {

}
```

- Run `terraform import` on terminal
```bash
terraform import edgedelta_config.imported_config <resource-id>
```

## Schema

| Name           | Description                                                                                                                             | Type   | Default | Required |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------|--------|---------|----------|
| conf_id        | The pre-existing unique configuration ID. When not specified in resource schema, a new Edge Delta config will be created on the first  `terraform apply` | String | ""      | no       |
| config_content | Configuration file data                                                                                                                 | String | n/a     | yes      |

## Outputs

| Name | Description | Type |
|------|-------------|------|
| tag  | Configuration instance tag. The output value is the exact value of the `tag` key in the `config_content`. | String |
| id | When a resource is created, ID is set to the active configuration ID of the config instance. Using `id` instead of `config_id` as the configuration ID output is highly encouraged. | String