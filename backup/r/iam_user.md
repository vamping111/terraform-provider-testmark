---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "aws_iam_user"
description: |-
  Manages an IAM user.
---

[iam-users-and-projects]: https://docs.cloud.croc.ru/en/services/iam/iam.html
[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Resource: aws_iam_user

Manages an IAM user. For details about IAM users, see the [user documentation][iam-users-and-projects].

## Example Usage

### Predefined Password

```terraform
resource "aws_iam_user" "example" {
  name     = "tf-user"
  password = "********"
  email    = "example@mail.com"
}
```

### Generated Password

```terraform
resource "aws_iam_user" "example" {
  name = "tf-user"
}

output "user-password" {
  value     = aws_iam_user.example.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional, Editable) The displayed name of the user.
  If no value is specified, `name` will be used as the displayed name.
* `email` - (Optional, Editable) The email of the user.
* `name` - (Required) The name of the user. The value must start with a Latin letter and
  can only contain Latin letters, numbers, underscores (_), periods (.) and hyphens (-) (`^[a-zA-Z][a-zA-Z0-9_.-]*$`).
  The value must be 1 to 40 characters long.

~> **Note** User names are not case-sensitive. For example, you cannot create user names "TESTUSER" and "testuser" at the same time.

* `otp_required` - (Optional) Indicates whether the user is required to use two-factor authentication to log in to the web interface.
  Defaults to `false`.
* `password` - (Optional, Editable) The password of the user.
  If no value is specified, the password will be generated automatically.
* `phone` - (Optional, Editable) The phone number of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the user.
* `id` - The name of the user.
* `enabled` - Indicates whether the user is **not** locked.
* `last_login_date` - The time in [RFC3339 format] when the user last logged in to the web interface.
* `login` - The login of the user.
* `secret_key` - The secret key of the user.
* `update_date` - The time in [RFC3339 format] when the user was last updated.
* `user_id` - The ID of the user.

~> **Note** `password` and `secret_key` are exported only once when the user is created and will not be updated afterwards.

->  **Unsupported attributes**
These attributes are currently unsupported:

* `force_destroy` - The user will be destroyed even if its IAM access keys, login profile or MFA devices
  are not managed by Terraform. Always `false`.
* `path` - The path to the user. Always `""`.
* `permissions_boundary` - The ARN of the policy that is used to limit permissions for the user. Always empty.

## Import

IAM user can be imported using `name`, e.g.,

```
$ terraform import aws_iam_user.example user-name
```
