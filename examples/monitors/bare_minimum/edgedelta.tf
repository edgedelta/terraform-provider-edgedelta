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
    value = edgedelta_monitor.bare_minimum_skyline.id
}

resource "edgedelta_monitor" "bare_minimum_skyline" {
    name    = "skyline-example-monitor"
    type    = "pattern-skyline"
    payload = file("payload.json") 
}

