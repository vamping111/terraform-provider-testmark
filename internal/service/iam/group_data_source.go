package iam

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

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
			"group_id": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	arn := d.Get("arn").(string)
	name := d.Get("name").(string)

	if arn == "" {
		groups, err := FindGroups(conn, name, "")

		if err != nil {
			return fmt.Errorf("error reading IAM group (%s): %w", GroupSearchDetails(arn, name), err)
		}

		if len(groups) == 0 {
			return fmt.Errorf("no IAM group found matching criteria (%s); try different search", GroupSearchDetails(arn, name))
		}

		if len(groups) > 1 {
			return fmt.Errorf("multiple IAM groups found matching criteria (%s); try different search", GroupSearchDetails(arn, name))
		}

		arn = aws.StringValue(groups[0].GroupArn)
	}

	group, users, err := FindGroupByArn(conn, arn)

	if err != nil {
		return fmt.Errorf("error reading IAM group (%s): %w", arn, err)
	}

	d.SetId(aws.StringValue(group.GroupArn))

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

	if err := d.Set("users", dataSourceGroupUsersRead(users)); err != nil {
		return fmt.Errorf("error setting users from IAM group (%s): %w", arn, err)
	}

	return nil
}

func dataSourceGroupUsersRead(iamUsers []*iam.User) []map[string]interface{} {
	users := make([]map[string]interface{}, 0, len(iamUsers))
	for _, i := range iamUsers {
		u := make(map[string]interface{})
		u["arn"] = aws.StringValue(i.UserArn)
		u["user_id"] = aws.StringValue(i.UserId)
		u["user_name"] = aws.StringValue(i.UserName)
		u["path"] = aws.StringValue(i.Path)
		users = append(users, u)
	}
	return users
}

func GroupSearchDetails(arn, name string) string {
	var groupDetails []string
	if arn != "" {
		groupDetails = append(groupDetails, fmt.Sprintf("ARN: %s", arn))
	}
	if name != "" {
		groupDetails = append(groupDetails, fmt.Sprintf("Name: %s", name))
	}

	return strings.Join(groupDetails, ", ")
}
