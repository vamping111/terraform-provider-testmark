package iam

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
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
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 40),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]*$`),
						"name must start with a Latin letter "+
							"and can only contain Latin letters, numbers, underscores (_), periods (.) and hyphens (-)",
					),
				),
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

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Get("name").(string)

	input := &iam.CreateProjectInput{
		ProjectName: aws.String(name),
	}

	if v, ok := d.GetOk("display_name"); ok {
		input.DisplayName = aws.String(v.(string))
	} else {
		input.DisplayName = aws.String(name)
	}

	log.Printf("[DEBUG] Creating IAM project: %s", input)
	output, err := conn.CreateProject(input)

	if err != nil {
		return fmt.Errorf("failed creating IAM project (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.Project.ProjectName))

	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Id()

	project, err := FindProjectByName(conn, name)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IAM project (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM project: %w", err)
	}

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

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Id()

	if d.HasChange("display_name") {
		input := &iam.UpdateProjectInput{
			ProjectName: aws.String(name),
			DisplayName: aws.String(d.Get("display_name").(string)),
		}

		log.Printf("[DEBUG] Modifying IAM project: %s", input)
		_, err := conn.UpdateProject(input)

		if err != nil {
			return fmt.Errorf("error modifying IAM project (%s): %w", name, err)
		}
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	name := d.Id()

	input := &iam.DeleteProjectInput{
		ProjectName: aws.String(name),
	}

	log.Printf("[DEBUG] Deleting IAM project: %s", input)
	_, err := conn.DeleteProject(input)

	if tfawserr.ErrCodeEquals(err, ProjectNotFoundCode) {
		log.Printf("[WARN] IAM project (%s) not found, removing from state", name)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting IAM project %s: %w", name, err)
	}

	return nil
}
