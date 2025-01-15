---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface"
description: |-
  Get information on a Network Interface resource.
---

# aws_network_interface

Use this data source to get information about a network interface.

## Example Usage

```terraform
data "aws_network_interface" "example" {
  id = "eni-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `id` – (Optional) The identifier for the network interface.
* `filter` – (Optional) One or more name/value pairs to filter.

For more information about filtering, see the [EC2 API documentation][describe-network-interfaces].

## Attributes Reference

See the [`aws_network_interface`][tf-network-interface] for details on the returned attributes.

Additionally, the following attributes are exported:

* `arn` - The ARN of the network interface.
* `association` - The association information for an Elastic IP address (IPv4) associated with the network interface. See supported fields below.
* `availability_zone` - The availability zone.
* `description` - Description of the network interface.
* `mac_address` - The MAC address.
* `owner_id` - The project ID.
* `private_dns_name` - The private DNS name.
* `private_ip` - The private IPv4 address of the network interface within the subnet.
* `private_ips` - The private IPv4 addresses associated with the network interface.
* `security_groups` - The list of security groups for the network interface.
* `subnet_id` - The ID of the subnet.
* `tags` - Any tags assigned to the network interface.
* `vpc_id` - The ID of the VPC.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `interface_type` - The type of interface. Always `"interface"`.
* `ipv6_addresses` - List of IPv6 addresses to assign to the ENI. Always empty.
* `requester_id` - The ID of the entity that launched the instance on your behalf. Always `""`.
* `outpost_arn` - The ARN of the Outpost. Always `""`.

### `association`

* `allocation_id` - The allocation ID.
* `association_id` - The association ID.
* `customer_owned_ip` - The customer-owned IP address.
* `ip_owner_id` - The ID of the elastic IP address owner.
* `public_dns_name` - The public DNS name.
* `public_ip` - The address of the elastic IP address bound to the network interface.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `carrier_ip` - The carrier IP address associated with the network interface. This attribute is only set when the network interface is in a subnet which is associated with a Wavelength Zone.

## Import

Elastic network interfaces can be imported using the `id`, e.g.,

```
$ terraform import aws_network_interface.test eni-12345678
```

[describe-network-interfaces]: https://docs.cloud.croc.ru/en/api/ec2/network_interfaces/DescribeNetworkInterfaces.html
[tf-network-interface]: network_interface.html
