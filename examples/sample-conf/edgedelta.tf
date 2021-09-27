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

resource "edgedelta_config" "conf_without_id" {
  org_id             = ""
  config_content     = file("")
  debug              = true
  api_endpoint       = "https://api.edgedelta.com"
  api_secret         = var.ED_API_SECRET
}

resource "edgedelta_config" "conf_with_id" {
  conf_id            = ""
  org_id             = ""
  config_content     = file("")
  debug              = true
  api_endpoint       = "https://api.edgedelta.com"
  api_secret         = var.ED_API_SECRET
}