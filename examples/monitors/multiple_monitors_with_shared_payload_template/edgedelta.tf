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

variable "monitor_tags" {
    type = list(string)
}

data "template_file" "monitor_payloads" {
    template = file("${path.module}/payload-template")
    for_each = var.monitor_tags
    vars = {
        tag = each.value
    }
}

resource "edgedelta_monitor" "existing_monitor_skyline" {
    count    = length(var.monitor_tags)
    name     = "skyline-example-monitor-${var.monitor_tags[count.index]}"
    type     = "pattern-skyline"
    enabled  = true
    payload  = data.template_file.monitor_payloads[count.index].rendered
}
