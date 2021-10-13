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

provider "edgedelta" {
    org_id     = "<your-organization-id>"
    api_secret = var.ED_API_TOKEN
}

output "monitor_id" {
    value = edgedelta_monitor.existing_monitor_skyline.id
}

resource "edgedelta_monitor" "existing_monitor_skyline" {
    monitor_id = "<your-existing-configuration-id>"
    name       = "skyline-example-monitor"
    type       = "pattern-skyline"
    enabled    = true
    payload    = file("payload.json") 
}

