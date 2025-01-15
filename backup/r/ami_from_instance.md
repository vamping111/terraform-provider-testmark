---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ami_from_instance"
description: |-
  Creates an Amazon Machine Image (AMI) from an EBS-backed EC2 instance
---

# Resource: aws_ami_from_instance

Resource allows the creation of image from existing EBS-backed EC2 instance.

The created image will refer to implicitly-created snapshots of the instance's
EBS volumes and mimick its assigned block device configuration at the time
the resource is created.

This resource is best applied to an instance that is stopped when this instance
is created, so that the contents of the created image are predictable. When
applied to an instance that is running, *the instance will be stopped before taking
the snapshots and then started back up again*, resulting in a period of
downtime.

Note that the source instance is inspected only at the initial creation of this
resource. Ongoing updates to the referenced instance will not be propagated into
the generated image. Users may taint or otherwise recreate the resource in order
to produce a fresh snapshot.

## Example Usage

```terraform
resource "aws_ami_from_instance" "example" {
  name               = "terraform-example"
  source_instance_id = "i-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A region-unique name for the image.
* `source_instance_id` - (Required) The id of the instance to use as the basis of the image.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Default `40 minutes`) Used when creating the image
* `update` - (Default `40 minutes`) Used when updating the image
* `delete` - (Default `90 minutes`) Used when deregistering the image

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the image.
* `id` - The ID of the created image.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `snapshot_without_reboot` - Whether the behavior of stopping the instance before snapshotting is overrided. Always empty.

This resource also exports a full set of attributes corresponding to the arguments of the
[`aws_ami`][tf-ami] resource, allowing the properties of the created image to be used elsewhere in the
configuration.

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[tf-ami]: ami.html
