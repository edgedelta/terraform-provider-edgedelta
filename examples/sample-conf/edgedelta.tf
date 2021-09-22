terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta.com/edgedelta/config"
      version = "0.0.1"
    }
  }
}

resource "edgedelta_config" "my_conf" {
  org_id             = ""
  path               = ""
  debug              = true
}

resource "edgedelta_config" "my_conf_with_id" {
  conf_id            = ""
  org_id             = ""
  path               = ""
  debug              = true
}