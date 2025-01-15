---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_connection"
description: |-
  Manages a Site-to-Site VPN connection. A Site-to-Site VPN connection is an Internet Protocol security (IPSec) VPN connection between a VPC and an on-premises network.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[sensitive-data]: https://www.terraform.io/docs/state/sensitive-data.html
[vpn-connections]: https://docs.cloud.croc.ru/en/services/networks/vpn_connections/operations.html

# Resource: aws_vpn_connection

-> **Unsupported resource**
This resource is currently unsupported.

Manages a Site-to-Site VPN connection. A Site-to-Site VPN connection is an Internet Protocol security (IPSec) VPN connection between a VPC and an on-premises network.

For more information about VPN connections, see [user documentation][vpn-connections].

~> **Note:** All arguments including `tunnel1_preshared_key` and `tunnel2_preshared_key` will be stored in the raw state as plain-text.
[Read more about sensitive data in state][sensitive-data].

~> **Note:** The CIDR blocks in the arguments `tunnel1_inside_cidr` and `tunnel2_inside_cidr` must have a prefix of /30 and be a part of a specific range.

-> The terms VPC, Internet Gateway, VPN Gateway are equivalent

## Example Usage

```terraform
resource "aws_vpc" "example" {
   cidr_block         = "172.16.8.0/24"
   enable_dns_support = true
   
   tags = {
     Name = "tf-vpc"
   }
}

resource "aws_customer_gateway" "example" {
  bgp_asn    = 65000
  ip_address = "172.0.0.1"
  type       = "ipsec.1"
  
  tags = {
    Name = "tf-customer-gateway"
  }
}

resource "aws_vpn_connection" "example" {
  vpn_gateway_id      = replace(aws_vpc.example.id, "/vpc/", "vgw")
  customer_gateway_id = aws_customer_gateway.example.id
  type                = aws_customer_gateway.example.type
  
  tags = {
    Name = "tf-vpn-connection"
  }
}
```

## Argument Reference

The following arguments are required:

* `customer_gateway_id` - (Required) The ID of the customer gateway.
* `type` - (Required) The type of VPN connection. Valid values are `ipsec.1`, `ipsec.legacy`.
* `vpn_gateway_id` - (Required) The ID of the VPN gateway.

Other arguments:

* `local_ipv4_network_cidr` - (Optional, Default `0.0.0.0/0`) The IPv4 CIDR on the customer gateway (on-premises) side of the VPN connection. Valid value must not fall within the range of 169.254.0.0/16.
* `remote_ipv4_network_cidr` - (Optional, Default `0.0.0.0/0`) The IPv4 CIDR on the cloud side of the VPN connection. Valid value must not fall within the range of 169.254.0.0/16.
* `tags` - (Optional) Tags to apply to the connection. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `tunnel1_ike_versions` - (Optional) The IKE versions that are permitted for the first VPN tunnel. Valid values are `ikev1 | ikev2`.
* `tunnel2_ike_versions` - (Optional) The IKE versions that are permitted for the second VPN tunnel. Valid values are `ikev1 | ikev2`.
* `tunnel1_inside_cidr` - (Optional) The CIDR block of the inside IP addresses for the first VPN tunnel. Valid value is a size /30 CIDR block from the 169.254.0.0/16 range.
* `tunnel2_inside_cidr` - (Optional) The CIDR block of the inside IP addresses for the second VPN tunnel. Valid value is a size /30 CIDR block from the 169.254.0.0/16 range.
* `tunnel1_preshared_key` - (Optional) The preshared key of the first VPN tunnel. The preshared key must be between 8 and 64 characters in length and cannot start with zero(0). Allowed characters are alphanumeric characters, periods(.) and underscores(_).
* `tunnel2_preshared_key` - (Optional) The preshared key of the second VPN tunnel. The preshared key must be between 8 and 64 characters in length and cannot start with zero(0). Allowed characters are alphanumeric characters, periods(.) and underscores(_).
* `tunnel1_phase1_dh_group_numbers` - (Optional) List of one or more Diffie-Hellman group numbers that are permitted for the first VPN tunnel for phase 1 IKE negotiations. Valid values are `2 | 5 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21`.
* `tunnel2_phase1_dh_group_numbers` - (Optional) List of one or more Diffie-Hellman group numbers that are permitted for the second VPN tunnel for phase 1 IKE negotiations. Valid values are `2 | 5 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21`.
* `tunnel1_phase1_encryption_algorithms` - (Optional) List of one or more encryption algorithms that are permitted for the first VPN tunnel for phase 1 IKE negotiations. Valid values are `AES128 | AES256 | AES128-GCM-16 | AES256-GCM-16`.
* `tunnel2_phase1_encryption_algorithms` - (Optional) List of one or more encryption algorithms that are permitted for the second VPN tunnel for phase 1 IKE negotiations. Valid values are `AES128 | AES256 | AES128-GCM-16 | AES256-GCM-16`.
* `tunnel1_phase1_integrity_algorithms` - (Optional) One or more integrity algorithms that are permitted for the first VPN tunnel for phase 1 IKE negotiations. Valid values are `SHA1 | SHA2-256 | SHA2-384 | SHA2-512`.
* `tunnel2_phase1_integrity_algorithms` - (Optional) One or more integrity algorithms that are permitted for the second VPN tunnel for phase 1 IKE negotiations. Valid values are `SHA1 | SHA2-256 | SHA2-384 | SHA2-512`.
* `tunnel1_phase1_lifetime_seconds` - (Optional, Default `28800`) The lifetime for phase 1 of the IKE negotiation for the first VPN tunnel, in seconds. Valid value is between `900` and `28800`.
* `tunnel2_phase1_lifetime_seconds` - (Optional, Default `28800`) The lifetime for phase 1 of the IKE negotiation for the second VPN tunnel, in seconds. Valid value is between `900` and `28800`.
* `tunnel1_phase2_dh_group_numbers` - (Optional) List of one or more Diffie-Hellman group numbers that are permitted for the first VPN tunnel for phase 2 IKE negotiations. Valid values are `0 | 2 | 5 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21`.
  The `0` value is applicable for Phase2 only and means that the PFS (Perfect Forward Secrecy) mode is disabled.
  In order to prevent the session encryption key from being compromised, do not disable PFS.
