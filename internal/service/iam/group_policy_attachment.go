package iam

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceGroupPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupPolicyAttachmentCreate,
		Read:   resourceGroupPolicyAttachmentRead,
		Delete: resourceGroupPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupPolicyAttachmentImport,
		},

		Schema: map[string]*schema.Schema{
			"group_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"policy_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	groupArn := d.Get("group_arn").(string)
	policyArn := d.Get("policy_arn").(string)

	// need group name for attaching policy
	group, _, err := FindGroupByArn(conn, groupArn)

	if err != nil {
		return fmt.Errorf("error reding IAM group (%s) for attaching policy: %w", groupArn, err)
	}

	input := &iam.AttachGroupPolicyInput{
		GroupName: group.GroupName,
		PolicyArn: aws.String(policyArn),
	}

	log.Printf("[DEBUG] Attaching IAM policy to IAM group: %s", input)
	if _, err := conn.AttachGroupPolicy(input); err != nil {
		return fmt.Errorf("error attaching IAM policy (%s) to IAM group (%s): %w", policyArn, groupArn, err)
	}

	d.SetId(composeUserGroupAttachmentId(groupArn, policyArn))

	return resourceGroupPolicyAttachmentRead(d, meta)
}

func resourceGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	id := d.Id()
	groupArn := d.Get("group_arn").(string)
	policyArn := d.Get("policy_arn").(string)

	_, err := FindGroupAttachedPolicy(conn, groupArn, policyArn)

	if !d.IsNewResource() && (tfresource.NotFound(err) || tfawserr.ErrCodeEquals(err, GroupNotFoundCode)) {
		log.Printf("[WARN] IAM group managed policy attachment (%s) not found, removing from state", id)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM group managed policy attachment (%s): %w", id, err)
	}

	return nil
}

func resourceGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	id := d.Id()
	groupArn := d.Get("group_arn").(string)
	policyArn := d.Get("policy_arn").(string)

	// need group name for detaching policy
	group, _, err := FindGroupByArn(conn, groupArn)

	if tfresource.NotFound(err) {
		log.Printf("[WARN] IAM group managed policy attachment (%s) not found, removing from state", id)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reding IAM group (%s) for detaching policy: %w", groupArn, err)
	}

	input := &iam.DetachGroupPolicyInput{
		GroupName: group.GroupName,
		PolicyArn: aws.String(policyArn),
	}

	log.Printf("[DEBUG] Detaching IAM policy from IAM group: %s", input)
	if _, err := conn.DetachGroupPolicy(input); err != nil {
		return fmt.Errorf("error detaching IAM policy (%s) from IAM group (%s): %w", policyArn, groupArn, err)
	}

	return nil
}

func resourceGroupPolicyAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), "#", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected <group-arn>#<policy_arn>", d.Id())
	}

	groupArn := idParts[0]
	policyArn := idParts[1]

	d.Set("group_arn", groupArn)
	d.Set("policy_arn", policyArn)

	d.SetId(composeUserGroupAttachmentId(groupArn, policyArn))

	return []*schema.ResourceData{d}, nil
}

// composeUserPolicyAttachmentId constructs an id for a group policy attachment.
//
// Format: groupArn#policyArn
func composeUserGroupAttachmentId(groupArn, policyArn string) string {
	return fmt.Sprintf("%s#%s", groupArn, policyArn)
}
