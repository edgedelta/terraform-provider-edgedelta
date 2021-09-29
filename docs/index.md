---
page_title: "Edge Delta Provider"
subcategory: ""
description: "Terraform provider for managing Edge Delta configurations"
  
---

# Edge Delta Provider

Terraform provider for managing Edge Delta configurations

## Example Usage

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