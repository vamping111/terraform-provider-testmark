package backup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourcePlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePlanRead,

		Schema: map[string]*schema.Schema{
			"plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePlanRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).BackupConn

	id := d.Get("plan_id").(string)

	resp, err := conn.GetBackupPlan(&backup.GetBackupPlanInput{
		BackupPlanId: aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("Error getting Backup Plan: %w", err)
	}

	d.SetId(aws.StringValue(resp.BackupPlanId))
	d.Set("arn", resp.BackupPlanArn)
	d.Set("name", resp.BackupPlan.BackupPlanName)
	d.Set("version", resp.VersionId)

	return nil
}
