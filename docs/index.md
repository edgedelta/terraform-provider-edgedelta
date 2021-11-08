---
page_title: "Edge Delta Provider"
subcategory: ""
description: "Terraform provider for managing Edge Delta configurations"
  
---

# Edge Delta Provider

## Overview

You can use this document to learn how to configure a Terraform provider to manage your Edge Delta configurations. 

***

## Step 1: Review Requirements

Review the following software requirements: 

* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.0
* [terraform-provider-edgedelta](https://github.com/edgedelta/terraform-provider-edgedelta) plugin >= 0.0.1

For additional agent-related requirements, please see [Pre-Installation Agent Requirements](https://docs.edgedelta.com/agent-requirements/). 

***

## Step 2: Review Schema Elements

Review the following schema elements that you can obtain from Edge Delta.

> **Note**
> 
> In this document, the steps below describe how to obtain these elements.


| Elements          | Description                                                                                                                                                  | Type               | Default                   | Required? |
|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|---------------------------|-----------|
| api_secret   | This key is the API token for your configuration file. Edge Delta recommends that you use Terraform variables to pass the token value in a resource schema.  | String (sensitive) | Not applicable            | Yes       |
| org_id       | This key is the unique organization ID for your Edge Delta account.                                                                                          | String             | Not applicable            | Yes       |
| api_endpoint | This key is the API base URL for Edge Delta.                                                                                                                 | String             | https://api.edgedelta.com | No        |

***

## Step 2: Create an Edge Delta Account

The plugin uses token authentication to access the API. As a result, to use this plugin, you must have an Edge Delta account and an API key.

After you obtain your API key, you must set the `api_secret` parameter of the provider in your `.tf` files with your API key.

Edge Delta recommends that you set the `api_secret` through an environment variable. In the sample configurations, `TF_VAR_ED_API_TOKEN` is used as the environment variable. You can see the example usage in the `examples/` folder in the project root.

1. Navigate to [admin.edgedelta.com](https://admin.edgedelta.com/), and then click **Sign up**.
2. Complete the missing fields, and then click **Register**. You will be redirected to the **Welcome to Edge Delta** screen in the Edge Delta Admin portal.
3. Click **Exit Set Up**.
4. In the pop-up window, click **X** not to access the demo environment. 
5. In the left-side navigation, click **My Organization**. 
6. In the table, under **Actions**, click the **Edit User Permissions** icon. 
7. Copy the organization ID. You will need this information in a later step. 

***

## Step 3: Create a Configuration

1. In the Edge Delta Admin portal, in the left-side navigation, click **Agent Settings**. 
2. Click **Create Configurations**. 
3. There are 2 ways to create a configuration. You can:
  * Use a template with default parameters. 
  * Use a visual editor to populate a YAML file. 
4. Click **Save**. 
  * You will be redirected to the **Agent Settings** page. 
  * Refresh the page until your newly created configuration appear. 
6. Under **Key**, copy the API key for your newly created configuration. 

***

## Step 4: Create a Resource Definition 

Create a resource definition for **edgedelta_config** or **edgedelta_monitor** in the **.tf** file.
  * The **edgedelta_config** resource allows you to manage your Edge Delta configurations.
  * The **edgedelta_monitor** resource allows you to manage your Edge Delta alert definitons (monitors). 

***

### Option 1: Create a Resource Definition for edgedelta_config

To learn how to create a resource definition for **edgedelta_config**, review the following example: 






  * To learn more, see [edgedelta_config (Resource)](resources/config.md).

***

### Option 2: Create a Resource Definition for edgedelta_monitor

To learn how to create a resource definition for **edgedelta_monitor**, review the following example: 







* To learn more, see [edgedelta_monitor (Resource)](resources/monitor.md).

***

## Review Sample Usage 

Review the following sample usage: 

> **Note**
> 
> There are others examples included in the `examples/` folder.


```bash
export TF_VAR_ED_API_TOKEN="<your-api-token-goes-here>"
```

```hcl
variable "ED_API_TOKEN" {
  type = string
}

provider "edgedelta" {
  org_id             = "22222222-2222-2222-2222-222222222222"
  api_secret         = var.ED_API_TOKEN
}
```


***
