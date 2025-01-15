---
subcategory: "EKS (Elastic Kubernetes)"
layout: "aws"
page_title: "aws_eks_cluster_auth"
description: |-
  Get an authentication token to communicate with an EKS Cluster
---

# Data Source: aws_eks_cluster_auth

Get an authentication token to communicate with an EKS cluster.

## Example Usage

```terraform
data "aws_eks_cluster_auth" "example" {
  name = "example"
}
```

## Argument Reference

* `name` - (Required) The name of the cluster

## Attributes Reference

* `id` - Name of the cluster.
* `token` - The token to use to authenticate with the cluster.
