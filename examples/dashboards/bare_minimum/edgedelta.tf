terraform {
  required_providers {
    edgedelta = {
      source  = "edgedelta/edgedelta"
      version = "0.0.9"
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

output "dashboard_id" {
  value       = edgedelta_dashboard.bare_minimum.dashboard_id
  description = "The dashboard ID of the edgedelta_dashboard instance"
}

output "dashboard_name" {
  value       = edgedelta_dashboard.bare_minimum.dashboard_name
  description = "The name of the edgedelta_dashboard instance"
}

resource "edgedelta_dashboard" "bare_minimum" {
  dashboard_name = "My Dashboard"
  description    = "A simple dashboard created via Terraform"
  tags           = ["terraform", "example"]

  definition = file("/path/to/the/dashboard/definition/file.json")
}
