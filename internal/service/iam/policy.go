package iam

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	policyDescriptionMaxLen = 1000
	policyNameMaxLen        = 128
	policyNamePrefixMaxLen  = policyNameMaxLen - resource.UniqueIDSuffixLength
)

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, policyDescriptionMaxLen),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validResourceName(policyNameMaxLen),
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validResourceName(policyNamePrefixMaxLen),
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     verify.ValidIAMPolicyJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(iam.GrantsTypeType_Values(), false),
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
	}

	request := &iam.CreatePolicyInput{
		Description: aws.String(d.Get("description").(string)),
		Document:    aws.String(policy),
		PolicyName:  aws.String(name),
		Type:        aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating IAM policy: %s", request)
	response, err := conn.CreatePolicy(request)

	if err != nil {
		return fmt.Errorf("error creating IAM policy %s: %w", name, err)
	}

	d.SetId(aws.StringValue(response.Policy.PolicyArn))

	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	arn := d.Id()

	policy, err := FindPolicyByArn(conn, arn)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM policy (%s) not found, removing from state", arn)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM policy %s: %w", d.Id(), err)
	}

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

	policyDocument := aws.StringValue(policy.Document)
	policyToSet, err := verify.SecondJSONUnlessEquivalent(d.Get("policy").(string), policyDocument)

	if err != nil {
		return fmt.Errorf("while setting policy (%s), encountered: %w", policyToSet, err)
	}

	policyToSet, err = structure.NormalizeJsonString(policyToSet)

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", policyToSet, err)
	}

	d.Set("policy", policyToSet)
	d.Set("policy_id", policy.PolicyId)
	d.Set("type", policy.Type)

	if policy.UpdateDate != nil {
		d.Set("update_date", aws.TimeValue(policy.UpdateDate).Format(time.RFC3339))
	} else {
		d.Set("update_date", nil)
	}

	// FIXME: Test after policy versions are supported in C2 IAM API.
	// // Retrieve policy
	//
	// getPolicyVersionRequest := &iam.GetPolicyVersionInput{
	// 	PolicyArn: aws.String(d.Id()),
	// 	// VersionId: policy.DefaultVersionId,
	// }
	//
	// // Handle IAM eventual consistency
	// var getPolicyVersionResponse *iam.GetPolicyVersionOutput
	// err = resource.Retry(propagationTimeout, func() *resource.RetryError {
	// 	var err error
	// 	getPolicyVersionResponse, err = conn.GetPolicyVersion(getPolicyVersionRequest)
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
	// 	getPolicyVersionResponse, err = conn.GetPolicyVersion(getPolicyVersionRequest)
	// }
	//
	// if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
	// 	log.Printf("[WARN] IAM Policy (%s) not found, removing from state", d.Id())
	// 	d.SetId("")
	// 	return nil
	// }
	//
	// if err != nil {
	// 	return fmt.Errorf("error reading IAM policy version %s: %w", d.Id(), err)
	// }
	//
	// var policyDocument string
	// if getPolicyVersionResponse != nil && getPolicyVersionResponse.PolicyVersion != nil {
	// 	var err error
	// 	policyDocument, err = url.QueryUnescape(aws.StringValue(getPolicyVersionResponse.PolicyVersion.Document))
	// 	if err != nil {
	// 		return fmt.Errorf("error parsing IAM policy (%s) document: %w", d.Id(), err)
	// 	}
	// }
	//
	// policyToSet, err := verify.SecondJSONUnlessEquivalent(d.Get("policy").(string), policyDocument)
	//
	// if err != nil {
	// 	return fmt.Errorf("while setting policy (%s), encountered: %w", policyToSet, err)
	// }
	//
	// policyToSet, err = structure.NormalizeJsonString(policyToSet)
	//
	// if err != nil {
	// 	return fmt.Errorf("policy (%s) is invalid JSON: %w", policyToSet, err)
	// }
	//
	// d.Set("policy", policyToSet)

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	arn := d.Id()

	if d.HasChanges("description", "policy") {

		// need policy name for update
		policy, err := FindPolicyByArn(conn, arn)

		if err != nil {
			return fmt.Errorf("error reding IAM policy (%s) for updating: %w", arn, err)
		}

		input := &iam.UpdatePolicyInput{
			PolicyName: policy.PolicyName,
		}

		if d.HasChange("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		if d.HasChange("policy") {
			policyDocument, err := structure.NormalizeJsonString(d.Get("policy").(string))

			if err != nil {
				return fmt.Errorf("policy (%s) is invalid JSON: %w", policyDocument, err)
			}

			input.Document = aws.String(policyDocument)
		}

		log.Printf("[DEBUG] Modifying IAM policy: %s", input)
		_, err = conn.UpdatePolicy(input)

		if err != nil {
			return fmt.Errorf("error modifying IAM policy (%s): %s", arn, err)
		}
	}

	// FIXME: Test after policy versions are supported in C2 IAM API.
	// if d.HasChangesExcept("tags", "tags_all") {
	//
	// 	if err := policyPruneVersions(d.Id(), conn); err != nil {
	// 		return err
	// 	}
	//
	// 	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))
	//
	// 	if err != nil {
	// 		return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
	// 	}
	//
	// 	request := &iam.CreatePolicyVersionInput{
	// 		PolicyArn:      aws.String(d.Id()),
	// 		PolicyDocument: aws.String(policy),
	// 		SetAsDefault:   aws.Bool(true),
	// 	}
	//
	// 	if _, err := conn.CreatePolicyVersion(request); err != nil {
	// 		return fmt.Errorf("error updating IAM policy %s: %w", d.Id(), err)
	// 	}
	// }

	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	arn := d.Id()

	// FIXME: Test after policy versions are supported in C2 IAM API.
	// if err := PolicyDeleteNondefaultVersions(d.Id(), conn); err != nil {
	// 	return err
	// }

	// need policy name for delete
	policy, err := FindPolicyByArn(conn, arn)

	if tfresource.NotFound(err) {
		log.Printf("[WARN] IAM policy (%s) not found, removing from state", arn)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reding IAM policy (%s) for deleting: %w", arn, err)
	}

	request := &iam.DeletePolicyInput{
		PolicyName: policy.PolicyName,
	}

	log.Printf("[DEBUG] Deleting IAM policy: %s", request)
	_, err = conn.DeletePolicy(request)

	if tfawserr.ErrCodeEquals(err, PolicyNotFoundCode) {
		log.Printf("[WARN] IAM policy (%s) not found, removing from state", arn)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM policy (%s): %w", arn, err)
	}

	return nil
}

