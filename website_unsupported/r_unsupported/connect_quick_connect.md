---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_quick_connect"
description: |-
  Provides details about a specific Amazon Quick Connect
---

# Resource: aws_connect_quick_connect

Provides an Amazon Connect Quick Connect resource. For more information see
[Amazon Connect: Getting Started](https://docs.aws.amazon.com/connect/latest/adminguide/amazon-connect-get-started.html)

## Example Usage

```terraform
resource "aws_connect_quick_connect" "test" {
  instance_id = "aaaaaaaa-bbbb-cccc-dddd-111111111111"
  name        = "Example Name"
  description = "quick connect phone number"

  quick_connect_config {
    quick_connect_type = "PHONE_NUMBER"

    phone_config {
      phone_number = "+12345678912"
    }
  }

  tags = {
    "Name" = "Example Quick Connect"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Specifies the description of the Quick Connect.
* `instance_id` - (Required) Specifies the identifier of the hosting Amazon Connect Instance.
* `name` - (Required) Specifies the name of the Quick Connect.
* `quick_connect_config` - (Required) A block that defines the configuration information for the Quick Connect: `quick_connect_type` and one of `phone_config`, `queue_config`, `user_config` . The Quick Connect Config block is documented below.
* `tags` - (Optional) Tags to apply to the Quick Connect. If configured with a provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

A `quick_connect_config` block supports the following arguments:

* `quick_connect_type` - (Required) Specifies the configuration type of the quick connect. valid values are `PHONE_NUMBER`, `QUEUE`, `USER`.
* `phone_config` - (Optional) Specifies the phone configuration of the Quick Connect. This is required only if `quick_connect_type` is `PHONE_NUMBER`. The `phone_config` block is documented below.
* `queue_config` - (Optional) Specifies the queue configuration of the Quick Connect. This is required only if `quick_connect_type` is `QUEUE`. The `queue_config` block is documented below.
* `user_config` - (Optional) Specifies the user configuration of the Quick Connect. This is required only if `quick_connect_type` is `USER`. The `user_config` block is documented below.

A `phone_config` block supports the following arguments:

* `phone_number` - (Required) Specifies the phone number in in E.164 format.

A `queue_config` block supports the following arguments:

* `contact_flow_id` - (Required) Specifies the identifier of the contact flow.
* `queue_id` - (Required) Specifies the identifier for the queue.

A `user_config` block supports the following arguments:

* `contact_flow_id` - (Required) Specifies the identifier of the contact flow.
* `user_id` - (Required) Specifies the identifier for the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the Quick Connect.
* `quick_connect_id` - The identifier for the Quick Connect.
* `id` - The identifier of the hosting Amazon Connect Instance and identifier of the Quick Connect separated by a colon (`:`).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block).

## Import

Amazon Connect Quick Connects can be imported using the `instance_id` and `quick_connect_id` separated by a colon (`:`), e.g.,

```
$ terraform import aws_connect_quick_connect.example f1288a1f-6193-445a-b47e-af739b2:c1d4e5f6-1b3c-1b3c-1b3c-c1d4e5f6c1d4e5
```