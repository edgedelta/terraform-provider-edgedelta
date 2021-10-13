terraform {
    required_providers {
            edgedelta = {
            source  = "edgedelta.com/edgedelta/config"
            version = "0.0.1"
        }
    }
}

variable "ED_API_TOKEN" {
    type = string
}

variable "monitor_names" {
    type = list(string)
}

provider "edgedelta" {
    org_id     = "<your-organization-id>"
    api_secret = var.ED_API_TOKEN
}

output "monitor_id" {
    value = edgedelta_monitor.existing_monitor_skyline.id
}

resource "edgedelta_monitor" "existing_monitor_skyline" {
    for_each = var.monitor_names
    name     = "skyline-example-monitor-${each.value}"
    type     = "pattern-skyline"
    enabled  = true
    payload  = file("payload.json") 
}

