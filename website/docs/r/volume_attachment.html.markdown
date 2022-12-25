---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "AWS: aws_volume_attachment"
description: |-
  Provides an EBS Volume Attachment
---

# Resource: aws_volume_attachment

Provides an EBS volume attachment as a top level resource, to attach and detach volumes from instances.

~> **NOTE on EBS block devices:** If you use `ebs_block_device` on an `aws_instance`, Terraform will assume management over the full set of non-root EBS block devices for the instance, and treats additional block devices as drift. For this reason, `ebs_block_device` cannot be mixed with external `aws_ebs_volume` + `aws_volume_attachment` resources for a given instance.

## Example Usage

```terraform
variable instance_id {}

resource "aws_ebs_volume" "example" {
  availability_zone = "ru-msk-vol52"
  size              = 1
}

resource "aws_volume_attachment" "ebs_att" {
  device_name = "disk2"
  volume_id   = aws_ebs_volume.example.id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `device_name` - (Required) The device name to expose to the instance.
* `instance_id` - (Required) ID of the instance to attach to.
* `volume_id` - (Required) ID of the volume to be attached.
* `stop_instance_before_detaching` - (Optional) Set this to true to ensure
that the target instance is stopped before trying to detach the volume.
Stops the instance, if it is not already stopped.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `device_name` - The device name exposed to the instance.
* `instance_id` - ID of the instance.
* `volume_id` - ID of the volume.

->  **Unsupported attributes**
These exported attributes are currently unsupported by CROC Cloud:

* `force_detach` - Whether force volume detaching is enabled. Always empty.
* `skip_destroy` - Whether volume attachment will be removed from Terraform state instead of detach if the instance is destroyed. Always empty.

## Import

EBS volume attachments can be imported using `DEVICE_NAME:VOLUME_ID:INSTANCE_ID`, e.g.,

```
$ terraform import aws_volume_attachment.example disk1:vol-12345678:i-12345678
```
