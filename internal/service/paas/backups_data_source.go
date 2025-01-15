package paas

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func DataSourceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupsRead,

		Schema: map[string]*schema.Schema{
			"backup_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(services.ServiceClassValues(), false),
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(services.ServiceTypeValues(), false),
			},
		},
	}
}

func dataSourceBackupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceClass := d.Get("service_class").(string)
	serviceId := d.Get("service_id").(string)
	serviceType := d.Get("service_type").(string)

	backups, err := FindBackups(conn, serviceClass, serviceId, serviceType)

	if err != nil {
		return diag.Errorf("error reading PaaS Service Backups: %s", err)
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	var ids []string
	for _, backup := range backups {
		ids = append(ids, aws.StringValue(backup.Id))
	}

	d.Set("backup_ids", ids)

	return nil
}