// policyPruneVersions deletes the oldest versions.
//
// Old versions are deleted until there are 4 or less remaining, which means at
// least one more can be created before hitting the maximum of 5.
//
// The default version is never deleted.
//
//nolint:unused // Policy versions are unsupported.
func policyPruneVersions(arn string, conn *iam.IAM) error {
	versions, err := policyListVersions(arn, conn)
	if err != nil {
		return err
	}
	if len(versions) < 5 {
		return nil
	}

	var oldestVersion *iam.PolicyVersion

	for _, version := range versions {
		if *version.IsDefaultVersion {
			continue
		}
		if oldestVersion == nil ||
			version.CreateDate.Before(*oldestVersion.CreateDate) {
			oldestVersion = version
		}
	}

	err1 := policyDeleteVersion(arn, aws.StringValue(oldestVersion.VersionId), conn)
	return err1
}

func PolicyDeleteNondefaultVersions(arn string, conn *iam.IAM) error {
	versions, err := policyListVersions(arn, conn)
	if err != nil {
		return err
	}

	for _, version := range versions {
		if *version.IsDefaultVersion {
			continue
		}
		if err := policyDeleteVersion(arn, aws.StringValue(version.VersionId), conn); err != nil {
			return err
		}
	}

	return nil
}

func policyDeleteVersion(arn, versionID string, conn *iam.IAM) error {
	request := &iam.DeletePolicyVersionInput{
		PolicyArn: aws.String(arn),
		VersionId: aws.String(versionID),
	}

	_, err := conn.DeletePolicyVersion(request)
	if err != nil {
		return fmt.Errorf("Error deleting version %s from IAM policy %s: %w", versionID, arn, err)
	}
	return nil
}

func policyListVersions(arn string, conn *iam.IAM) ([]*iam.PolicyVersion, error) {
	request := &iam.ListPolicyVersionsInput{
		PolicyArn: aws.String(arn),
	}

	response, err := conn.ListPolicyVersions(request)
	if err != nil {
		return nil, fmt.Errorf("Error listing versions for IAM policy %s: %w", arn, err)
	}
	return response.Versions, nil
}