* `tunnel2_phase2_dh_group_numbers` - (Optional) List of one or more Diffie-Hellman group numbers that are permitted for the second VPN tunnel for phase 2 IKE negotiations. Valid values are `0 | 2 | 5 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21`.
  The `0` value is applicable for Phase2 only and means that the PFS (Perfect Forward Secrecy) mode is disabled.
  In order to prevent the session encryption key from being compromised, do not disable PFS.
* `tunnel1_phase2_encryption_algorithms` - (Optional) List of one or more encryption algorithms that are permitted for the first VPN tunnel for phase 2 IKE negotiations. Valid values are `AES128 | AES256 | AES128-GCM-16 | AES256-GCM-16`.
* `tunnel2_phase2_encryption_algorithms` - (Optional) List of one or more encryption algorithms that are permitted for the second VPN tunnel for phase 2 IKE negotiations. Valid values are `AES128 | AES256 | AES128-GCM-16 | AES256-GCM-16`.
* `tunnel1_phase2_integrity_algorithms` - (Optional) List of one or more integrity algorithms that are permitted for the first VPN tunnel for phase 2 IKE negotiations. Valid values are `SHA1 | SHA2-256 | SHA2-384 | SHA2-512`.
* `tunnel2_phase2_integrity_algorithms` - (Optional) List of one or more integrity algorithms that are permitted for the second VPN tunnel for phase 2 IKE negotiations. Valid values are `SHA1 | SHA2-256 | SHA2-384 | SHA2-512`.
* `tunnel1_phase2_lifetime_seconds` - (Optional, Default `3600`) The lifetime for phase 2 of the IKE negotiation for the first VPN tunnel, in seconds. Valid value is between `900` and `3600`.
* `tunnel2_phase2_lifetime_seconds` - (Optional, Default `3600`) The lifetime for phase 2 of the IKE negotiation for the second VPN tunnel, in seconds. Valid value is between `900` and `3600`.
* `tunnel1_replay_window_size` - (Optional, Default `1024`) The number of packets in an IKE replay window for the first VPN tunnel. Valid value is between `64` and `2048`.
* `tunnel2_replay_window_size` - (Optional, Default `1024`) The number of packets in an IKE replay window for the second VPN tunnel. Valid value is between `64` and `2048`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the VPN connection.
* `id` - The ID of the VPN connection.
* `customer_gateway_configuration` - The configuration information for the VPN connection's customer gateway (in the native XML format).
* `customer_gateway_id` - The ID of the customer gateway to which the connection is attached.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `tunnel1_address` - The public IP address of the first VPN tunnel.
* `tunnel1_cgw_inside_address` - The RFC 6890 link-local address of the first VPN tunnel (Customer Gateway Side).
* `tunnel1_vgw_inside_address` - The RFC 6890 link-local address of the first VPN tunnel (VPN Gateway Side).
* `tunnel1_preshared_key` - The preshared key of the first VPN tunnel.
* `tunnel1_bgp_asn` - The bgp asn number of the first VPN tunnel.
* `tunnel1_bgp_holdtime` - The bgp holdtime of the first VPN tunnel.
* `tunnel2_address` - The public IP address of the second VPN tunnel.
* `tunnel2_cgw_inside_address` - The RFC 6890 link-local address of the second VPN tunnel (Customer Gateway Side).
* `tunnel2_vgw_inside_address` - The RFC 6890 link-local address of the second VPN tunnel (VPN Gateway Side).
* `tunnel2_preshared_key` - The preshared key of the second VPN tunnel.
* `tunnel2_bgp_asn` - The bgp asn number of the second VPN tunnel.
* `tunnel2_bgp_holdtime` - The bgp holdtime of the second VPN tunnel.
* `vgw_telemetry` - Telemetry for the VPN tunnels. Detailed below.
* `vpn_gateway_id` - The ID of the virtual private gateway to which the connection is attached.

