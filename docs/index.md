---
page_title: "Edge Delta Provider"
subcategory: ""
description: "Terraform provider for managing Edge Delta configurations"
  
---

# Edge Delta Provider

Terraform provider for managing Edge Delta configurations

## Example Usage

There are multiple examples included in the `examples/` folder. A simple usage is as follows:

```bash
export TF_VAR_ED_API_SECRET="<your-api-secret-goes-here>"
```

```hcl
variable "ED_API_SECRET" {
  type = string
}

provider "edgedelta" {
  org_id             = "22222222-2222-2222-2222-222222222222"
  api_secret         = var.ED_API_SECRET
}
```

## Schema

| Name         | Description                                                                                                        | Type               | Default                   | Required |
|--------------|--------------------------------------------------------------------------------------------------------------------|--------------------|---------------------------|----------|
| api_secret   | API secret. User is  **highly encouraged**  to use terraform variables to pass the secret value in resource schema | String,  Sensitive | n/a                       | yes      |
| org_id       | Unique organization ID                                                                                             | String             | n/a                       | yes      |
| api_endpoint | API base URL                                                                                                       | String             | https://api.edgedelta.com | no       |

## Requirements

### Software

* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.0
* [terraform-provider-edgedelta](https://github.com/edgedelta/terraform-provider-edgedelta) plugin >= 0.0.1

### Permissions

The plugin uses token authentication to access the API. In order to use this plugin, you must have an [Edge Delta](https://edgedelta.com) account and save your API secret somewhere secure. Then, you need to set the `api_secret` parameter of the provider in your `.tf` files using your API secret.

We **highly encourage** you to set the `api_secret` through an environment variable. In the example configurations, we have used `TF_VAR_ED_API_SECRET` as the environment variable. You can see the example usage in the `examples/` folder in the project root.