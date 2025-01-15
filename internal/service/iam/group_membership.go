package iam

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Update: resourceGroupMembershipUpdate,
		Delete: resourceGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"group_arn": {
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	groupArn := d.Get("group_arn").(string)
	projectName := d.Get("project").(string)
	userNames := flex.ExpandStringSet(d.Get("users").(*schema.Set))

	if err := validateIfUsersExist(conn, aws.StringValueSlice(userNames)); err != nil {
		return fmt.Errorf("invalid user names: %w", err)
	}

	if err := addUsersToGroup(conn, userNames, projectName, groupArn); err != nil {
		return fmt.Errorf("error creating IAM group membership (%s): %w", groupArn, err)
	}

	d.SetId(d.Get("name").(string))
	return resourceGroupMembershipRead(d, meta)
}

func resourceGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	groupArn := d.Get("group_arn").(string)
	projectName := d.Get("project").(string)
	userNames := d.Get("users").(*schema.Set)

	group, users, err := FindGroupByArn(conn, groupArn)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM group membership (%s) not found, removing from state", groupArn)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM group membership (%s): %w", groupArn, err)
	}

	var userNamesToSet []string
	for _, user := range users {
		if !userNames.Contains(aws.StringValue(user.UserName)) {
			continue
		}

		if aws.StringValue(group.Type) == iam.GrantsTypeTypeProject {
			for _, project := range user.Projects {
				if projectName == aws.StringValue(project.ProjectName) {
					userNamesToSet = append(userNamesToSet, aws.StringValue(user.UserName))
				}
			}
		} else {
			userNamesToSet = append(userNamesToSet, aws.StringValue(user.UserName))
		}
	}

	d.Set("users", userNamesToSet)

	return nil
}

func resourceGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	if d.HasChange("users") {
		groupArn := d.Get("group_arn").(string)
		projectName := d.Get("project").(string)

		o, n := d.GetChange("users")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := flex.ExpandStringSet(os.Difference(ns))
		add := flex.ExpandStringSet(ns.Difference(os))

		if err := validateIfUsersExist(conn, aws.StringValueSlice(add)); err != nil {
			return fmt.Errorf("invalid user names: %w", err)
		}

		if err := removeUsersFromGroup(conn, remove, projectName, groupArn); err != nil {
			return fmt.Errorf("error deleting IAM group membership (%s): %w", groupArn, err)
		}

		if err := addUsersToGroup(conn, add, projectName, groupArn); err != nil {
			return fmt.Errorf("error creating IAM group membership (%s): %w", groupArn, err)
		}
	}

	return resourceGroupMembershipRead(d, meta)
}

func resourceGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	groupArn := d.Get("group_arn").(string)
	projectName := d.Get("project").(string)
	userNames := flex.ExpandStringSet(d.Get("users").(*schema.Set))

	err := removeUsersFromGroup(conn, userNames, projectName, groupArn)

	if tfawserr.ErrCodeEquals(err, GroupNotFoundCode, ProjectNotFoundCode) {
		log.Printf("[WARN] IAM group membership (%s) not found, removing from state", groupArn)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM group membership (%s): %w", groupArn, err)
	}

	return nil
}

func validateIfUsersExist(conn *iam.IAM, userNames []string) error {
	availableUsers, err := FindUsers(conn, "", "")

	if err != nil {
		return fmt.Errorf("error fetching IAM users %w", err)
	}

	var availableUserNames []string
	for _, user := range availableUsers {
		availableUserNames = append(availableUserNames, aws.StringValue(user.UserName))
	}

	var invalidUserNames []string
	for _, userName := range userNames {
		if !slices.Contains(availableUserNames, userName) {
			invalidUserNames = append(invalidUserNames, userName)
		}
	}

	if len(invalidUserNames) > 0 {
		formattedUserNames := strings.Join(invalidUserNames, "\n")
		return fmt.Errorf("Specified IAM users aren't found:\n%s", formattedUserNames)
	}

	return nil
}

func removeUsersFromGroup(conn *iam.IAM, users []*string, projectName, groupArn string) error {
	for _, userName := range users {
		input := &iam.RemoveUserFromGroupInput{
			UserName: userName,
			GroupArn: aws.String(groupArn),
		}

		if projectName != "" {
			input.ProjectName = aws.String(projectName)
		}

		if _, err := conn.RemoveUserFromGroup(input); err != nil {
			if tfawserr.ErrCodeEquals(err, UserNotFoundCode) {
				continue
			}

			return err
		}
	}
	return nil
}

func addUsersToGroup(conn *iam.IAM, users []*string, projectName, groupArn string) error {
	for _, userName := range users {
		input := &iam.AddUserToGroupInput{
			UserName: userName,
			GroupArn: aws.String(groupArn),
		}

		if projectName != "" {
			input.ProjectName = aws.String(projectName)
		}

		if _, err := conn.AddUserToGroup(input); err != nil {
			return err
		}
	}
	return nil
}
