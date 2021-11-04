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
  org_id             = "<your-organization-id>"
  api_secret         = var.ED_API_TOKEN
}

output "instance_conf_id" {
  value              = edgedelta_config.bare_minimum.id
  description        = "The config ID created of the edgedelta_config instance"
}

output "instance_tag" {
  value              = edgedelta_config.bare_minimum.tag
  description        = "The tag created of the edgedelta_config instance"
}

resource "edgedelta_config" "bare_minimum" {
  config_content     = file("/path/to/the/agent/configuration/file.yml")
}