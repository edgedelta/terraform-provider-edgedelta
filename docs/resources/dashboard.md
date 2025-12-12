# edgedelta_dashboard Resource

Manages an EdgeDelta dashboard.

## Example Usage

### Basic Dashboard

```hcl
resource "edgedelta_dashboard" "monitoring" {
  dashboard_name = "Infrastructure Monitoring"
  description    = "Main monitoring dashboard for production infrastructure"
  tags           = ["infrastructure", "monitoring", "production"]

  definition = file("${path.module}/dashboards/monitoring.json")
}
```

### Dashboard with Inline Definition

```hcl
resource "edgedelta_dashboard" "simple" {
  dashboard_name = "Simple Dashboard"

  definition = jsonencode({
    panels = [
      {
        type  = "graph"
        title = "Requests per Second"
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

### Required

* `dashboard_name` - (Required) Name of the dashboard.

### Optional

* `description` - (Optional) Description of the dashboard.
* `tags` - (Optional) List of searchable tags for the dashboard.
* `definition` - (Optional) Dashboard definition as a JSON string. Use `file()` to load from a file or `jsonencode()` for inline definitions. The provider will suppress diffs for semantically equivalent JSON.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `dashboard_id` - Unique identifier for the dashboard.
* `creator` - User ID who created the dashboard.
* `updater` - User ID who last updated the dashboard.
* `created` - UTC timestamp of dashboard creation.
* `updated` - UTC timestamp of last update.

## Import

Dashboards can be imported using the dashboard ID:

```shell
terraform import edgedelta_dashboard.example <dashboard_id>
```

Multiple dashboards can be imported using comma-separated IDs:

```shell
terraform import edgedelta_dashboard.example <id1>,<id2>,<id3>
```

All dashboards can be imported using `*`:

```shell
terraform import edgedelta_dashboard.example "*"
```
