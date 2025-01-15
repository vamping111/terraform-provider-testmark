package directconnect

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceTransitVirtualInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransitVirtualInterfaceCreate,
		Read:   resourceTransitVirtualInterfaceRead,
		Update: resourceTransitVirtualInterfaceUpdate,
		Delete: resourceTransitVirtualInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTransitVirtualInterfaceImport,
		},

		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  directconnect.AddressFamilyIpv4,
				ValidateFunc: validation.StringInSlice([]string{
					directconnect.AddressFamilyIpv4,
				}, false),
			},
			"amazon_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"amazon_side_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"bgp_auth_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"customer_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dx_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"jumbo_frame_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1500, 8500}),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sitelink_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceTransitVirtualInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn
	connEC2 := meta.(*conns.AWSClient).EC2Conn

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	req := &directconnect.CreateTransitVirtualInterfaceInput{
		ConnectionId: aws.String(d.Get("connection_id").(string)),
		NewTransitVirtualInterface: &directconnect.NewTransitVirtualInterface{
			AddressFamily:          aws.String(d.Get("address_family").(string)),
			Asn:                    aws.Int64(int64(d.Get("bgp_asn").(int))),
			DirectConnectGatewayId: aws.String(d.Get("dx_gateway_id").(string)),
			EnableSiteLink:         aws.Bool(d.Get("sitelink_enabled").(bool)),
			Mtu:                    aws.Int64(int64(d.Get("mtu").(int))),
			VirtualInterfaceName:   aws.String(d.Get("name").(string)),
			Vlan:                   aws.Int64(int64(d.Get("vlan").(int))),
		},
	}
	if v, ok := d.GetOk("amazon_address"); ok {
		req.NewTransitVirtualInterface.AmazonAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("bgp_auth_key"); ok {
		req.NewTransitVirtualInterface.AuthKey = aws.String(v.(string))
	}
	if v, ok := d.GetOk("customer_address"); ok {
		req.NewTransitVirtualInterface.CustomerAddress = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Direct Connect transit virtual interface: %s", req)
	resp, err := conn.CreateTransitVirtualInterface(req)
	if err != nil {
		return fmt.Errorf("error creating Direct Connect transit virtual interface: %s", err)
	}

	d.SetId(aws.StringValue(resp.VirtualInterface.VirtualInterfaceId))

	// FIXME: Provide tags in CreateTransitVirtualInterface request after Direct Connect API support.
	if len(tags) > 0 {
		if err := ec2.CreateTags(connEC2, d.Id(), tags.IgnoreAWS()); err != nil {
			return fmt.Errorf("error creating tags for Direct Connect transit virtual interface: %s", err)
		}
	}

	if err := dxTransitVirtualInterfaceWaitUntilAvailable(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	return resourceTransitVirtualInterfaceRead(d, meta)
}

func resourceTransitVirtualInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn
	connEC2 := meta.(*conns.AWSClient).EC2Conn

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	vif, err := dxVirtualInterfaceRead(d.Id(), conn)
	if err != nil {
		return err
	}
	if vif == nil {
		log.Printf("[WARN] Direct Connect transit virtual interface (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("address_family", vif.AddressFamily)
	d.Set("amazon_address", vif.AmazonAddress)
	d.Set("amazon_side_asn", strconv.FormatInt(aws.Int64Value(vif.AmazonSideAsn), 10))
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "directconnect",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("dxvif/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("aws_device", vif.AwsDeviceV2)
	d.Set("bgp_asn", vif.Asn)
	d.Set("bgp_auth_key", vif.AuthKey)
	d.Set("connection_id", vif.ConnectionId)
	d.Set("customer_address", vif.CustomerAddress)
	d.Set("dx_gateway_id", vif.DirectConnectGatewayId)
	d.Set("jumbo_frame_capable", vif.JumboFrameCapable)
	d.Set("mtu", vif.Mtu)
	d.Set("name", vif.VirtualInterfaceName)
	d.Set("sitelink_enabled", vif.SiteLinkEnabled)
	d.Set("vlan", vif.Vlan)

	// FIXME: Use directconnect.ListTags after DescribeTags is supported in Direct Connect API.
	tags, err := ec2.ListTags(connEC2, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for Direct Connect transit virtual interface (%s): %s", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceTransitVirtualInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := dxVirtualInterfaceUpdate(d, meta); err != nil {
		return err
	}

	if err := dxTransitVirtualInterfaceWaitUntilAvailable(meta.(*conns.AWSClient).DirectConnectConn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTransitVirtualInterfaceRead(d, meta)
}

func resourceTransitVirtualInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	return dxVirtualInterfaceDelete(d, meta)
}

func resourceTransitVirtualInterfaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	vif, err := dxVirtualInterfaceRead(d.Id(), conn)
	if err != nil {
		return nil, err
	}
	if vif == nil {
		return nil, fmt.Errorf("virtual interface (%s) not found", d.Id())
	}

	if vifType := aws.StringValue(vif.VirtualInterfaceType); vifType != "transit" {
		return nil, fmt.Errorf("virtual interface (%s) has incorrect type: %s", d.Id(), vifType)
	}

	return []*schema.ResourceData{d}, nil
}

func dxTransitVirtualInterfaceWaitUntilAvailable(conn *directconnect.DirectConnect, vifId string, timeout time.Duration) error {
	return dxVirtualInterfaceWaitUntilAvailable(
		conn,
		vifId,
		timeout,
		[]string{
			directconnect.VirtualInterfaceStatePending,
		},
		[]string{
			directconnect.VirtualInterfaceStateAvailable,
			directconnect.VirtualInterfaceStateDown,
		})
}
