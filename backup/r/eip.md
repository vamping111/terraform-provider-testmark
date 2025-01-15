---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip"
description: |-
  Provides an Elastic IP resource.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[elastic-ips]: https://docs.cloud.croc.ru/en/services/networks/addresses/operations.html
[timeouts]: https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts

# Resource: aws_eip

Provides an Elastic IP resource.

For more information about EIPs, see [user documentation][elastic-ips].

## Example Usage

### Single EIP associated with an instance

```terraform
resource "aws_eip" "example" {
  instance = "i-12345678"
  vpc      = true
}
```

### Attaching an EIP to an instance with a pre-assigned private ip

```terraform
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "tf_test_subnet" {
  vpc_id     = aws_vpc.default.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_instance" "foo" {
  ami           = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"

  private_ip = "10.0.0.12"
  subnet_id  = aws_subnet.tf_test_subnet.id
}

resource "aws_eip" "bar" {
  vpc = true

  associate_with_private_ip = "10.0.0.12"
}
```

### Allocating EIP from the BYOIP pool

```terraform
resource "aws_eip" "byoip-ip" {
  vpc              = true
  public_ipv4_pool = "ipv4pool-ec2-012345"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Optional) IP address from an EC2 BYOIP pool. This option is only available for VPC EIPs.
* `associate_with_private_ip` - (Optional) User-specified primary or secondary private IP address to associate with the elastic IP address.
  If no private IP address is specified, the elastic IP address is associated with the primary private IP address.
* `instance` - (Optional) EC2 instance ID.
* `network_interface` - (Optional) Network interface ID to associate with.
* `public_ipv4_pool` - (Optional) EC2 IPv4 address pool identifier. This option is only available for VPC EIPs.
* `tags` - (Optional) Map of tags to assign to the resource. Tags can only be applied to EIPs in a VPC. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `vpc` - (Optional) Boolean if the EIP is in a VPC or not.

~> **Note** You can specify either the `instance` ID or the `network_interface` ID, but not both.

~> **Note** Specifying both `public_ipv4_pool` and `address` won't cause an error but `address` will be used in the
case both options are defined as the api only requires one or the other.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - ID that the cloud assigns to represent the allocation of the elastic IP address for use with instances in a VPC.
* `association_id` - ID representing the association of the address with an instance in a VPC.
* `domain` - Indicates if this EIP is for use in VPC (`vpc`).
* `id` - Contains the EIP allocation ID.
* `private_ip` - Contains the private IP address. Can be `""` if `associate_with_private_ip` is specified.
* `public_ip` - Contains the public IP address.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `carrier_ip` - Carrier IP address. Always `""`.
* `customer_owned_ip` - Customer owned IP. Always `""`.
* `customer_owned_ipv4_pool` - (Optional) ID  of a customer-owned address pool. Always `""`.
* `network_border_group` - (Optional) Location from which the IP address is advertised. Always `""`.
* `private_dns` - The Private DNS associated with the Elastic IP address (if in VPC). Computed by provider.
* `public_dns` - Public DNS associated with the Elastic IP address. Computed by provider.

~> **Note** The data source computes the `public_dns` and `private_dns` attributes according to the [AWS VPC DNS Guide](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html#vpc-dns-hostnames) as they are not available with the EC2 API.

## Timeouts

`aws_eip` provides the following [timeouts] configuration options:

- `read` - (Default `15 minutes`) How long to wait querying for information about EIPs.
- `update` - (Default `5 minutes`) How long to wait for an EIP to be updated.
- `delete` - (Default `3 minutes`) How long to wait for an EIP to be deleted.

## Import

EIPs in a VPC can be imported using their allocation ID, e.g.,

```
$ terraform import aws_eip.bar eipalloc-1234567
```

EIPs can be imported using their public IP, e.g.,

```
$ terraform import aws_eip.bar 1.1.1.1
```
