---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_connection"
description: |-
  Provides information about a Direct Connect connection.
---

# Data Source: aws_dx_connection

Provides information about a Direct Connect connection.

## Example Usage

```terraform
data "aws_dx_connection" "selected" {
  name = "tf-dxconn-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the connection.
* `aws_device` - The ID of the device to which the connection is attached.
* `bandwidth` - The bandwidth of the connection.
* `id` - The ID of the connection.
* `location` - The physical site, where the connection terminates.
* `owner_account_id` - The ID of the project that owns the connection.
* `tags` - Tags assigned to the connection.

->  **Unsupported attributes**
These attributes are currently unsupported:

* `provider_name` - The name of the service provider associated with the connection. Always `""`.
