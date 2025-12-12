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

output "dashboard_id" {
  value       = edgedelta_dashboard.existing.dashboard_id
  description = "The dashboard ID of the edgedelta_dashboard instance"
}

output "dashboard_name" {
  value       = edgedelta_dashboard.existing.dashboard_name
  description = "The name of the edgedelta_dashboard instance"
}

# Import an existing dashboard using: terraform import edgedelta_dashboard.existing <dashboard-id>
resource "edgedelta_dashboard" "existing" {
  dashboard_name = "My Existing Dashboard"
  description    = "An existing dashboard managed via Terraform"
  tags           = ["terraform", "imported"]

  definition = file("/path/to/the/dashboard/definition/file.json")
}
