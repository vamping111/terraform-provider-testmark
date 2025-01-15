package ec2

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceSecurityGroup() *schema.Resource {
	//lintignore:R011
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  SecurityGroupMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validation.StringLenBetween(0, 255),
			},

			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(0, 100),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ingress": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"protocol": {
							Type:      schema.TypeString,
							Required:  true,
							StateFunc: ProtocolStateFunc,
						},

						"cidr_blocks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: verify.ValidCIDRNetworkAddress,
							},
						},

						"ipv6_cidr_blocks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: verify.ValidCIDRNetworkAddress,
							},
						},

						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"security_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"self": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validSecurityGroupRuleDescription,
						},
					},
				},
				Set: SecurityGroupRuleHash,
			},

			"egress": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"protocol": {
							Type:      schema.TypeString,
							Required:  true,
							StateFunc: ProtocolStateFunc,
						},

						"cidr_blocks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: verify.ValidCIDRNetworkAddress,
							},
						},

						"ipv6_cidr_blocks": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: verify.ValidCIDRNetworkAddress,
							},
						},

						"prefix_list_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"security_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"self": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validSecurityGroupRuleDescription,
						},
					},
				},
				Set: SecurityGroupRuleHash,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),

			"revoke_rules_on_delete": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	groupName := create.Name(d.Get("name").(string), d.Get("name_prefix").(string))
	input := &ec2.CreateSecurityGroupInput{
		GroupName: aws.String(groupName),
	}

	if v := d.Get("description"); v != nil {
		input.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.TagSpecifications = ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeSecurityGroup)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		input.VpcId = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Security Group: %s", input)
	output, err := conn.CreateSecurityGroup(input)

	if err != nil {
		return fmt.Errorf("error creating Security Group (%s): %w", groupName, err)
	}

	d.SetId(aws.StringValue(output.GroupId))

	// Wait for the security group to truly exist
	group, err := WaitSecurityGroupCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("error waiting for Security Group (%s) create: %w", d.Id(), err)
	}

	// AWS defaults all Security Groups to have an ALLOW ALL egress rule. Here we
	// revoke that rule, so users don't unknowingly have/use it.
	if aws.StringValue(group.VpcId) != "" {
		log.Printf("[DEBUG] Revoking default egress rule for Security Group for %s", d.Id())

		input := &ec2.RevokeSecurityGroupEgressInput{
			GroupId: output.GroupId,
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort: aws.Int64(0),
					ToPort:   aws.Int64(0),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
					IpProtocol: aws.String("-1"),
				},
			},
		}

		if _, err = conn.RevokeSecurityGroupEgress(input); err != nil {
			return fmt.Errorf("Error revoking default egress rule for Security Group (%s): %w", d.Id(), err)
		}

		log.Printf("[DEBUG] Revoking default IPv6 egress rule for Security Group for %s", d.Id())
		input = &ec2.RevokeSecurityGroupEgressInput{
			GroupId: output.GroupId,
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort: aws.Int64(0),
					ToPort:   aws.Int64(0),
					Ipv6Ranges: []*ec2.Ipv6Range{
						{
							CidrIpv6: aws.String("::/0"),
						},
					},
					IpProtocol: aws.String("-1"),
				},
			},
		}

		_, err = conn.RevokeSecurityGroupEgress(input)
		if err != nil {
			//If we have a NotFound or InvalidParameterValue, then we are trying to remove the default IPv6 egress of a non-IPv6
			//enabled SG
			if !tfawserr.ErrCodeEquals(err, ErrCodeInvalidPermissionNotFound) && !tfawserr.ErrMessageContains(err, ErrCodeInvalidParameterValue, "remote-ipv6-range") {
				return fmt.Errorf("Error revoking default IPv6 egress rule for Security Group (%s): %w", d.Id(), err)
			}
		}

	}

	return resourceSecurityGroupUpdate(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	sg, err := FindSecurityGroupByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Security Group %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Security Group (%s): %w", d.Id(), err)
	}

	remoteIngressRules := SecurityGroupIPPermGather(d.Id(), sg.IpPermissions, sg.OwnerId)
	remoteEgressRules := SecurityGroupIPPermGather(d.Id(), sg.IpPermissionsEgress, sg.OwnerId)

	localIngressRules := d.Get("ingress").(*schema.Set).List()
	localEgressRules := d.Get("egress").(*schema.Set).List()

	// Loop through the local state of rules, doing a match against the remote
	// ruleSet we built above.
	ingressRules := MatchRules("ingress", localIngressRules, remoteIngressRules)
	egressRules := MatchRules("egress", localEgressRules, remoteEgressRules)

	sgArn := arn.ARN{
		AccountID: aws.StringValue(sg.OwnerId),
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Resource:  fmt.Sprintf("security-group/%s", aws.StringValue(sg.GroupId)),
		Service:   ec2.ServiceName,
	}

	d.Set("arn", sgArn.String())
	d.Set("description", sg.Description)
	d.Set("name", sg.GroupName)
	d.Set("name_prefix", create.NamePrefixFromName(aws.StringValue(sg.GroupName)))
	d.Set("owner_id", sg.OwnerId)
	d.Set("vpc_id", sg.VpcId)

	if err := d.Set("ingress", ingressRules); err != nil {
		return fmt.Errorf("error setting ingress: %w", err)
	}

	if err := d.Set("egress", egressRules); err != nil {
		return fmt.Errorf("error setting egress: %w", err)
	}

	tags := KeyValueTags(sg.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	group, err := FindSecurityGroupByID(conn, d.Id())
	if err != nil {
		return fmt.Errorf("error updating Security Group (%s): %w", d.Id(), err)
	}

	if d.HasChange("description") && !d.IsNewResource() {
		_, err := conn.ModifySecurityGroupAttribute(&ec2.ModifySecurityGroupAttributeInput{
			Description: &ec2.AttributeValue{
				Value: aws.String(d.Get("description").(string)),
			},
			GroupId: aws.String(d.Id()),
		})

		if err != nil {
			return fmt.Errorf("error updating Security Group (%s) description: %w", d.Id(), err)
		}
	}

	err = resourceSecurityGroupUpdateRules(d, "ingress", meta, group)
	if err != nil {
		return fmt.Errorf("error updating Security Group (%s): %w", d.Id(), err)
	}

	if d.Get("vpc_id") != nil {
		err = resourceSecurityGroupUpdateRules(d, "egress", meta, group)
		if err != nil {
			return fmt.Errorf("error updating Security Group (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") && !d.IsNewResource() {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating Security Group (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	log.Printf("[DEBUG] Security Group destroy: %v", d.Id())

	if err := deleteLingeringLambdaENIs(conn, "group-id", d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error deleting Lambda ENIs using Security Group (%s): %w", d.Id(), err)
	}

	// conditionally revoke rules first before attempting to delete the group
	if v := d.Get("revoke_rules_on_delete").(bool); v {
		if err := forceRevokeSecurityGroupRules(conn, d); err != nil {
			return err
		}
	}
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(d.Id()),
	}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteSecurityGroup(input)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, "InvalidGroup.NotFound") {
				return nil
			}

			// If it is a dependency violation, we want to retry
			if tfawserr.ErrMessageContains(err, "DependencyViolation", "has a dependent object") {
				return resource.RetryableError(err)
			}

			if tfawserr.ErrCodeEquals(err, "DependencyViolation") {
				return resource.RetryableError(err)
			}

			if tfawserr.ErrCodeEquals(err, "InvalidGroup.InUse") {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteSecurityGroup(input)
		if tfawserr.ErrCodeEquals(err, "InvalidGroup.NotFound") {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("Error deleting security group: %w", err)
	}
	return nil
}

// Revoke all ingress/egress rules that a Security Group has
func forceRevokeSecurityGroupRules(conn *ec2.EC2, d *schema.ResourceData) error {
	group, err := FindSecurityGroupByID(conn, d.Id())
	if err != nil {
		return err
	}

	if len(group.IpPermissions) > 0 {
		req := &ec2.RevokeSecurityGroupIngressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissions,
		}
		if aws.StringValue(group.VpcId) == "" {
			req.GroupId = nil
			req.GroupName = group.GroupName
		}
		_, err = conn.RevokeSecurityGroupIngress(req)

		if err != nil {
			return fmt.Errorf("error revoking Security Group (%s) rules: %w", aws.StringValue(group.GroupId), err)
		}
	}

	if len(group.IpPermissionsEgress) > 0 {
		req := &ec2.RevokeSecurityGroupEgressInput{
			GroupId:       group.GroupId,
			IpPermissions: group.IpPermissionsEgress,
		}
		_, err = conn.RevokeSecurityGroupEgress(req)

		if err != nil {
			return fmt.Errorf("error revoking Security Group (%s) rules: %w", aws.StringValue(group.GroupId), err)
		}
	}

	return nil
}

func SecurityGroupRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["from_port"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["to_port"].(int)))
	p := ProtocolForValue(m["protocol"].(string))
	buf.WriteString(fmt.Sprintf("%s-", p))
	buf.WriteString(fmt.Sprintf("%t-", m["self"].(bool)))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["cidr_blocks"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["ipv6_cidr_blocks"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["prefix_list_ids"]; ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if v, ok := m["security_groups"]; ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if m["description"].(string) != "" {
		buf.WriteString(fmt.Sprintf("%s-", m["description"].(string)))
	}

	return create.StringHashcode(buf.String())
}

func SecurityGroupIPPermGather(groupId string, permissions []*ec2.IpPermission, ownerId *string) []map[string]interface{} {
	ruleMap := make(map[string]map[string]interface{})
	for _, perm := range permissions {
		if len(perm.IpRanges) > 0 {
			for _, ip := range perm.IpRanges {
				desc := aws.StringValue(ip.Description)

				rule := initSecurityGroupRule(ruleMap, perm, desc)

				raw, ok := rule["cidr_blocks"]
				if !ok {
					raw = make([]string, 0)
				}
				list := raw.([]string)

				rule["cidr_blocks"] = append(list, *ip.CidrIp)
			}
		}

		if len(perm.Ipv6Ranges) > 0 {
			for _, ip := range perm.Ipv6Ranges {
				desc := aws.StringValue(ip.Description)

				rule := initSecurityGroupRule(ruleMap, perm, desc)

				raw, ok := rule["ipv6_cidr_blocks"]
				if !ok {
					raw = make([]string, 0)
				}
				list := raw.([]string)

				rule["ipv6_cidr_blocks"] = append(list, *ip.CidrIpv6)
			}
		}

		if len(perm.PrefixListIds) > 0 {
			for _, pl := range perm.PrefixListIds {
				desc := aws.StringValue(pl.Description)

				rule := initSecurityGroupRule(ruleMap, perm, desc)

				raw, ok := rule["prefix_list_ids"]
				if !ok {
					raw = make([]string, 0)
				}
				list := raw.([]string)

				rule["prefix_list_ids"] = append(list, *pl.PrefixListId)
			}
		}

		groups := FlattenSecurityGroups(perm.UserIdGroupPairs, ownerId)
		if len(groups) > 0 {
			for _, g := range groups {
				desc := aws.StringValue(g.Description)

				rule := initSecurityGroupRule(ruleMap, perm, desc)

				if aws.StringValue(g.GroupId) == groupId {
					rule["self"] = true
					continue
				}

				raw, ok := rule["security_groups"]
				if !ok {
					raw = schema.NewSet(schema.HashString, nil)
				}
				list := raw.(*schema.Set)

				if g.GroupName != nil {
					list.Add(*g.GroupName)
				} else {
					list.Add(*g.GroupId)
				}
				rule["security_groups"] = list
			}
		}

	}

	rules := make([]map[string]interface{}, 0, len(ruleMap))
	for _, m := range ruleMap {
		rules = append(rules, m)
	}

	return rules
}

func resourceSecurityGroupUpdateRules(
	d *schema.ResourceData, ruleset string,
	meta interface{}, group *ec2.SecurityGroup) error {

	if d.HasChange(ruleset) {
		o, n := d.GetChange(ruleset)
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := SecurityGroupExpandRules(o.(*schema.Set))
		ns := SecurityGroupExpandRules(n.(*schema.Set))

		remove, err := ExpandIPPerms(group, SecurityGroupCollapseRules(ruleset, os.Difference(ns).List()))
		if err != nil {
			return err
		}
		add, err := ExpandIPPerms(group, SecurityGroupCollapseRules(ruleset, ns.Difference(os).List()))
		if err != nil {
			return err
		}

		// TODO: We need to handle partial state better in the in-between
		// in this update.

		// TODO: It'd be nicer to authorize before removing, but then we have
		// to deal with complicated unrolling to get individual CIDR blocks
		// to avoid authorizing already authorized sources. Removing before
		// adding is easier here, and Terraform should be fast enough to
		// not have service issues.

		if len(remove) > 0 || len(add) > 0 {
			conn := meta.(*conns.AWSClient).EC2Conn

			var err error
			if len(remove) > 0 {
				log.Printf("[DEBUG] Revoking security group %#v %s rule: %#v",
					group, ruleset, remove)

				if ruleset == "egress" {
					req := &ec2.RevokeSecurityGroupEgressInput{
						GroupId:       group.GroupId,
						IpPermissions: remove,
					}
					_, err = conn.RevokeSecurityGroupEgress(req)
				} else {
					req := &ec2.RevokeSecurityGroupIngressInput{
						GroupId:       group.GroupId,
						IpPermissions: remove,
					}
					if aws.StringValue(group.VpcId) == "" {
						req.GroupId = nil
						req.GroupName = group.GroupName
					}
					_, err = conn.RevokeSecurityGroupIngress(req)
				}

				if err != nil {
					return fmt.Errorf("error revoking Security Group (%s) rules: %w", ruleset, err)
				}
			}

			if len(add) > 0 {
				log.Printf("[DEBUG] Authorizing security group %#v %s rule: %#v",
					group, ruleset, add)
				// Authorize the new rules
				if ruleset == "egress" {
					req := &ec2.AuthorizeSecurityGroupEgressInput{
						GroupId:       group.GroupId,
						IpPermissions: add,
					}
					_, err = conn.AuthorizeSecurityGroupEgress(req)
				} else {
					req := &ec2.AuthorizeSecurityGroupIngressInput{
						GroupId:       group.GroupId,
						IpPermissions: add,
					}
					if aws.StringValue(group.VpcId) == "" {
						req.GroupId = nil
						req.GroupName = group.GroupName
					}

					_, err = conn.AuthorizeSecurityGroupIngress(req)
				}

				if err != nil {
					return fmt.Errorf("error authorizing Security Group (%s) rules: %w", ruleset, err)
				}
			}
		}
	}
	return nil
}

// MatchRules receives the group id, type of rules, and the local / remote maps
// of rules. We iterate through the local set of rules trying to find a matching
// remote rule, which may be structured differently because of how AWS
// aggregates the rules under the to, from, and type.
//
// Matching rules are written to state, with their elements removed from the
// remote set
//
// If no match is found, we'll write the remote rule to state and let the graph
// sort things out
func MatchRules(rType string, local []interface{}, remote []map[string]interface{}) []map[string]interface{} {
	// For each local ip or security_group, we need to match against the remote
	// ruleSet until all ips or security_groups are found

	// saves represents the rules that have been identified to be saved to state,
	// in the appropriate d.Set("{ingress,egress}") call.
	var saves []map[string]interface{}
	for _, raw := range local {
		l := raw.(map[string]interface{})

		var selfVal bool
		if v, ok := l["self"]; ok {
			selfVal = v.(bool)
		}

		// matching against self is required to detect rules that only include self
		// as the rule. SecurityGroupIPPermGather parses the group out
		// and replaces it with self if it's ID is found
		localHash := idHash(rType, l["protocol"].(string), int64(l["to_port"].(int)), int64(l["from_port"].(int)), selfVal)

		// loop remote rules, looking for a matching hash
		for _, r := range remote {
			var remoteSelfVal bool
			if v, ok := r["self"]; ok {
				remoteSelfVal = v.(bool)
			}

			// hash this remote rule and compare it for a match consideration with the
			// local rule we're examining
			rHash := idHash(rType, r["protocol"].(string), r["to_port"].(int64), r["from_port"].(int64), remoteSelfVal)
			if rHash == localHash {
				var numExpectedCidrs, numExpectedIpv6Cidrs, numExpectedPrefixLists, numExpectedSGs, numRemoteCidrs, numRemoteIpv6Cidrs, numRemotePrefixLists, numRemoteSGs int
				var matchingCidrs []string
				var matchingIpv6Cidrs []string
				var matchingSGs []string
				var matchingPrefixLists []string

				// grab the local/remote cidr and sg groups, capturing the expected and
				// actual counts
				lcRaw, ok := l["cidr_blocks"]
				if ok {
					numExpectedCidrs = len(l["cidr_blocks"].([]interface{}))
				}
				liRaw, ok := l["ipv6_cidr_blocks"]
				if ok {
					numExpectedIpv6Cidrs = len(l["ipv6_cidr_blocks"].([]interface{}))
				}
				lpRaw, ok := l["prefix_list_ids"]
				if ok {
					numExpectedPrefixLists = len(l["prefix_list_ids"].([]interface{}))
				}
				lsRaw, ok := l["security_groups"]
				if ok {
					numExpectedSGs = len(l["security_groups"].(*schema.Set).List())
				}

				rcRaw, ok := r["cidr_blocks"]
				if ok {
					numRemoteCidrs = len(r["cidr_blocks"].([]string))
				}
				riRaw, ok := r["ipv6_cidr_blocks"]
				if ok {
					numRemoteIpv6Cidrs = len(r["ipv6_cidr_blocks"].([]string))
				}
				rpRaw, ok := r["prefix_list_ids"]
				if ok {
					numRemotePrefixLists = len(r["prefix_list_ids"].([]string))
				}

				rsRaw, ok := r["security_groups"]
				if ok {
					numRemoteSGs = len(r["security_groups"].(*schema.Set).List())
				}

				// check some early failures
				if numExpectedCidrs > numRemoteCidrs {
					log.Printf("[DEBUG] Local rule has more CIDR blocks, continuing (%d/%d)", numExpectedCidrs, numRemoteCidrs)
					continue
				}
				if numExpectedIpv6Cidrs > numRemoteIpv6Cidrs {
					log.Printf("[DEBUG] Local rule has more IPV6 CIDR blocks, continuing (%d/%d)", numExpectedIpv6Cidrs, numRemoteIpv6Cidrs)
					continue
				}
				if numExpectedPrefixLists > numRemotePrefixLists {
					log.Printf("[DEBUG] Local rule has more prefix lists, continuing (%d/%d)", numExpectedPrefixLists, numRemotePrefixLists)
					continue
				}
				if numExpectedSGs > numRemoteSGs {
					log.Printf("[DEBUG] Local rule has more Security Groups, continuing (%d/%d)", numExpectedSGs, numRemoteSGs)
					continue
				}

				// match CIDRs by converting both to sets, and using Set methods
				var localCidrs []interface{}
				if lcRaw != nil {
					localCidrs = lcRaw.([]interface{})
				}
				localCidrSet := schema.NewSet(schema.HashString, localCidrs)

				// remote cidrs are presented as a slice of strings, so we need to
				// reformat them into a slice of interfaces to be used in creating the
				// remote cidr set
				var remoteCidrs []string
				if rcRaw != nil {
					remoteCidrs = rcRaw.([]string)
				}
				// convert remote cidrs to a set, for easy comparisons
				var list []interface{}
				for _, s := range remoteCidrs {
					list = append(list, s)
				}
				remoteCidrSet := schema.NewSet(schema.HashString, list)

				// Build up a list of local cidrs that are found in the remote set
				for _, s := range localCidrSet.List() {
					if remoteCidrSet.Contains(s) {
						matchingCidrs = append(matchingCidrs, s.(string))
					}
				}

				//IPV6 CIDRs
				var localIpv6Cidrs []interface{}
				if liRaw != nil {
					localIpv6Cidrs = liRaw.([]interface{})
				}
				localIpv6CidrSet := schema.NewSet(schema.HashString, localIpv6Cidrs)

				var remoteIpv6Cidrs []string
				if riRaw != nil {
					remoteIpv6Cidrs = riRaw.([]string)
				}
				var listIpv6 []interface{}
				for _, s := range remoteIpv6Cidrs {
					listIpv6 = append(listIpv6, s)
				}
				remoteIpv6CidrSet := schema.NewSet(schema.HashString, listIpv6)

				for _, s := range localIpv6CidrSet.List() {
					if remoteIpv6CidrSet.Contains(s) {
						matchingIpv6Cidrs = append(matchingIpv6Cidrs, s.(string))
					}
				}

				// match prefix lists by converting both to sets, and using Set methods
				var localPrefixLists []interface{}
				if lpRaw != nil {
					localPrefixLists = lpRaw.([]interface{})
				}
				localPrefixListsSet := schema.NewSet(schema.HashString, localPrefixLists)

				// remote prefix lists are presented as a slice of strings, so we need to
				// reformat them into a slice of interfaces to be used in creating the
				// remote prefix list set
				var remotePrefixLists []string
				if rpRaw != nil {
					remotePrefixLists = rpRaw.([]string)
				}
				// convert remote prefix lists to a set, for easy comparison
				list = nil
				for _, s := range remotePrefixLists {
					list = append(list, s)
				}
				remotePrefixListsSet := schema.NewSet(schema.HashString, list)

				// Build up a list of local prefix lists that are found in the remote set
				for _, s := range localPrefixListsSet.List() {
					if remotePrefixListsSet.Contains(s) {
						matchingPrefixLists = append(matchingPrefixLists, s.(string))
					}
				}

				// match SGs. Both local and remote are already sets
				var localSGSet *schema.Set
				if lsRaw == nil {
					localSGSet = schema.NewSet(schema.HashString, nil)
				} else {
					localSGSet = lsRaw.(*schema.Set)
				}

				var remoteSGSet *schema.Set
				if rsRaw == nil {
					remoteSGSet = schema.NewSet(schema.HashString, nil)
				} else {
					remoteSGSet = rsRaw.(*schema.Set)
				}

				// Build up a list of local security groups that are found in the remote set
				for _, s := range localSGSet.List() {
					if remoteSGSet.Contains(s) {
						matchingSGs = append(matchingSGs, s.(string))
					}
				}

				// compare equalities for matches.
				// If we found the number of cidrs and number of sgs, we declare a
				// match, and then remove those elements from the remote rule, so that
				// this remote rule can still be considered by other local rules
				if numExpectedCidrs == len(matchingCidrs) {
					if numExpectedIpv6Cidrs == len(matchingIpv6Cidrs) {
						if numExpectedPrefixLists == len(matchingPrefixLists) {
							if numExpectedSGs == len(matchingSGs) {
								// confirm that self references match
								var lSelf bool
								var rSelf bool
								if _, ok := l["self"]; ok {
									lSelf = l["self"].(bool)
								}
								if _, ok := r["self"]; ok {
									rSelf = r["self"].(bool)
								}
								if rSelf == lSelf {
									delete(r, "self")
									// pop local cidrs from remote
									diffCidr := remoteCidrSet.Difference(localCidrSet)
									var newCidr []string
									for _, cRaw := range diffCidr.List() {
										newCidr = append(newCidr, cRaw.(string))
									}

									// reassigning
									if len(newCidr) > 0 {
										r["cidr_blocks"] = newCidr
									} else {
										delete(r, "cidr_blocks")
									}

									//// IPV6
									//// Comparison
									diffIpv6Cidr := remoteIpv6CidrSet.Difference(localIpv6CidrSet)
									var newIpv6Cidr []string
									for _, cRaw := range diffIpv6Cidr.List() {
										newIpv6Cidr = append(newIpv6Cidr, cRaw.(string))
									}

									// reassigning
									if len(newIpv6Cidr) > 0 {
										r["ipv6_cidr_blocks"] = newIpv6Cidr
									} else {
										delete(r, "ipv6_cidr_blocks")
									}

									// pop local prefix lists from remote
									diffPrefixLists := remotePrefixListsSet.Difference(localPrefixListsSet)
									var newPrefixLists []string
									for _, pRaw := range diffPrefixLists.List() {
										newPrefixLists = append(newPrefixLists, pRaw.(string))
									}

									// reassigning
									if len(newPrefixLists) > 0 {
										r["prefix_list_ids"] = newPrefixLists
									} else {
										delete(r, "prefix_list_ids")
									}

									// pop local sgs from remote
									diffSGs := remoteSGSet.Difference(localSGSet)
									if len(diffSGs.List()) > 0 {
										r["security_groups"] = diffSGs
									} else {
										delete(r, "security_groups")
									}

									// copy over any remote rule description
									if _, ok := r["description"]; ok {
										l["description"] = r["description"]
									}

									saves = append(saves, l)
								}
							}
						}

					}
				}
			}
		}
	}
	// Here we catch any remote rules that have not been stripped of all self,
	// cidrs, and security groups. We'll add remote rules here that have not been
	// matched locally, and let the graph sort things out. This will happen when
	// rules are added externally to Terraform
	for _, r := range remote {
		var lenCidr, lenIpv6Cidr, lenPrefixLists, lenSGs int
		if rCidrs, ok := r["cidr_blocks"]; ok {
			lenCidr = len(rCidrs.([]string))
		}
		if rIpv6Cidrs, ok := r["ipv6_cidr_blocks"]; ok {
			lenIpv6Cidr = len(rIpv6Cidrs.([]string))
		}
		if rPrefixLists, ok := r["prefix_list_ids"]; ok {
			lenPrefixLists = len(rPrefixLists.([]string))
		}
		if rawSGs, ok := r["security_groups"]; ok {
			lenSGs = len(rawSGs.(*schema.Set).List())
		}

		if _, ok := r["self"]; ok {
			if r["self"].(bool) {
				lenSGs++
			}
		}

		if lenSGs+lenCidr+lenIpv6Cidr+lenPrefixLists > 0 {
			log.Printf("[DEBUG] Found a remote Rule that wasn't empty: (%#v)", r)
			saves = append(saves, r)
		}
	}

	return saves
}

// Duplicate ingress/egress block structure and fill out all
// the required fields
func resourceSecurityGroupCopyRule(src map[string]interface{}, self bool, k string, v interface{}) map[string]interface{} {
	var keys_to_copy = []string{"description", "from_port", "to_port", "protocol"}

	dst := make(map[string]interface{})
	for _, key := range keys_to_copy {
		if val, ok := src[key]; ok {
			dst[key] = val
		}
	}
	if k != "" {
		dst[k] = v
	}
	if _, ok := src["self"]; ok {
		dst["self"] = self
	}
	return dst
}

// Given a set of SG rules (ingress/egress blocks), this function
// will group the rules by from_port/to_port/protocol/description
// tuples. This is inverse operation of
// SecurityGroupExpandRules()
//
// For more detail, see comments for
// SecurityGroupExpandRules()
func SecurityGroupCollapseRules(ruleset string, rules []interface{}) []interface{} {

	var keys_to_collapse = []string{"cidr_blocks", "ipv6_cidr_blocks", "prefix_list_ids", "security_groups"}

	collapsed := make(map[string]map[string]interface{})

	for _, rule := range rules {
		r := rule.(map[string]interface{})

		ruleHash := idCollapseHash(ruleset, r["protocol"].(string), int64(r["to_port"].(int)), int64(r["from_port"].(int)), r["description"].(string))

		if _, ok := collapsed[ruleHash]; ok {
			if v, ok := r["self"]; ok && v.(bool) {
				collapsed[ruleHash]["self"] = r["self"]
			}
		} else {
			collapsed[ruleHash] = r
			continue
		}

		for _, key := range keys_to_collapse {
			if _, ok := r[key]; ok {
				if _, ok := collapsed[ruleHash][key]; ok {
					if key == "security_groups" {
						collapsed[ruleHash][key] = collapsed[ruleHash][key].(*schema.Set).Union(r[key].(*schema.Set))
					} else {
						collapsed[ruleHash][key] = append(collapsed[ruleHash][key].([]interface{}), r[key].([]interface{})...)
					}
				} else {
					collapsed[ruleHash][key] = r[key]
				}
			}
		}
	}

	values := make([]interface{}, 0, len(collapsed))
	for _, val := range collapsed {
		values = append(values, val)
	}
	return values
}

// SecurityGroupExpandRules works in pair with
// SecurityGroupCollapseRules and is used as a
// workaround for the problem explained in
// https://github.com/hashicorp/terraform-provider-aws/pull/4726
//
// This function converts every ingress/egress block that
// contains multiple rules to multiple blocks with only one
// rule. Doing a Difference operation on such a normalized
// set helps to avoid unnecessary removal of unchanged
// rules during the Apply step.
//
// For example, in terraform syntax, the following block:
//
//	ingress {
//	  from_port = 80
//	  to_port = 80
//	  protocol = "tcp"
//	  cidr_blocks = [
//	    "192.168.0.1/32",
//	    "192.168.0.2/32",
//	  ]
//	}
//
// will be converted to the two blocks below:
//
//	ingress {
//	  from_port = 80
//	  to_port = 80
//	  protocol = "tcp"
//	  cidr_blocks = [ "192.168.0.1/32" ]
//	}
//
//	ingress {
//	  from_port = 80
//	  to_port = 80
//	  protocol = "tcp"
//	  cidr_blocks = [ "192.168.0.2/32" ]
//	}
//
// Then the Difference operation is executed on the new set
// to find which rules got modified, and the resulting set
// is then passed to SecurityGroupCollapseRules
// to convert the "diff" back to a more compact form for
// execution. Such compact form helps reduce the number of
// API calls.
func SecurityGroupExpandRules(rules *schema.Set) *schema.Set {
	var keys_to_expand = []string{"cidr_blocks", "ipv6_cidr_blocks", "prefix_list_ids", "security_groups"}

	normalized := schema.NewSet(SecurityGroupRuleHash, nil)

	for _, rawRule := range rules.List() {
		rule := rawRule.(map[string]interface{})

		if v, ok := rule["self"]; ok && v.(bool) {
			new_rule := resourceSecurityGroupCopyRule(rule, true, "", nil)
			normalized.Add(new_rule)
		}
		for _, key := range keys_to_expand {
			item, exists := rule[key]
			if exists {
				var list []interface{}
				if key == "security_groups" {
					list = item.(*schema.Set).List()
				} else {
					list = item.([]interface{})
				}
				for _, v := range list {
					var new_rule map[string]interface{}
					if key == "security_groups" {
						new_v := schema.NewSet(schema.HashString, nil)
						new_v.Add(v)
						new_rule = resourceSecurityGroupCopyRule(rule, false, key, new_v)
					} else {
						new_v := make([]interface{}, 0)
						new_v = append(new_v, v)
						new_rule = resourceSecurityGroupCopyRule(rule, false, key, new_v)
					}
					normalized.Add(new_rule)
				}
			}
		}
	}

	return normalized
}

// Convert type-to_port-from_port-protocol-description tuple
// to a hash to use as a key in Set.
func idCollapseHash(rType, protocol string, toPort, fromPort int64, description string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", rType))
	buf.WriteString(fmt.Sprintf("%d-", toPort))
	buf.WriteString(fmt.Sprintf("%d-", fromPort))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(protocol)))
	buf.WriteString(fmt.Sprintf("%s-", description))

	return fmt.Sprintf("rule-%d", create.StringHashcode(buf.String()))
}

