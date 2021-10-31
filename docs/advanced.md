# Provider Eng Documentation

This document provides detailed information about the advanced usage of Edge Delta provider.

## Imports

Terraform CLI offers functionality to import the existing resources to your tfstate and resource definition file, to prevent deleting and re-creating the remote resource to be able to define and manage the same resource on local.

Currently `edgedelta_config` and `edgedelta_monitor` resources support importing. To import your resources, please follow the instructions below:

#### Get resource ID

1. Go to [app.edgedelta.com](https://app.edgedelta.com) and log in to your account
2. Go to the page of the specific resource that you want to import
   1. For `edgedelta_config`, go to Agent Settings under Data Pipeline
   2. For `edgedelta_monitor`, go to Monitors under Management
3. Find the ID of the specific resource you want to import
   1. For `edgedelta_config`, the IDs are listed under the "Key" column
   2. For `edgedelta_monitor`, open the developer tools and go to network tab for Chrome (or the equivalent one for any browser)
   3. Refresh the page and find the xhr request with name `alert_definitions`
   4. Inspect the response and find the `id` of the specific monitor you want to import

#### Create the resource definition

1. Create a skeleton resource definition in your `.tf` file. An example resoruce skeleton can be found below:

```hcl
resource "edgedelta_config" "imported_config" {

}
```

2. Run `terraform import` with the name of the skeleton resource you have just created and the resource ID you have from the previous steps:

```bash
terraform import edgedelta_config.imported_config <resource-id>
```

3. If you have done the previous steps correctly, terraform will fetch the resource data from Edge Delta and store it in the `.tfstate` file. You should be able to see the resource data in your terraform state now. Feel free to check your state file to make sure `terraform import` ran correctly.

#### Sync the resource definition

1. Run `terraform show` in your terminal. This command will show the data of your resources in the current state file in the hcl format.
2. Copy the resource definition you have imported recently from the output of the previous command, and use it to fill up the resource skeleton.
3. Run `terraform apply` to see that there is no diff between the resource in the state and the one in the `.tf` file.