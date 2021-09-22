terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta.com/edgedelta/config"
      version = "0.0.1"
    }
  }
}

resource "edgedelta_config" "my_conf" {
  conf_id            = ""
  org_id             = ""
  config_data        = file("")
  debug              = true
}