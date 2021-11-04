terraform {
    required_providers {
            edgedelta = {
            source  = "edgedelta/edgedelta"
            version = "0.0.4"
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
    name     = "correlated-signal-example-monitor"
    type     = "correlated-signal"
    enabled  = true
    payload  = file("payload.json")
    creator  = "creator@email.domain"
}
