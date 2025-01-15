---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "aws_ebs_volume"
description: |-
  Provides an elastic block storage resource.
---

# Resource: aws_ebs_volume

Manages a single EBS volume.

## Example Usage

```terraform
resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 40

  tags = {
    Name = "HelloWorld"
  }
}
```

~> **Note** At least one of `size` or `snapshot_id` is required when specifying an EBS volume

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The AZ where the EBS volume will exist.
* `iops` - (Optional) The amount of IOPS to provision for the disk. Only valid for `type` of `io2`.
* `size` - (Optional) The size of the drive in GiB.
* `snapshot_id` (Optional) A snapshot to base the EBS volume on.
* `type` - (Optional) The type of EBS volume. Can be `st2`, `gp2` or `io2` (Default: `st2`).
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the volume.
* `id` - The volume ID (e.g., vol-12345678).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `throughput` - The throughput that the volume supports, in MiB/s.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `encrypted` - Whether the disk is encrypted. Always `false`.
* `kms_key_id` - The ARN for the KMS encryption key. Always `""`.
* `multi_attach_enabled` - Whether EBS Multi-Attach is enabled. Always `false`.
* `outpost_arn` - The Amazon Resource Name (ARN) of the Outpost. Always `""`.

## Timeouts

`aws_ebs_volume` provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `5 minutes`) Used for creating volumes. This includes the time required for the volume to become available
- `update` - (Default `5 minutes`) Used for `size`, `type`, or `iops` volume changes
- `delete` - (Default `5 minutes`) Used for destroying volumes

## Import

EBS Volumes can be imported using the `id`, e.g.,

```
$ terraform import aws_ebs_volume.id vol-12345678
```

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
