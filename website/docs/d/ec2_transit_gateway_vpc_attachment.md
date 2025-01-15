---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "aws_ec2_transit_gateway_vpc_attachment"
description: |-
  Provides information about a transit gateway VPC attachment.
---

[describe-tgw-vpc-attachments]: https://docs.cloud.croc.ru/en/api/ec2/transit_gateways/DescribeTransitGatewayVpcAttachments.html

# Data Source: aws_ec2_transit_gateway_vpc_attachment

Provides information about a transit gateway VPC attachment.

## Example Usage

### By Filter

```terraform
data "aws_ec2_transit_gateway_vpc_attachment" "selected" {
  filter {
    name   = "vpc-id"
    values = ["vpc-12345678"]
  }
}
```

### By Identifier

```terraform
data "aws_ec2_transit_gateway_vpc_attachment" "selected" {
  id = "tgw-attach-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more configuration blocks containing name-values filters.
  The structure of this block is [described below](#filter).
* `id` - (Optional) The ID of the transit gateway VPC attachment.

### filter

* `name` - (Required) The name of the field to filter by it.
  Valid values can be found in the [EC2 API documentation][describe-tgw-vpc-attachments].
* `values` - (Required) List of one or more values for the filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the transit gateway attachment.
* `subnet_ids` - List of subnet IDs.
* `transit_gateway_id` - The ID of the transit gateway.
* `tags` - Map of tags assigned to the transit gateway VPC attachment.
* `vpc_id` - The ID of the VPC.
* `vpc_owner_id` - The ID of the project that owns the VPC.

->  **Unsupported attributes**
These attributes are currently unsupported:

* `appliance_mode_support` - Whether Appliance Mode support is enabled. Always empty.
* `dns_support` - Whether DNS support is enabled. Always empty.
* `ipv6_support` - Whether IPv6 support is enabled. Always empty.
