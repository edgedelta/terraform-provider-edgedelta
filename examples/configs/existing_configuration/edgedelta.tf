terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta/edgedelta"
      version = "0.0.10"
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

output "instance_conf_id" {
  value       = edgedelta_config.conf_with_id.id
  description = "The config ID of the edgedelta_config instance"
}

output "instance_tag" {
  value       = edgedelta_config.conf_with_id.tag
  description = "The tag of the edgedelta_config instance"
}

resource "edgedelta_config" "conf_with_id" {
  conf_id        = "<your-existing-configuration-id>"
  config_content = file("/path/to/the/agent/configuration/file.yml")
  environment    = "MacOS"
  fleet_type     = "Edge"
}