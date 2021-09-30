terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta.com/edgedelta/config"
      version = "0.0.1"
    }
  }
}

variable "ED_API_SECRET" {
  type = string
}

provider "edgedelta" {
  org_id             = "<your-organization-id>"
  api_secret         = var.ED_API_SECRET
}

resource "edgedelta_config" "bare_minimum" {
  config_content     = file("/path/to/the/agent/configuration/file.yml")
}