// Creates a unique hash for the type, ports, and protocol, used as a key in
// maps
func idHash(rType, protocol string, toPort, fromPort int64, self bool) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", rType))
	buf.WriteString(fmt.Sprintf("%d-", toPort))
	buf.WriteString(fmt.Sprintf("%d-", fromPort))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(protocol)))
	buf.WriteString(fmt.Sprintf("%t-", self))

	return fmt.Sprintf("rule-%d", create.StringHashcode(buf.String()))
}

// ProtocolStateFunc ensures we only store a string in any protocol field
func ProtocolStateFunc(v interface{}) string {
	switch v := v.(type) {
	case string:
		p := ProtocolForValue(v)
		return p
	default:
		log.Printf("[WARN] Non String value given for Protocol: %#v", v)
		return ""
	}
}

// ProtocolForValue converts a valid Internet Protocol number into it's name
// representation. If a name is given, it validates that it's a proper protocol
// name. Names/numbers are as defined at
// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
func ProtocolForValue(v string) string {
	// special case -1
	protocol := strings.ToLower(v)
	if protocol == "-1" || protocol == "all" {
		return "-1"
	}
	// if it's a name like tcp, return that
	if _, ok := sgProtocolIntegers()[protocol]; ok {
		return protocol
	}
	// convert to int, look for that value
	p, err := strconv.Atoi(protocol)
	if err != nil {
		// we were unable to convert to int, suggesting a string name, but it wasn't
		// found above
		log.Printf("[WARN] Unable to determine valid protocol: %s", err)
		return protocol
	}

	for k, v := range sgProtocolIntegers() {
		if p == v {
			// guard against protocolIntegers sometime in the future not having lower
			// case ids in the map
			return strings.ToLower(k)
		}
	}

	// fall through
	log.Printf("[WARN] Unable to determine valid protocol: no matching protocols found")
	return protocol
}

