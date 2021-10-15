terraform {
    required_providers {
            edgedelta = {
            source  = "edgedelta.com/edgedelta/config"
            version = "0.0.1"
        }
    }
}

provider "edgedelta" {
    org_id     = "<your-organization-id>"
    api_secret = var.ED_API_TOKEN
}

variable "ED_API_TOKEN" {
    type = string
}

variable "monitor_tags" {
    type = list(string)
}

data "template_file" "pattern_skyline_payloads" {
    template = file("${path.module}/payload-template")
    for_each = var.monitor_tags
    vars = {
        tag = each.value
    }
}

resource "edgedelta_monitor" "example_monitors" {
    count    = length(var.monitor_tags)
    name     = "pattern-skyline-example-monitor-${var.monitor_tags[count.index]}"
    type     = "pattern-skyline"
    enabled  = true
    payload  = data.template_file.pattern_skyline_payloads[count.index].rendered
}
