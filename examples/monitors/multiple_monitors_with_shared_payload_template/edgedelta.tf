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

data "template_file" "pattern_skyline_payloads" {
    template = file("${path.module}/payload-template-skyline")
    for_each = var.monitor_tags
    vars = {
        tag = each.value
    }
}

data "template_file" "pattern_check_payloads" {
    template = file("${path.module}/payload-template-check")
    for_each = var.monitor_tags
    vars = {
        tag = each.value
    }
}

data "template_file" "correlated_signal_payloads" {
    template = file("${path.module}/payload-template-signal")
    for_each = var.monitor_tags
    vars = {
        tag = each.value
    }
}

resource "edgedelta_monitor" "monitor_pattern_skyline" {
    count    = length(var.monitor_tags)
    name     = "pattern-skyline-example-monitor-${var.monitor_tags[count.index]}"
    type     = "pattern-skyline"
    enabled  = true
    payload  = data.template_file.pattern_skyline_payloads[count.index].rendered
}


resource "edgedelta_monitor" "monitor_pattern_check" {
    count    = length(var.monitor_tags)
    name     = "pattern-check-example-monitor-${var.monitor_tags[count.index]}"
    type     = "pattern-check"
    enabled  = true
    payload  = data.template_file.pattern_check_payloads[count.index].rendered
}

resource "edgedelta_monitor" "monitor_correlated_signal" {
    count    = length(var.monitor_tags)
    name     = "correlated-signal-example-monitor-${var.monitor_tags[count.index]}"
    type     = "correlated-signal"
    enabled  = true
    payload  = data.template_file.correlated_signal_payloads[count.index].rendered
}
