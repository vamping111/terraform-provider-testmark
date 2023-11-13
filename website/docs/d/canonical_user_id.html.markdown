---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "CROC Cloud: aws_canonical_user_id"
description: |-
  Provides the canonical user ID (CROC Cloud S3 User ID) associated with the provider
  connection to CROC Cloud.
---

# Data Source: aws_canonical_user_id

The Canonical User ID data source allows access to the CROC Cloud S3 User ID
for the effective account in which Terraform is working.  

~> **NOTE:** To use this data source, you must have the `s3:ListAllMyBuckets` permission.

## Example Usage

```terraform
data "aws_canonical_user_id" "current" {}

output "canonical_user_id" {
  value = data.aws_canonical_user_id.current.id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The canonical user ID (CROC Cloud S3 User ID) associated with the provider connection to CROC Cloud.
* `display_name` - The human-friendly name linked to the canonical user ID. The bucket owner's display name.
