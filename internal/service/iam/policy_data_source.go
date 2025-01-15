package iam

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  verify.ValidARN,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"arn"},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	arn := d.Get("arn").(string)
	name := d.Get("name").(string)
	pathPrefix := d.Get("path_prefix").(string)

	var policy *iam.Policy

	if arn == "" {
		policies, err := FindPolicies(conn, arn, name, pathPrefix)

		if err != nil {
			return fmt.Errorf("error reading IAM policy (%s): %w", PolicySearchDetails(arn, name, pathPrefix), err)
		}

		if len(policies) == 0 {
			return fmt.Errorf("no IAM policy found matching criteria (%s); try different search", PolicySearchDetails(arn, name, pathPrefix))
		}

		if len(policies) > 1 {
			return fmt.Errorf("multiple IAM policies found matching criteria (%s); try different search", PolicySearchDetails(arn, name, pathPrefix))
		}

		arn = aws.StringValue(policies[0].PolicyArn)
	}

	policy, err := FindPolicyByArn(conn, arn)

	if err != nil {
		return fmt.Errorf("error reading IAM policy (%s): %w", arn, err)
	}

	d.SetId(aws.StringValue(policy.PolicyArn))

	d.Set("arn", policy.PolicyArn)

	if policy.CreateDate != nil {
		d.Set("create_date", aws.TimeValue(policy.CreateDate).Format(time.RFC3339))
	} else {
		d.Set("create_date", nil)
	}

	d.Set("description", policy.Description)
	d.Set("name", policy.PolicyName)
	d.Set("owner", policy.Owner)
	d.Set("path", policy.Path)
	d.Set("policy", policy.Document)
	d.Set("policy_id", policy.PolicyId)

	if err := d.Set("tags", KeyValueTags(policy.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	d.Set("type", policy.Type)

	if policy.UpdateDate != nil {
		d.Set("update_date", aws.TimeValue(policy.UpdateDate).Format(time.RFC3339))
	} else {
		d.Set("update_date", nil)
	}

	// FIXME: Test after policy versions are supported in C2 IAM API.
	// // Retrieve policy
	// policyVersionInput := &iam.GetPolicyVersionInput{
	// 	PolicyArn: policy.Arn,
	// 	VersionId: policy.DefaultVersionId,
	// }
	//
	// // Handle IAM eventual consistency
	// var policyVersionOutput *iam.GetPolicyVersionOutput
	// err = resource.Retry(propagationTimeout, func() *resource.RetryError {
	// 	var err error
	// 	policyVersionOutput, err = conn.GetPolicyVersion(policyVersionInput)
	//
	// 	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
	// 		return resource.RetryableError(err)
	// 	}
	//
	// 	if err != nil {
	// 		return resource.NonRetryableError(err)
	// 	}
	//
	// 	return nil
	// })
	//
	// if tfresource.TimedOut(err) {
	// 	policyVersionOutput, err = conn.GetPolicyVersion(policyVersionInput)
	// }
	//
	// if err != nil {
	// 	return fmt.Errorf("error reading IAM Policy (%s) version: %w", policyArn, err)
	// }
	//
	// if policyVersionOutput == nil || policyVersionOutput.PolicyVersion == nil {
	// 	return fmt.Errorf("error reading IAM Policy (%s) version: empty output", policyArn)
	// }
	//
	// policyVersion := policyVersionOutput.PolicyVersion
	//
	// var policyDocument string
	// if policyVersion != nil {
	// 	policyDocument, err = url.QueryUnescape(aws.StringValue(policyVersion.Document))
	// 	if err != nil {
	// 		return fmt.Errorf("error parsing IAM Policy (%s) document: %w", policyArn, err)
	// 	}
	// }
	//
	// d.Set("policy", policyDocument)

	return nil
}

// PolicySearchDetails returns the configured search criteria as a printable string
func PolicySearchDetails(arn, name, pathPrefix string) string {
	var policyDetails []string
	if arn != "" {
		policyDetails = append(policyDetails, fmt.Sprintf("ARN: %s", arn))
	}
	if name != "" {
		policyDetails = append(policyDetails, fmt.Sprintf("Name: %s", name))
	}
	if pathPrefix != "" {
		policyDetails = append(policyDetails, fmt.Sprintf("PathPrefix: %s", pathPrefix))
	}

	return strings.Join(policyDetails, ", ")
}
