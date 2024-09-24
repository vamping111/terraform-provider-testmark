package iam

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	groupNameMaxLen = 128
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		// FIXME: Test after UpdateGroup is supported in C2 IAM API.
		// Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
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
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validResourceName(groupNameMaxLen),
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
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(iam.GrantsTypeType_Values(), false),
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Get("name").(string)

	request := &iam.CreateGroupInput{
		GroupName: aws.String(name),
		Type:      aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating IAM group: %s", request)
	createResp, err := conn.CreateGroup(request)

	if err != nil {
		return fmt.Errorf("error creating IAM group %s: %w", name, err)
	}

	d.SetId(aws.StringValue(createResp.Group.GroupArn))

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	arn := d.Id()

	group, _, err := FindGroupByArn(conn, arn)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM group (%s) not found, removing from state", arn)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM group (%s): %w", arn, err)
	}

	d.Set("arn", group.GroupArn)

	if group.CreateDate != nil {
		d.Set("create_date", aws.TimeValue(group.CreateDate).Format(time.RFC3339))
	} else {
		d.Set("create_date", nil)
	}

	d.Set("group_id", group.GroupId)
	d.Set("name", group.GroupName)
	d.Set("owner", group.Owner)
	d.Set("path", group.Path)
	d.Set("type", group.Type)

	return nil
}

//nolint:unused // UpdateGroup is unsupported.
func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChanges("name", "path") {
		conn := meta.(*conns.AWSClient).IAMConn
		on, nn := d.GetChange("name")
		_, np := d.GetChange("path")

		request := &iam.UpdateGroupInput{
			GroupName:    aws.String(on.(string)),
			NewGroupName: aws.String(nn.(string)),
			NewPath:      aws.String(np.(string)),
		}
		_, err := conn.UpdateGroup(request)
		if err != nil {
			return fmt.Errorf("Error updating IAM Group %s: %s", d.Id(), err)
		}
		d.SetId(nn.(string))
		return resourceGroupRead(d, meta)
	}
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	arn := d.Id()

	// need group name for delete
	group, _, err := FindGroupByArn(conn, arn)

	if tfresource.NotFound(err) {
		log.Printf("[WARN] IAM group (%s) not found, removing from state", arn)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reding IAM group (%s) for deleting: %w", arn, err)
	}

	input := &iam.DeleteGroupInput{
		GroupName: group.GroupName,
	}

	log.Printf("[DEBUG] Deleting IAM group: %s", input)
	_, err = conn.DeleteGroup(input)

	if tfawserr.ErrCodeEquals(err, GroupNotFoundCode) {
		log.Printf("[WARN] IAM group (%s) not found, removing from state", arn)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM group (%s): %w", arn, err)
	}

	return nil
}

func DeleteGroupPolicyAttachments(conn *iam.IAM, groupName string) error {
	var attachedPolicies []*iam.AttachedPolicy
	input := &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupName),
	}

	err := conn.ListAttachedGroupPoliciesPages(input, func(page *iam.ListAttachedGroupPoliciesOutput, lastPage bool) bool {
		attachedPolicies = append(attachedPolicies, page.AttachedPolicies...)

		return !lastPage
	})

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing IAM Group (%s) policy attachments for deletion: %w", groupName, err)
	}

	for _, attachedPolicy := range attachedPolicies {
		input := &iam.DetachGroupPolicyInput{
			GroupName: aws.String(groupName),
			PolicyArn: attachedPolicy.PolicyArn,
		}

		_, err := conn.DetachGroupPolicy(input)

		if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error detaching IAM Group (%s) policy (%s): %w", groupName, aws.StringValue(attachedPolicy.PolicyArn), err)
		}
	}

	return nil
}

func DeleteGroupPolicies(conn *iam.IAM, groupName string) error {
	var inlinePolicies []*string
	input := &iam.ListGroupPoliciesInput{
		GroupName: aws.String(groupName),
	}

	err := conn.ListGroupPoliciesPages(input, func(page *iam.ListGroupPoliciesOutput, lastPage bool) bool {
		inlinePolicies = append(inlinePolicies, page.PolicyNames...)
		return !lastPage
	})

	if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing IAM Group (%s) inline policies for deletion: %w", groupName, err)
	}

	for _, policyName := range inlinePolicies {
		input := &iam.DeleteGroupPolicyInput{
			GroupName:  aws.String(groupName),
			PolicyName: policyName,
		}

		_, err := conn.DeleteGroupPolicy(input)

		if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error deleting IAM Group (%s) inline policy (%s): %w", groupName, aws.StringValue(policyName), err)
		}
	}

	return nil
}
