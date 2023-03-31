---
subcategory: "EBS (EC2)"
layout: "aws"
page_title: "AWS: aws_ebs_snapshot_ids"
description: |-
  Provides a list of EBS snapshot IDs.
---

# Data Source: aws_ebs_snapshot_ids

Use this data source to get a list of EBS snapshot IDs matching the specified criteria.

## Example Usage

```terraform
data "aws_ebs_snapshot_ids" "ebs_snapshot_ids" {
  owners = ["self"]

  filter {
    name   = "volume-size"
    values = ["40"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `owners` - (Optional) Returns the snapshots owned by the specified owner ID. Multiple owners can be specified.
* `restorable_by_user_ids` - (Optional) One or more CROC Cloud project IDs that can create volumes from the snapshot.
* `filter` - (Optional) One or more name/value pairs to filter.

For more information about filtering, see the [EC2 API documentation][describe-snapshots].

## Attributes Reference

* `id` - Region (for example, `croc`).
* `ids` - Set of EBS snapshot IDs, sorted by creation time in descending order.

[describe-snapshots]: https://docs.cloud.croc.ru/en/api/ec2/snapshots/DescribeSnapshots.html
[tf-ebs-snapshot]: ebs_snapshot.html
