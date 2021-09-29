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
  org_id             = ""
  api_secret         = var.ED_API_SECRET
  api_endpoint       = "https://api.edgedelta.com"
}

resource "edgedelta_config" "conf_without_id" {
  config_content     = file("")
}

resource "edgedelta_config" "conf_with_id" {
  conf_id            = ""
  config_content     = file("")
}