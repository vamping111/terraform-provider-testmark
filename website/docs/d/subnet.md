---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnet"
description: |-
    Provides details about a specific VPC subnet
---

# Data Source: aws_subnet

`aws_subnet` provides details about a specific VPC subnet.

This resource can prove useful when a module accepts a subnet ID as an input variable and needs to, for example, determine the ID of the VPC that the subnet belongs to.

## Example Usage

The following example shows how one might accept a subnet ID as a variable and use this data source to obtain the data necessary to create a security group that allows connections from hosts in that subnet.

```terraform
variable "subnet_id" {}

data "aws_subnet" "selected" {
  id = var.subnet_id
}

resource "aws_security_group" "subnet" {
  vpc_id = data.aws_subnet.selected.vpc_id

  ingress {
    cidr_blocks = [data.aws_subnet.selected.cidr_block]
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
  }
}
```

### Filter Example

If you want to match against tag `Name`, use:

```terraform
data "aws_subnet" "selected" {
  filter {
    name   = "tag:Name"
    values = ["yakdriver"]
  }
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available subnets in the current region. The given filters must match exactly one subnet whose data will be exported as attributes.

The following arguments are optional:

* `availability_zone` - (Optional) Availability zone where the subnet must reside.
* `default_for_az` - (Optional) Whether the desired subnet must be the default subnet for its associated availability zone.
* `filter` - (Optional) Configuration block. Detailed below.
* `id` - (Optional) ID of the specific subnet to retrieve.
* `state` - (Optional) State that the desired subnet must have.
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired subnet.
* `vpc_id` - (Optional) ID of the VPC that the desired subnet belongs to.

### filter

This block allows for complex filters. You can use one or more `filter` blocks.

The following arguments are required:

* `name` - (Required) The name of the field to filter by it.
* `values` - (Required) Set of values that are accepted for the given field. A subnet will be selected if any one of the given values matches.

For more information about filtering, see the [EC2 API documentation][describe-subnets].

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `map_public_ip_on_launch` - Indicates whether public IP addresses will be associated with instances created in this subnet. Addresses are associated only if there are available allocated Elastic IP addresses.

->  **Unsupported attributes**
These attributes are currently unsupported:

* `arn` - ARN of the subnet. Always `""`.
* `assign_ipv6_address_on_creation` - Whether an IPv6 address is assigned on creation. Always `false`.
* `availability_zone_id` - AZ ID of the subnet. Always `""`.
* `customer_owned_ipv4_pool` - Identifier of customer owned IPv4 address pool. Always `""`.
* `enable_dns64` - Indicates whether DNS queries made to the Amazon-provided DNS Resolver in this subnet return synthetic IPv6 addresses for IPv4-only destinations. Always `false`.
* `enable_resource_name_dns_aaaa_record_on_launch` - Indicates whether to respond to DNS queries for instance hostnames with DNS AAAA records. Always `false`.
* `enable_resource_name_dns_a_record_on_launch` - Indicates whether to respond to DNS queries for instance hostnames with DNS A records. Always `false`.
* `ipv6_cidr_block_association_id` - Association ID of the IPv6 CIDR block. Always `""`.
* `ipv6_native` - Indicates whether this is an IPv6-only subnet. Always `false`.
* `map_customer_owned_ip_on_launch` - Whether customer owned IP addresses are assigned on network interface creation. Always `false`.
* `outpost_arn` - ARN of the Outpost. Always `""`.
* `owner_id` - The ID of the project that owns the subnet. Always `""`.
* `private_dns_hostname_type_on_launch` - The type of hostnames assigned to instances in the subnet at launch. Always `""`.

[describe-subnets]: https://docs.cloud.croc.ru/en/api/ec2/subnets/DescribeSubnets.html
