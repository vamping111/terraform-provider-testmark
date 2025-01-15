---
subcategory: "Direct Connect"
layout: "aws"
page_title: "aws_dx_transit_virtual_interface"
description: |-
  Manages a Direct Connect transit virtual interface.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_dx_transit_virtual_interface

Manages a Direct Connect transit virtual interface.
A transit virtual interface is a VLAN that transports traffic from a [Direct Connect gateway](dx_gateway.html) to one or more [transit gateways](ec2_transit_gateway.html).

## Example Usage

```terraform
data "aws_dx_connection" "selected" {
  name = "tf-dxconn-example"
}

resource "aws_dx_gateway" "example" {
  name            = "tf-dxgw-example"
  amazon_side_asn = "64512"
}

resource "aws_dx_transit_virtual_interface" "example" {
  name          = "tf-dxvif-example"
  connection_id = data.aws_dx_connection.selected.id
  dx_gateway_id = aws_dx_gateway.example.id
  vlan          = "4094"
  bgp_asn       = "65352"
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Optional) The address family for the BGP peer. Valid value is `ipv4 `. Defaults to `ipv4`.
* `amazon_address` - (Optional) The IPv4 CIDR address of the connection endpoint on the cloud side.
* `bgp_asn` - (Required) The BGP ASN on the client side.
* `bgp_auth_key` - (Optional, Sensitive) The authentication key for BGP configuration.
* `connection_id` - (Required) The ID of the Direct Connect connection (or LAG) on which the virtual interface has to be created.
* `customer_address` - (Optional) The IPv4 CIDR address of the connection endpoint on the client side.
* `dx_gateway_id` - (Required) The ID of the Direct Connect gateway that the virtual interface must be connected to.
* `name` - (Required) The name for the virtual interface.
* `tags` - (Optional, Editable) Tags assigned to the resource. If there is a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider level.
* `vlan` - (Required) The VLAN ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `amazon_side_asn` - The ASN for the cloud side of the connection.
* `arn` - The ARN of the virtual interface.
* `aws_device` - The ID of the device to which the connection (or LAG) is attached.
* `id` - The ID of the virtual interface.
* `tags_all` - Tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

->  **Unsupported attributes**
These attributes are currently unsupported:

* `jumbo_frame_capable` - Indicates whether jumbo frames (8500 MTU) are supported. Always `false`.
* `mtu` - The maximum transmission unit (MTU) is the size, in bytes, of the largest permissible packet that can be passed over the connection. Always `0`.
* `sitelink_enabled` - Indicates whether to enable or disable SiteLink.

## Timeouts

`aws_dx_transit_virtual_interface` provides the following [Timeouts][timeouts] configuration options:

- `create` - (Default `10 minutes`) Timeout for creating virtual interface
- `update` - (Default `10 minutes`) Timeout for virtual interface modifications
- `delete` - (Default `10 minutes`) Timeout for destroying virtual interface

## Import

Direct Connect transit virtual interfaces can be imported using `id`, e.g.,

```
$ terraform import aws_dx_transit_virtual_interface.example dxvif-12345678
```