### vgw_telemetry

* `accepted_route_count` - The number of accepted routes.
* `last_status_change` - The date and time of the last change in status.
* `outside_ip_address` - The Internet-routable IP address of the virtual private gateway's outside interface.
* `status` - The status of the VPN tunnel.
* `status_message` - If an error occurs, a description of the error.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `core_network_arn` - The ARN of the core network. Always `""`.
* `core_network_attachment_arn` - The ARN of the core network attachment. Always `""`.
* `enable_acceleration` - Indicate whether to enable acceleration for the VPN connection. Always `false`.
* `local_ipv6_network_cidr` - The IPv6 CIDR on the customer gateway (on-premises) side of the VPN connection. Always `""`.
* `remote_ipv6_network_cidr` - The IPv6 CIDR on the customer gateway (on-premises) side of the VPN connection. Always `""`.
* `routes` - The static routes associated with the VPN connection. Always empty.
    * `destination_cidr_block` - The CIDR block associated with the local subnet of the customer data center.
    * `source` - Indicates how the routes were provided.
    * `state` - The current state of the static route.
* `static_routes_only` - Whether the VPN connection uses static routes exclusively. Static routes must be used for devices that don't support BGP. Always `false`.
* `transit_gateway_attachment_id` - When associated with an EC2 transit gateway (`transit_gateway_id` argument), the attachment ID. Always `""`.
* `transit_gateway_id` - The ID of the EC2 transit gateway. Always `""`.
* `tunnel_inside_ip_version` - Indicate whether the VPN tunnels process IPv4 or IPv6 traffic. Always `""`.
* `tunnel1_dpd_timeout_action` - The action to take after DPD timeout occurs for the first VPN tunnel. Always empty.
* `tunnel2_dpd_timeout_action` - The action to take after DPD timeout occurs for the second VPN tunnel. Always empty.
* `tunnel1_dpd_timeout_seconds` - The number of seconds after which a DPD timeout occurs for the first VPN tunnel. Always empty.
* `tunnel2_dpd_timeout_seconds` - The number of seconds after which a DPD timeout occurs for the second VPN tunnel. Always empty.
* `tunnel1_inside_ipv6_cidr` - The range of inside IPv6 addresses for the first VPN tunnel. Always empty.
* `tunnel2_inside_ipv6_cidr` - The range of inside IPv6 addresses for the second VPN tunnel. Always empty.
* `tunnel1_rekey_fuzz_percentage` - The percentage of the rekey window for the first VPN tunnel (determined by `tunnel1_rekey_margin_time_seconds`) during which the rekey time is randomly selected. Always empty.
* `tunnel2_rekey_fuzz_percentage` - The percentage of the rekey window for the second VPN tunnel (determined by `tunnel2_rekey_margin_time_seconds`) during which the rekey time is randomly selected. Always empty.
* `tunnel1_rekey_margin_time_seconds` - The margin time, in seconds, before the phase 2 lifetime expires, during which the AWS side of the first VPN connection performs an IKE rekey. Always empty.
* `tunnel2_rekey_margin_time_seconds` - The margin time, in seconds, before the phase 2 lifetime expires, during which the AWS side of the second VPN connection performs an IKE rekey. Always empty.
* `tunnel1_startup_action` - The action to take when the establishing the tunnel for the first VPN connection. By default, your customer gateway device must initiate the IKE negotiation and bring up the tunnel. Always empty.
* `tunnel2_startup_action` - The action to take when the establishing the tunnel for the second VPN connection. By default, your customer gateway device must initiate the IKE negotiation and bring up the tunnel. Always empty.
* `vgw_telemetry`:
    * `certificate_arn` - The ARN of the VPN tunnel endpoint certificate. Always `""`.

## Import

VPN connections can be imported using the ID of VPN connection, e.g.,

```
$ terraform import aws_vpn_connection.example vpn-12345678
```
