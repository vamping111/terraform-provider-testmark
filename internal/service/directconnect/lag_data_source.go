package directconnect

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceLag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLagRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceLagRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn
	connEC2 := meta.(*conns.AWSClient).EC2Conn

	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	var lags []*directconnect.Lag
	input := &directconnect.DescribeLagsInput{}
	name := d.Get("name").(string)

	// DescribeLags is not paginated.
	output, err := conn.DescribeLags(input)

	if err != nil {
		return fmt.Errorf("error reading Direct Connect LAGs: %w", err)
	}

	for _, lag := range output.Lags {
		if aws.StringValue(lag.LagName) == name {
			lags = append(lags, lag)
		}
	}

	switch count := len(lags); count {
	case 0:
		return fmt.Errorf("no matching Direct Connect LAG found")
	case 1:
	default:
		return fmt.Errorf("%d Direct Connect LAGs matched; use additional constraints to reduce matches to a single Direct Connect LAG", count)
	}

	lag := lags[0]

	d.SetId(aws.StringValue(lag.LagId))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    aws.StringValue(lag.Region),
		Service:   "directconnect",
		AccountID: aws.StringValue(lag.OwnerAccount),
		Resource:  fmt.Sprintf("dxlag/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("aws_device", lag.AwsDeviceV2)
	d.Set("bandwidth", lag.ConnectionsBandwidth)
	d.Set("location", lag.Location)
	d.Set("name", lag.LagName)
	d.Set("owner_account_id", lag.OwnerAccount)
	d.Set("provider_name", lag.ProviderName)

	// FIXME: Use directconnect.ListTags after DescribeTags is supported in Direct Connect API.
	tags, err := ec2.ListTags(connEC2, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for Direct Connect LAG (%s): %w", arn, err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
