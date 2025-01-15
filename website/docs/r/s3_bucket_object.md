---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket_object"
description: |-
  Provides an S3 object resource.
---

[canned-acl]: https://docs.cloud.croc.ru/en/api/s3/acl.html#cannedacl
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[w3c cache_control]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9
[w3c content_disposition]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1
[w3c content_encoding]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11

# Resource: aws_s3_bucket_object

~> **Note** The `aws_s3_bucket_object` resource is DEPRECATED and will be removed in a future version! Use [`aws_s3_object`](s3_object.md) instead, where new features and fixes will be added. When replacing `aws_s3_bucket_object` with `aws_s3_object` in your configuration, on the next apply, Terraform will recreate the object. If you prefer to not have Terraform recreate the object, import the object using `aws_s3_object`.

Provides an S3 object resource.

## Example Usage

### Uploading a file to a bucket

```terraform
resource "aws_s3_object" "example" {
  bucket = "tf-example"
  key    = "new_object_key"
  source = "path/to/file"

  # The filemd5() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
  # etag = "${md5(file("path/to/file"))}"
  etag = filemd5("path/to/file")
}
```

## Argument Reference

-> **Note** If you specify `content_encoding` you are responsible for encoding the body appropriately. `source`, `content`, and `content_base64` all expect already encoded/compressed bytes.

The following arguments are required:

* `bucket` - (Required) Name of the bucket to put the file in.
* `key` - (Required) Name of the object once it is in the bucket.

The following arguments are optional:

* `acl` - (Optional) [Canned ACL][canned-acl] to apply. Valid values are `private`, `public-read`, `public-read-write`, `authenticated-read`. Defaults to `private`.
* `cache_control` - (Optional) Caching behavior along the request/reply chain Read [w3c cache_control] for further details.
* `content_base64` - (Optional, Conflicts with `source` and `content`) Base64-encoded data that will be decoded and uploaded as raw bytes for the object content. This allows safely uploading non-UTF8 binary data, but is recommended only for small content such as the result of the `gzipbase64` function with small text strings. For larger objects, use `source` to stream the content from a disk file.
* `content_disposition` - (Optional) Presentational information for the object. Read [w3c content_disposition] for further information.
* `content_encoding` - (Optional) Content encodings that have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [w3c content_encoding] for further information.
* `content_language` - (Optional) Language the content is in e.g., en-US or en-GB.
* `content_type` - (Optional) Standard MIME type describing the format of the object data, e.g., application/octet-stream. All Valid MIME Types are valid for this input.
* `content` - (Optional, Conflicts with `source` and `content_base64`) Literal string value to use as the object content, which will be uploaded as UTF-8-encoded text.
* `etag` - (Optional) Triggers updates when the value changes. The only meaningful value is `filemd5("path/to/file")` (Terraform 0.11.12 or later) or `${md5(file("path/to/file"))}` (Terraform 0.11.11 or earlier). If an object is larger than 5 MB, it will be uploaded as a Multipart Upload, and therefore the ETag will not be an MD5 digest (see `source_hash` instead).
* `metadata` - (Optional) Map of keys/values to provision metadata (will be automatically prefixed by `x-amz-meta-`, note that only lowercase label are currently supported by the AWS Go API).
* `source_hash` - (Optional) Triggers updates like `etag` but useful to address `etag` encryption limitations. Set using `filemd5("path/to/source")` (Terraform 0.11.12 or later). (The value is only stored in state and not saved by AWS.)
* `source` - (Optional, Conflicts with `content` and `content_base64`) Path to a file that will be read and uploaded as raw bytes for the object content.
* `tags` - (Optional) Map of tags to assign to the object. If configured with a provider [`default_tags` configuration block][default-tags] present, tags with matching keys will overwrite those defined at the provider-level.
* `website_redirect` - (Optional) Target URL for [website redirect](http://docs.aws.amazon.com/AmazonS3/latest/dev/how-to-page-redirect.html).

If no content is provided through `source`, `content` or `content_base64`, then the object will be empty.

-> **Note** Terraform ignores all leading `/`s in the object's `key` and treats multiple `/`s in the rest of the object's `key` as a single `/`, so values of `/index.html` and `index.html` correspond to the same S3 object as do `first//second///third//` and `first/second/third/`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `etag` - ETag generated for the object (an MD5 sum of the object content). For plaintext objects the hash is an MD5 digest of the object data. For objects created by either the Multipart Upload or Part Copy operation, the hash is not an MD5 digest, regardless of the method of encryption.
* `id` - `key` of the resource supplied above.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block][default-tags].
* `version_id` - Unique version ID value for the object, if bucket versioning is enabled.

->  **Unsupported attributes**
These exported attributes are currently unsupported:

* `bucket_key_enabled` - Whether to use [Amazon S3 Bucket Keys](https://docs.aws.amazon.com/AmazonS3/latest/dev/bucket-key.html) for SSE-KMS. Always `false`.
* `force_destroy` - Whether to allow the object to be deleted by removing any legal hold on any object version. This value should be set to `true` only if the bucket has S3 object lock enabled. Always `false`.
* `kms_key_id` - The ARN for the KMS encryption key. Always empty.
* `object_lock_legal_hold_status` - [Legal hold](https://docs.aws.amazon.com/AmazonS3/latest/dev/object-lock-overview.html#object-lock-legal-holds) status that you want to apply to the specified object. Always `""`.
* `object_lock_mode` - Object lock [retention mode](https://docs.aws.amazon.com/AmazonS3/latest/dev/object-lock-overview.html#object-lock-retention-modes) that you want to apply to this object. Always `""`.
* `object_lock_retain_until_date` - Date and time, in [RFC3339 format](https://tools.ietf.org/html/rfc3339#section-5.8), when this object's object lock will [expire](https://docs.aws.amazon.com/AmazonS3/latest/dev/object-lock-overview.html#object-lock-retention-periods). Always `""`.
* `server_side_encryption` - Server-side encryption of the object in S3. Always `""`.
* `storage_class` - [Storage Class](https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html#AmazonS3-PutObject-request-header-StorageClass) for the object. Always `STANDARD`.

## Import

Objects can be imported using the `id`. The `id` is the bucket name and the key together e.g.,

```
$ terraform import aws_s3_bucket_object.example some-bucket-name/some/key.txt
```

Additionally, s3 url syntax can be used, e.g.,

```
$ terraform import aws_s3_bucket_object.example s3://some-bucket-name/some/key.txt
```
