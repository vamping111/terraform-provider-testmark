---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_vpcs"
description: |-
    Provides a list of VPC Ids in a region
---

# Data Source: aws_vpcs

This resource can be useful for getting back a list of VPC Ids for a region.

The following example retrieves a list of VPC IDs with a custom tag of `service` set to a value of "production".

## Example Usage

The following shows outputing all VPC IDs.

```terraform
data "aws_vpcs" "foo" {
  tags = {
    service = "production"
  }
}

output "foo" {
  value = data.aws_vpcs.foo.ids
}
```

## Argument Reference

* `tags` - (Optional) A map of tags, each pair of which must exactly match
  a pair on the desired VPCs.
* `filter` - (Optional) Custom filter block as described below.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field.
  A VPC will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-vpcs].

## Attributes Reference

* `id` - The region (e.g., `region-1`).
* `ids` - A list of all the VPC IDs found.

[describe-vpcs]: https://docs.cloud.croc.ru/en/api/ec2/vpcs/DescribeVpcs.html
