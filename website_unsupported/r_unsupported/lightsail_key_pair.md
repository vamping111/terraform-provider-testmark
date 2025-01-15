---
subcategory: "Lightsail"
layout: "aws"
page_title: "AWS: aws_lightsail_key_pair"
description: |-
  Provides an Lightsail Key Pair
---

# Resource: aws_lightsail_key_pair

Provides a Lightsail Key Pair, for use with Lightsail Instances. These key pairs
are separate from EC2 Key Pairs, and must be created or imported for use with
Lightsail.

~> **Note:** Lightsail is currently only supported in a limited number of AWS Regions, please see ["Regions and Availability Zones in Amazon Lightsail"](https://lightsail.aws.amazon.com/ls/docs/overview/article/understanding-regions-and-availability-zones-in-amazon-lightsail) for more details

## Example Usage

### Create New Key Pair

```terraform
# Create a new Lightsail Key Pair
resource "aws_lightsail_key_pair" "lg_key_pair" {
  name = "lg_key_pair"
}
```

### Create New Key Pair with PGP Encrypted Private Key

```terraform
resource "aws_lightsail_key_pair" "lg_key_pair" {
  name    = "lg_key_pair"
  pgp_key = "keybase:keybaseusername"
}
```

### Existing Public Key Import

```terraform
resource "aws_lightsail_key_pair" "lg_key_pair" {
  name       = "importing"
  public_key = file("~/.ssh/id_rsa.pub")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the Lightsail Key Pair. If omitted, a unique
name will be generated by Terraform
* `pgp_key` – (Optional) An optional PGP key to encrypt the resulting private
key material. Only used when creating a new key pair
* `public_key` - (Required) The public key material. This public key will be
imported into Lightsail

~> **NOTE:** a PGP key is not required, however it is strongly encouraged.
Without a PGP key, the private key material will be stored in state unencrypted.
`pgp_key` is ignored if `public_key` is supplied.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name used for this key pair
* `arn` - The ARN of the Lightsail key pair
* `fingerprint` - The MD5 public key fingerprint as specified in section 4 of RFC 4716.
* `public_key` - the public key, base64 encoded
* `private_key` - the private key, base64 encoded. This is only populated
when creating a new key, and when no `pgp_key` is provided
* `encrypted_private_key` – the private key material, base 64 encoded and
encrypted with the given `pgp_key`. This is only populated when creating a new
key and `pgp_key` is supplied
* `encrypted_fingerprint` - The MD5 public key fingerprint for the encrypted
private key

## Import

Lightsail Key Pairs cannot be imported, because the private and public key are
only available on initial creation.