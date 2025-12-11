![](logo.png)

Terraform Provider Edge Delta
==================

Edge Delta provider requires Terraform 0.13.0 and later.

* [Terraform Website](https://www.terraform.io)
* [Edge Delta Provider Documentation](docs/index.md)
* [Edge Delta Provider Usage Examples](examples/)

## Usage Example

> When using the Edge Delta Provider, the recommended approach to pass the `api_secret` is to use an environment variable instead of explicitly passing the secret in the provider block.

```hcl
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
  org_id             = "<your-organization-id>"
  api_secret         = var.ED_API_TOKEN
}

resource "edgedelta_config" "bare_minimum" {
  config_content     = file("/path/to/the/agent/configuration/file.yml")
}

# Dashboard example
resource "edgedelta_dashboard" "monitoring" {
  dashboard_name = "Infrastructure Monitoring"
  description    = "Main monitoring dashboard"
  tags           = ["infrastructure", "monitoring"]

  definition = file("/path/to/dashboard.json")
}
```

## Available Resources

| Resource | Description |
|----------|-------------|
| `edgedelta_config` | Manages agent configurations |
| `edgedelta_dashboard` | Manages dashboards |

Further [usage documentation is available in the provider repo](docs/index.md).

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.13.0+
* [Go](https://golang.org/doc/install) version 1.25.3 (to build the provider plugin)

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.25.3 is required).

First clone the repository to: `$GOPATH/src/edgedelta.com/edgedelta/terraform-provider-edgedelta`

```bash
repo_path="$GOPATH/src/edgedelta.com/edgedelta"
mkdir -p $repo_path
cd $repo_path
$ git clone git@github.com:edgedelta/terraform-provider-edgedelta
$ cd terraform-provider-edgedelta
```

Once inside the provider directory, you can compile the provider by running `make`, which will build the provider and put the provider binary in the `~/.terraform.d/plugins` directory.
