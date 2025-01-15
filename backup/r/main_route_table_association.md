---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_main_route_table_association"
description: |-
  Provides a resource for managing the main routing table of a VPC.
---

# Resource: aws_main_route_table_association

Provides a resource for managing the main routing table of a VPC.

~> **Note** **Do not** use both `aws_default_route_table` to manage a default route table **and** `aws_main_route_table_association` with the same VPC due to possible route conflicts. See [aws_default_route_table][tf-default-route-table] documentation for more details.
For more information, see the documentation on [Route Tables][route-tables]. For information about managing normal route tables in Terraform, see [`aws_route_table`][tf-route-table].

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_main_route_table_association" "example" {
  vpc_id         = aws_vpc.example.id
  route_table_id = aws_route_table.example.id
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) ID of the VPC whose main route table should be set
* `route_table_id` - (Required) ID of the Route Table to set as the new
  main route table for the target VPC

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Route Table Association
* `original_route_table_id` - Used internally, see __Notes__ below

## Notes

On VPC creation, the cloud always creates an initial Main Route Table. This
resource records the ID of that Route Table under `original_route_table_id`.
The "Delete" action for a `main_route_table_association` consists of resetting
this original table as the Main Route Table for the VPC. You'll see this
additional Route Table in the cloud console; it must remain intact in order for
the `main_route_table_association` delete to work properly.

[route-tables]: https://docs.cloud.croc.ru/en/services/networks/routetables.html
[tf-route-table]: route_table.html
[tf-default-route-table]: default_route_table.html
