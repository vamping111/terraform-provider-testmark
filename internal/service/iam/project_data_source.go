package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	name := d.Get("name").(string)
	project, err := FindProjectByName(conn, name)

	if err != nil {
		return fmt.Errorf("error reading IAM project: %w", err)
	}

	d.SetId(aws.StringValue(project.ProjectName))

	d.Set("arn", project.ProjectArn)

	if project.CreateDate != nil {
		d.Set("create_date", aws.TimeValue(project.CreateDate).Format(time.RFC3339))
	} else {
		d.Set("create_date", nil)
	}

	d.Set("display_name", project.DisplayName)
	d.Set("name", project.ProjectName)
	d.Set("project_id", project.ProjectId)
	d.Set("s3_email", project.S3Email)
	d.Set("state", project.State)

	return nil
}
