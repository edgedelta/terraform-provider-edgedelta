terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta.com/edgedelta/config"
      version = "0.0.1"
    }
  }
}

resource "edgedelta_config" "my_conf" {
  conf_id            = "123"
  org_id             = "123123"
  path               = "/Users/<username>/test.yaml"
  debug              = true
}