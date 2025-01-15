---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group_rule"
description: |-
  Provides a security group rule resource.
---

# Resource: aws_security_group_rule

Provides a security group rule resource. Represents a single `ingress` or
`egress` group rule, which can be added to external security groups.

~> **Note on Security Groups and Security Group Rules:** Terraform currently
provides both a standalone security group rule resource (a single `ingress` or
`egress` rule), and a [`aws_security_group`][tf-security-group] with `ingress` and `egress` rules
defined in-line. At this time you cannot use a security group with in-line rules
in conjunction with any security group rule resources. Doing so will cause
a conflict of rule settings and will overwrite rules.

~> **Note** Setting `protocol = "all"` or `protocol = -1` with `from_port` and `to_port` will result in the EC2 API creating a security group rule with all ports open. This API behavior cannot be controlled by Terraform and may generate warnings in the future.

## Example Usage

Basic usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_security_group" "example" {
  name        = "test_security_group"
  description = "test_security_group"
  vpc_id      = aws_vpc.example.id
}

resource "aws_security_group_rule" "example" {
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = [aws_vpc.example.cidr_block]
  security_group_id = aws_security_group.example.id
}
```

## Argument Reference

The following arguments are required:

* `from_port` - (Required) Start port (or ICMP type number if protocol is "icmp" or "icmpv6").
* `protocol` - (Required) Protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
* `security_group_id` - (Required) Security group to apply this rule to.
* `to_port` - (Required) End port (or ICMP code if protocol is "icmp").
* `type` - (Required) Type of rule being created. Valid options are `ingress` (inbound)
or `egress` (outbound).

The following arguments are optional:

* `cidr_blocks` - (Optional) List of CIDR blocks. Cannot be specified with `source_security_group_id` or `self`.
* `description` - (Optional) Description of the rule.
* `ipv6_cidr_blocks` - (Optional) List of IPv6 CIDR blocks. Cannot be specified with `source_security_group_id` or `self`.
* `self` - (Optional) Whether the security group itself will be added as a source to this ingress rule. Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `source_security_group_id`.
* `source_security_group_id` - (Optional) Security group id to allow access to/from, depending on the `type`. Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `self`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the security group rule.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `prefix_list_ids` - List of prefix list IDs (for allowing access to VPC endpoints). Always empty.

## Import

-> **Unsupported operation**
Import security group rules is currently unsupported.

Security group rules can be imported using the `security_group_id`, `type`, `protocol`, `from_port`, `to_port`, and source(s)/destination(s) (e.g., `cidr_block`) separated by underscores (`_`). All parts are required.

Not all rule permissions (e.g., not all of a rule's CIDR blocks) need to be imported for Terraform to manage rule permissions. However, importing some of a rule's permissions but not others, and then making changes to the rule will result in the creation of an additional rule to capture the updated permissions. Rule permissions that were not imported are left intact in the original rule.

Import an ingress rule in security group `sg-12345678` for TCP port 8000 with an IPv4 destination CIDR of `10.0.3.0/24`:

```console
$ terraform import aws_security_group_rule.ingress sg-12345678_ingress_tcp_8000_8000_10.0.3.0/24
```

Import a rule with various IPv4 and IPv6 source CIDR blocks:

```console
$ terraform import aws_security_group_rule.ingress sg-12345678_ingress_tcp_100_121_10.1.0.0/16_2001:db8::/48_10.2.0.0/16_2002:db8::/48
```

Import a rule, applicable to all ports, with a protocol other than TCP/UDP/ICMP/ICMPV6/ALL, e.g., Multicast Transport Protocol (MTP), using the IANA protocol number, e.g., 92.

```console
$ terraform import aws_security_group_rule.ingress sg-12345678_ingress_92_0_65536_10.0.3.0/24_10.0.4.0/24
```

Import a default any/any egress rule to 0.0.0.0/0:

```console
$ terraform import aws_security_group_rule.default_egress sg-12345678_egress_all_0_0_0.0.0.0/0
```

Import a rule applicable to all protocols and ports with a security group source:

```console
$ terraform import aws_security_group_rule.ingress_rule sg-12345678_ingress_all_0_65536_sg-61766572
```

Import a rule that has itself and an IPv6 CIDR block as sources:

```console
$ terraform import aws_security_group_rule.rule_name sg-12345678_ingress_tcp_80_80_self_2001:db8::/48
```

[tf-security-group]: security_group.html
