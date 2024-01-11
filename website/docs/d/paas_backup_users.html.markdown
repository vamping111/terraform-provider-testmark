---
subcategory: "PaaS"
layout: "aws"
page_title: "CROC Cloud: aws_paas_backup_users"
description: |-
  Provides information about users with the PaaS Backup User project privileges.
---

[RFC3339 format]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.8

# Data Source: aws_paas_backup_users

Provides information about users with the PaaS Backup User project privileges.

## Example Usage

```terraform
data "aws_paas_backup_users" "selected" {
  active_only = true
}

output "backup-user-logins" {
  value = data.aws_paas_backup_users.selected.users[*].login
}
```

## Argument Reference

The following arguments are supported:

* `active_only` - (Optional) Indicates whether to filter only active users. Defaults to `false`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The region (e.g., `croc`).
* `users` - List of users. Each user has the following structure:
    * `email` - The user email.
    * `enabled` - Indicates whether the user is active.
    * `id` - The ID of the user.
    * `last_login_time` - The time when the user was last active in [RFC3339 format].
    * `login` - The user login.
    * `modify_time` - The time when user info was last modified in [RFC3339 format].
    * `name` - The user name.
    * `grants` - List of user's project privileges.
