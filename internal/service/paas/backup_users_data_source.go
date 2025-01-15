package paas

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceBackupUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupUsersRead,

		Schema: map[string]*schema.Schema{
			"active_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBackupUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	users, err := FindBackupUsers(conn)

	if err != nil {
		return diag.Errorf("error listing Backup Users: %s", err)
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	var filtered []*paas.BackupUser
	if v, ok := d.Get("active_only").(bool); ok && v {
		for _, user := range users {
			if aws.BoolValue(user.Enabled) {
				filtered = append(filtered, user)
			}
		}
	} else {
		filtered = users[:]
	}

	d.Set("users", flattenBackupUsers(filtered))

	return nil
}

func flattenBackupUsers(users []*paas.BackupUser) []map[string]interface{} {
	if users == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, user := range users {
		if user == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := user.Email; v != nil {
			tfMap["email"] = v
		}

		if v := user.Enabled; v != nil {
			tfMap["enabled"] = v
		}

		if v := user.Id; v != nil {
			tfMap["id"] = v
		}

		if v := user.Login; v != nil {
			tfMap["login"] = v
		}

		if v := user.Name; v != nil {
			tfMap["name"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}
