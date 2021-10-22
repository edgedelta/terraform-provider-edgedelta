---
page_title: "edgedelta_monitor Resource - terraform-provider-edgedelta"
subcategory: ""
description: "The edgedelta_monitor resource allows you to manage your Edge Delta alert definitons"
  
---

# edgedelta_monitor (Resource)

The `edgedelta_monitor` resource allows you to manage your monitors (alert definitions)

## Example Usage

You can define your monitors with or without specifying the monitor ID. If you have an existing monitor, using the monitor ID in the resource schema will be sufficient. If you don't specify the monitor ID in the resource schema, the provider will create a new alert definition using the name and the payload specified, and save the new monitor's ID in Terraform state. 

```bash
export TF_VAR_ED_API_TOKEN="<your-api-token-goes-here>"
```

```hcl
variable "ED_API_TOKEN" {
  type = string
}

provider "edgedelta" {
    org_id     = "<your-organization-id>"
    api_secret = var.ED_API_TOKEN
}

output "bare_minimum_monitor_id" {
    value = edgedelta_monitor.bare_minimum.id
}

output "existing_monitor_id" {
    value = edgedelta_monitor.existing_monitor.id
}

resource "edgedelta_monitor" "bare_minimum" {
    name    = "example-monitor"
    type    = "pattern-skyline"
    payload = file("payload1.json") 
    creator = "creator-mail@example.org"
}

resource "edgedelta_monitor" "existing_monitor" {
    monitor_id = "<your-monitor-id>"
    name       = "existing-monitor"
    type       = "pattern-skyline"
    payload    = file("payload2.json") 
    creator = "creator-mail@example.org"
}
```

## Schema

|Name|Description|Type|Default|Required|
|-|-|-|-|-|
|name|Name of the monitor|String|n/a|yes|
|type|Type of the monitor. Must be one of `pattern-check`, `pattern-skyline` and `correlated-signal`.|String|n/a|yes|
|payload|The monitor payload provides detailed information about the alert defintion. Payload must be in `JSON` format. The schema of the payload may be differ depending on the type of the monitor.|String|n/a|yes|
|creator|Mail address of the monitor creator.|String|n/a|yes|
|enabled|Monitor activity flag. When set to `false`, the monitor will be flagged as inactive, and active otherwise.|Boolean|n/a|yes|
|monitor_id|Unique ID of the monitor. If not provided, a new monitor will be created by the provider with the specified name, type and the payload.|String|""|no|

### Outputs

|Name|Description|Type|
|-|-|-|
|id|When a resource is created, the ID is set to the active monitor ID of the monitor instance.|String|