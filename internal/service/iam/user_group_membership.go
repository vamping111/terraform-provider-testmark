package iam

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupMembershipCreate,
		Read:   resourceUserGroupMembershipRead,
		Update: resourceUserGroupMembershipUpdate,
		Delete: resourceUserGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserGroupMembershipImport,
		},

		Schema: map[string]*schema.Schema{
			"group_arns": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	user := d.Get("user").(string)
	project := d.Get("project").(string)
	groupArns := flex.ExpandStringSet(d.Get("group_arns").(*schema.Set))

	var groupType string
	if project == "" {
		groupType = iam.GrantsTypeTypeGlobal
	} else {
		groupType = iam.GrantsTypeTypeProject
	}

	if err := validateIfGroupsExist(conn, aws.StringValueSlice(groupArns), groupType); err != nil {
		return fmt.Errorf("invalid group ARNs: %w", err)
	}

	if err := addUserToGroups(conn, user, project, groupArns); err != nil {
		return fmt.Errorf("error creating IAM user group membership (%s): %w", user, err)
	}

	// lintignore:R015 // Allow legacy unstable ID usage in managed resource
	d.SetId(resource.UniqueId())

	return resourceUserGroupMembershipRead(d, meta)
}

func resourceUserGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	user := d.Get("user").(string)
	project := d.Get("project").(string)
	groupArns := d.Get("group_arns").(*schema.Set)

	var groups []*iam.Group
	var err error
	if project == "" {
		groups, err = FindUserGlobalGroups(conn, user)
	} else {
		groups, err = FindUserProjectGroups(conn, user, project)
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, UserNotFoundCode) {
		log.Printf("[WARN] IAM user group membership (%s) not found, removing from state", user)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM user group membership (%s): %w", user, err)
	}

	var groupArnsToSet []string
	for _, group := range groups {
		if groupArns.Contains(aws.StringValue(group.GroupArn)) {
			groupArnsToSet = append(groupArnsToSet, aws.StringValue(group.GroupArn))
		}
	}

	d.Set("group_arns", groupArnsToSet)

	return nil
}

func resourceUserGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	if d.HasChange("group_arns") {
		user := d.Get("user").(string)
		project := d.Get("project").(string)

		o, n := d.GetChange("group_arns")
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

		var groupType string
		if project == "" {
			groupType = iam.GrantsTypeTypeGlobal
		} else {
			groupType = iam.GrantsTypeTypeProject
		}

		if err := validateIfGroupsExist(conn, aws.StringValueSlice(add), groupType); err != nil {
			return fmt.Errorf("invalid group ARNs: %w", err)
		}

		if err := removeUserFromGroups(conn, user, project, remove); err != nil {
			return fmt.Errorf("error deleting IAM user group membership (%s): %w", user, err)
		}

		if err := addUserToGroups(conn, user, project, add); err != nil {
			return fmt.Errorf("error creating IAM user group membership (%s): %w", user, err)
		}
	}

	return resourceUserGroupMembershipRead(d, meta)
}

func resourceUserGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	user := d.Get("user").(string)
	project := d.Get("project").(string)
	groupArns := flex.ExpandStringSet(d.Get("group_arns").(*schema.Set))

	err := removeUserFromGroups(conn, user, project, groupArns)

	if tfawserr.ErrCodeEquals(err, UserNotFoundCode, ProjectNotFoundCode) {
		log.Printf("[WARN] IAM user group membership (%s) not found, removing from state", user)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM user group membership (%s): %w", user, err)
	}

	return nil
}

func validateIfGroupsExist(conn *iam.IAM, groupArns []string, groupType string) error {
	availableGroups, err := FindGroups(conn, "", groupType)

	if err != nil {
		return fmt.Errorf("error fetching IAM groups %w", err)
	}

	var availableGroupArns []string
	for _, group := range availableGroups {
		availableGroupArns = append(availableGroupArns, aws.StringValue(group.GroupArn))
	}

	var invalidGroupArns []string
	for _, groupArn := range groupArns {
		if !slices.Contains(availableGroupArns, groupArn) {
			invalidGroupArns = append(invalidGroupArns, groupArn)
		}
	}

	if len(invalidGroupArns) > 0 {
		formattedArns := strings.Join(invalidGroupArns, "\n")
		return fmt.Errorf("Specified IAM groups with type %q aren't found:\n%s", groupType, formattedArns)
	}

	return nil
}

func removeUserFromGroups(conn *iam.IAM, userName, projectName string, groupArns []*string) error {
	for _, groupArn := range groupArns {
		input := &iam.RemoveUserFromGroupInput{
			UserName: aws.String(userName),
			GroupArn: groupArn,
		}

		if projectName != "" {
			input.ProjectName = aws.String(projectName)
		}

		if _, err := conn.RemoveUserFromGroup(input); err != nil {
			if tfawserr.ErrCodeEquals(err, GroupNotFoundCode) {
				continue
			}

			return err
		}
	}

	return nil
}

func addUserToGroups(conn *iam.IAM, userName string, projectName string, groupArns []*string) error {
	for _, groupArn := range groupArns {
		input := &iam.AddUserToGroupInput{
			UserName: aws.String(userName),
			GroupArn: groupArn,
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

func resourceUserGroupMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), "#")
	if len(idParts) < 2 {
		return nil, fmt.Errorf(
			"unexpected format of ID (%q), expected <user-name>[#project-name]#<group-arn1>#... ", d.Id(),
		)
	}

	userName := idParts[0]

	var projectName string
	var groupArns []string

	if arn.IsARN(idParts[1]) {
		groupArns = idParts[1:]
	} else {
		if len(idParts) == 2 {
			return nil, fmt.Errorf(
				"unexpected format of ID (%q): no IAM group ARNs found, "+
					"expected <user-name>[#project-name]#<group-arn1>#... ", d.Id(),
			)
		}

		projectName = idParts[1]
		groupArns = idParts[2:]
	}

	d.Set("user", userName)
	d.Set("project", projectName)
	d.Set("group_arns", groupArns)

	// lintignore:R015 // Allow legacy unstable ID usage in managed resource
	d.SetId(resource.UniqueId())

	return []*schema.ResourceData{d}, nil
}
