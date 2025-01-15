---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_lag"
description: |-
  Provides information about a Direct Connect link aggregation group (LAG).
---

# Data Source: aws_dx_lag

Provides information about a Direct Connect link aggregation group (LAG).

## Example Usage

```terraform
data "aws_dx_lag" "selected" {
  name = "tf-dxlag-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the LAG.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the LAG.
* `aws_device` - The ID of the device to which the LAG is attached.
* `bandwidth` - The bandwidth of each physical connection in the LAG.
* `id` - The ID of the LAG.
* `location` - The physical site, where the connection terminates.
* `owner_account_id` - The ID of the project that owns the LAG.
* `tags` - Tags assigned to the LAG.

->  **Unsupported attributes**
These attributes are currently unsupported:

* `provider_name` - The name of the service provider associated with the LAG. Always `""`.
