---
layout: "aws"
page_title: "Provider: AWS"
description: |-
  Use the Amazon Web Services (AWS) provider to interact with the many resources supported by AWS. You must configure the provider with the proper credentials before you can use it.
---

# CROC Cloud Provider

Use the CROC Cloud Provider based on the Terraform AWS Provider to interact with cloud services.
The provider needs to be configured with the proper credentials before you can use it.

Use the navigation to the left to read about the available resources.

To learn the basics of Terraform using this provider, follow the
hands-on [get started tutorials][hashicorp-tutorials] on HashiCorp's Learn platform.

Examples of using AWS services with Terraform: [AWS services tutorials][aws-tutorials].

[hashicorp-tutorials]: https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code?in=terraform/aws-get-started&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS
[aws-tutorials]: https://learn.hashicorp.com/collections/terraform/aws?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS

## Example Usage

Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    aws = {
      source  = "c2devel/croccloud"
      version = "4.14.0-CROC1"
    }
  }
}

# Configure the croccloud provider.
# Section name is `aws` due to backward compatibility.
provider "aws" {
  region = "croc"
}

# Create a VPC
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}
```

## Authentication and Configuration

Configuration for the CROC Cloud Provider can be derived from several sources,
which are applied in the following order:

1. Parameters in the provider configuration
2. Environment variables
3. Shared configuration and credentials files

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform
configuration and risks secret leakage should this file ever be committed to a
public version control system.

Credentials can be provided by adding an `access_key`, `secret_key` to the `aws` provider block.

Usage:

```terraform
provider "aws" {
  region     = "croc"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment Variables

Credentials can be provided by using the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` environment variables.
The region can be set using the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variables.

For example:

```terraform
provider "aws" {}
```

```sh
$ export AWS_ACCESS_KEY_ID="my-access-key"
$ export AWS_SECRET_ACCESS_KEY="my-secret-key"
$ export AWS_REGION="croc"
$ terraform plan
```

### Shared Configuration and Credentials Files

CROC Cloud Provider can use [AWS shared configuration and credentials files][aws-configure-files].
By default, these files are located at `$HOME/.aws/config` and `$HOME/.aws/credentials` on Linux and macOS,
and `"%USERPROFILE%\.aws\config"` and `"%USERPROFILE%\.aws\credentials"` on Windows.

[aws-configure-files]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html

If no named profile is specified, the `default` profile is used.
Use the `profile` parameter or `AWS_PROFILE`, `AWS_DEFAULT_PROFILE` environment variables to specify a named profile.

The locations of the shared configuration and credentials files can be configured using either
the parameters `shared_config_files` and `shared_credentials_files`
or the environment variables `AWS_CONFIG_FILE` and `AWS_SHARED_CREDENTIALS_FILE`.

For example:

```terraform
provider "aws" {
  shared_config_files      = ["/Users/tf_user/.aws/conf"]
  shared_credentials_files = ["/Users/tf_user/.aws/creds"]
  profile                  = "customprofile"
}
```

## CROC Cloud Provider Full Configuration

For more information about CROC Cloud Provider configuration, see documentation on [using Terraform in CROC Cloud][terraform].

[terraform]: https://docs.cloud.croc.ru/en/api/tools/terraform.html
