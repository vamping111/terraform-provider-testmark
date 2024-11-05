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

func ResourceUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserPolicyAttachmentCreate,
		Read:   resourceUserPolicyAttachmentRead,
		Delete: resourceUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserPolicyAttachmentImport,
		},

		Schema: map[string]*schema.Schema{
			"policy_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	userName := d.Get("user").(string)
	policyArn := d.Get("policy_arn").(string)

	input := &iam.AttachUserPolicyInput{
		UserName:  aws.String(userName),
		PolicyArn: aws.String(policyArn),
	}

	var projectName string
	if v, ok := d.GetOk("project"); ok {
		projectName = v.(string)
		input.ProjectName = aws.String(projectName)
	}

	log.Printf("[DEBUG] Attaching IAM policy to IAM user: %s", input)
	_, err := conn.AttachUserPolicy(input)
	if err != nil {
		if projectName == "" {
			return fmt.Errorf("error attaching IAM policy (%s) to IAM user (%s): %w", policyArn, userName, err)
		}

		return fmt.Errorf(
			"error attaching IAM policy (%s) to IAM user (%s) for IAM project (%s): %w",
			policyArn, userName, projectName, err,
		)
	}

	d.SetId(composeUserPolicyAttachmentId(userName, projectName, policyArn))

	return resourceUserPolicyAttachmentRead(d, meta)
}

func resourceUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	id := d.Id()
	userName := d.Get("user").(string)
	projectName := d.Get("project").(string)
	policyArn := d.Get("policy_arn").(string)

	var err error
	if projectName == "" {
		_, err = FindUserAttachedGlobalPolicy(conn, userName, policyArn)
	} else {
		_, err = FindUserAttachedProjectPolicy(conn, userName, policyArn, projectName)
	}

	if !d.IsNewResource() && (tfresource.NotFound(err) || tfawserr.ErrCodeEquals(err, UserNotFoundCode)) {
		log.Printf("[WARN] IAM user managed policy attachment (%s) not found, removing from state", id)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM user managed policy attachment (%s): %w", id, err)
	}

	return nil
}

func resourceUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	userName := d.Get("user").(string)
	policyArn := d.Get("policy_arn").(string)

	input := &iam.DetachUserPolicyInput{
		UserName:  aws.String(userName),
		PolicyArn: aws.String(policyArn),
	}

	var projectName string
	if v, ok := d.GetOk("project"); ok {
		projectName = v.(string)
		input.ProjectName = aws.String(projectName)
	}

	log.Printf("[DEBUG] Detaching IAM policy from IAM user: %s", input)
	_, err := conn.DetachUserPolicy(input)

	if tfawserr.ErrCodeEquals(err, UserNotFoundCode, ProjectNotFoundCode) {
		log.Printf("[WARN] IAM user managed policy attachment (%s) not found, removing from state", d.Id())
		return nil
	}

	if err != nil {
		if projectName == "" {
			return fmt.Errorf("error detaching IAM policy (%s) from IAM user (%s): %w", policyArn, userName, err)
		}

		return fmt.Errorf(
			"error detaching IAM policy (%s) from IAM user (%s) for IAM project (%s): %w",
			policyArn, userName, projectName, err,
		)
	}

	return nil
}

func resourceUserPolicyAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), "#", 3)
	if len(idParts) < 2 || len(idParts) > 3 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf(
			"unexpected format of ID (%q), expected <user-name>[#<project-name>]#<policy-arn>", d.Id(),
		)
	}

	userName := idParts[0]

	var projectName, policyArn string
	if len(idParts) == 3 {
		projectName = idParts[1]
		policyArn = idParts[2]
	} else {
		policyArn = idParts[1]
	}

	d.Set("user", userName)
	d.Set("policy_arn", policyArn)
	d.Set("project", projectName)

	d.SetId(composeUserPolicyAttachmentId(userName, projectName, policyArn))

	return []*schema.ResourceData{d}, nil
}

func DetachPolicyFromUser(conn *iam.IAM, user string, arn string) error {
	_, err := conn.DetachUserPolicy(&iam.DetachUserPolicyInput{
		UserName:  aws.String(user),
		PolicyArn: aws.String(arn),
	})
	return err
}

// composeUserPolicyAttachmentId constructs an id for a user policy attachment.
//
// Format: userName[#projectName]#policyArn
func composeUserPolicyAttachmentId(userName, projectName, policyArn string) string {
	if projectName == "" {
		return fmt.Sprintf("%s#%s", userName, policyArn)
	}

	return fmt.Sprintf("%s#%s#%s", userName, projectName, policyArn)
}