// a map of protocol names and their codes, defined at
// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml,
// documented to be supported by AWS Security Groups
// http://docs.aws.amazon.com/fr_fr/AWSEC2/latest/APIReference/API_IpPermission.html
// Similar to protocolIntegers() used by Network ACLs, but explicitly only
// supports "tcp", "udp", "icmp", "icmpv6", and "all"
func sgProtocolIntegers() map[string]int {
	return map[string]int{
		"icmpv6": 58,
		"udp":    17,
		"tcp":    6,
		"icmp":   1,
		"all":    -1,
	}
}

// The AWS Lambda service creates ENIs behind the scenes and keeps these around for a while
// which would prevent SGs attached to such ENIs from being destroyed
func deleteLingeringLambdaENIs(conn *ec2.EC2, filterName, resourceId string, timeout time.Duration) error {
	// AWS Lambda service team confirms P99 deletion time of ~35 minutes. Buffer for safety.
	if minimumTimeout := 45 * time.Minute; timeout < minimumTimeout {
		timeout = minimumTimeout
	}

	networkInterfaces, err := FindNetworkInterfaces(conn, &ec2.DescribeNetworkInterfacesInput{
		Filters: BuildAttributeFilterList(map[string]string{
			filterName:    resourceId,
			"description": "AWS Lambda VPC ENI*",
		}),
	})

	if err != nil {
		return fmt.Errorf("error listing EC2 Network Interfaces: %w", err)
	}

	for _, v := range networkInterfaces {
		networkInterfaceID := aws.StringValue(v.NetworkInterfaceId)

		if v.Attachment != nil && aws.StringValue(v.Attachment.InstanceOwnerId) == "amazon-aws" {
			networkInterface, err := WaitNetworkInterfaceAvailableAfterUse(conn, networkInterfaceID, timeout)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return fmt.Errorf("error waiting for Lambda ENI (%s) to become available for detachment: %w", networkInterfaceID, err)
			}

			v = networkInterface
		}

		if v.Attachment != nil {
			err = DetachNetworkInterface(conn, networkInterfaceID, aws.StringValue(v.Attachment.AttachmentId), timeout)

			if err != nil {
				return fmt.Errorf("error detaching Lambda ENI (%s): %w", networkInterfaceID, err)
			}
		}

		err = DeleteNetworkInterface(conn, networkInterfaceID)

		if err != nil {
			return fmt.Errorf("error deleting Lambda ENI (%s): %w", networkInterfaceID, err)
		}
	}

	return nil
}

func initSecurityGroupRule(ruleMap map[string]map[string]interface{}, perm *ec2.IpPermission, desc string) map[string]interface{} {
	var fromPort, toPort int64
	if v := perm.FromPort; v != nil {
		fromPort = aws.Int64Value(v)
	}
	if v := perm.ToPort; v != nil {
		toPort = aws.Int64Value(v)
	}
	k := fmt.Sprintf("%s-%d-%d-%s", *perm.IpProtocol, fromPort, toPort, desc)
	rule, ok := ruleMap[k]
	if !ok {
		rule = make(map[string]interface{})
		ruleMap[k] = rule
	}
	rule["protocol"] = aws.StringValue(perm.IpProtocol)
	rule["from_port"] = fromPort
	rule["to_port"] = toPort
	if desc != "" {
		rule["description"] = desc
	}

	return rule
}
