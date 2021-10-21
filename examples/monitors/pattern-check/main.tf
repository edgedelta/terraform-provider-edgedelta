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
    value = edgedelta_monitor.example_monitor.id
}

resource "edgedelta_monitor" "example_monitor" {
    name     = "pattern-check-example-monitor"
    type     = "pattern-check"
    enabled  = true
    payload  = file("payload.json")
}