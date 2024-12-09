---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_gateway_route_propagation"
description: |-
  Requests automatic route propagation between a VPN gateway and a route table.
---

# Resource: aws_vpn_gateway_route_propagation

Requests automatic route propagation between a VPN gateway and a route table.

~> **Note** This resource should not be used with a route table that has
the `propagating_vgws` argument set. If that argument is set, any route
propagation not explicitly listed in its value will be removed.

## Example Usage

-> The terms VPC, Internet Gateway, VPN Gateway are equivalent

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
}

data "aws_vpn_gateway" "selected" {
  id = aws_vpc.example.id # vpc_id can be used as vpn_gateway_id
}

resource "aws_vpn_gateway_route_propagation" "example" {
  vpn_gateway_id = data.aws_vpn_gateway.selected.id
  route_table_id = aws_route_table.example.id
}
```

## Argument Reference

The following arguments are required:

* `vpn_gateway_id` - ID of the VPN gateway to propagate routes from.
* `route_table_id` - ID of the route table to propagate routes into.

## Attributes Reference

No additional attributes are exported.

## Timeouts

`aws_vpn_gateway_route_propagation` provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `2 minutes`) Used for propagation creation.
- `delete` - (Default `2 minutes`) Used for propagation deletion.
