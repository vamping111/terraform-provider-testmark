package directconnect_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/directconnect"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDirectConnectTransitVirtualInterface_serial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic": testAccTransitVirtualInterface_basic,
		"tags":  testAccTransitVirtualInterface_Tags,
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccTransitVirtualInterface_basic(t *testing.T) {
	connectionNameEnvVar := "DX_CONNECTION_NAME"
	connectionName := os.Getenv(connectionNameEnvVar)
	if connectionName == "" {
		t.Skipf("Environment variable %s is not set", connectionNameEnvVar)
	}

	vlanEnvVar := "DX_VLAN"
	vlan, err := strconv.Atoi(os.Getenv(vlanEnvVar))
	if err != nil {
		t.Skipf("Environment variable %s is not set or its value is not a valid integer", vlanEnvVar)
	}

	var vif directconnect.VirtualInterface
	resourceName := "aws_dx_transit_virtual_interface.test"
	dxGatewayResourceName := "aws_dx_gateway.test"
	connectionDatasourceName := "data.aws_dx_connection.test"
	rName := fmt.Sprintf("tf-testacc-transit-vif-%s", sdkacctest.RandString(9))

	amzAsn := sdkacctest.RandIntRange(64512, 65534)
	var bgpAsn int
	for {
		bgpAsn = sdkacctest.RandIntRange(64512, 65534)
		if bgpAsn != amzAsn {
			break
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, directconnect.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckTransitVirtualInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxTransitVirtualInterfaceConfig_basic(connectionName, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitVirtualInterfaceExists(resourceName, &vif),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_address"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_side_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_device"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", strconv.Itoa(bgpAsn)),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_auth_key"),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id", connectionDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_address"),
					resource.TestCheckResourceAttrPair(resourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "vlan", strconv.Itoa(vlan)),
				),
			},
			{
				Config: testAccDxTransitVirtualInterfaceConfig_updated(connectionName, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitVirtualInterfaceExists(resourceName, &vif),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_address"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_side_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_device"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", strconv.Itoa(bgpAsn)),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_auth_key"),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id", connectionDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_address"),
					resource.TestCheckResourceAttrPair(resourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "vlan", strconv.Itoa(vlan)),
				),
			},
			// Test import.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitVirtualInterface_Tags(t *testing.T) {
	connectionNameEnvVar := "DX_CONNECTION_NAME"
	connectionName := os.Getenv(connectionNameEnvVar)
	if connectionName == "" {
		t.Skipf("Environment variable %s is not set", connectionNameEnvVar)
	}

	vlanEnvVar := "DX_VLAN"
	vlan, err := strconv.Atoi(os.Getenv(vlanEnvVar))
	if err != nil {
		t.Skipf("Environment variable %s is not set or its value is not valid integer", vlanEnvVar)
	}

	var vif directconnect.VirtualInterface
	resourceName := "aws_dx_transit_virtual_interface.test"
	dxGatewayResourceName := "aws_dx_gateway.test"
	connectionDatasourceName := "data.aws_dx_connection.test"
	rName := fmt.Sprintf("tf-testacc-transit-vif-%s", sdkacctest.RandString(9))

	amzAsn := sdkacctest.RandIntRange(64512, 65534)
	var bgpAsn int
	for {
		bgpAsn = sdkacctest.RandIntRange(64512, 65534)
		if bgpAsn != amzAsn {
			break
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, directconnect.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckTransitVirtualInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxTransitVirtualInterfaceConfig_tags(connectionName, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitVirtualInterfaceExists(resourceName, &vif),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_address"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_side_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_device"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", strconv.Itoa(bgpAsn)),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_auth_key"),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id", connectionDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_address"),
					resource.TestCheckResourceAttrPair(resourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2a"),
					resource.TestCheckResourceAttr(resourceName, "vlan", strconv.Itoa(vlan)),
				),
			},
			{
				Config: testAccDxTransitVirtualInterfaceConfig_tagsUpdated(connectionName, rName, amzAsn, bgpAsn, vlan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitVirtualInterfaceExists(resourceName, &vif),
					resource.TestCheckResourceAttr(resourceName, "address_family", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_address"),
					resource.TestCheckResourceAttrSet(resourceName, "amazon_side_asn"),
					resource.TestCheckResourceAttrSet(resourceName, "aws_device"),
					resource.TestCheckResourceAttr(resourceName, "bgp_asn", strconv.Itoa(bgpAsn)),
					resource.TestCheckResourceAttrSet(resourceName, "bgp_auth_key"),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id", connectionDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_address"),
					resource.TestCheckResourceAttrPair(resourceName, "dx_gateway_id", dxGatewayResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2b"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "Value3"),
					resource.TestCheckResourceAttr(resourceName, "vlan", strconv.Itoa(vlan)),
				),
			},
			// Test import.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTransitVirtualInterfaceExists(name string, vif *directconnect.VirtualInterface) resource.TestCheckFunc {
	return testAccCheckDxVirtualInterfaceExists(name, vif)
}

func testAccCheckTransitVirtualInterfaceDestroy(s *terraform.State) error {
	return testAccCheckDxVirtualInterfaceDestroy(s, "aws_dx_transit_virtual_interface")
}

func testAccDxTransitVirtualInterfaceConfig_base(rName string, amzAsn int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[1]q
  amazon_side_asn = %[2]d
}
`, rName, amzAsn)
}

func testAccDxTransitVirtualInterfaceConfig_basic(connectionName, rName string, amzAsn, bgpAsn, vlan int) string {
	return testAccDxTransitVirtualInterfaceConfig_base(rName, amzAsn) + fmt.Sprintf(`
data "aws_dx_connection" "test" {
  name = %[1]q
}

resource "aws_dx_transit_virtual_interface" "test" {
  address_family = "ipv4"
  bgp_asn        = %[3]d
  dx_gateway_id  = aws_dx_gateway.test.id
  connection_id  = data.aws_dx_connection.test.id
  name           = %[2]q
  vlan           = %[4]d
}
`, connectionName, rName, bgpAsn, vlan)
}

func testAccDxTransitVirtualInterfaceConfig_updated(connectionName, rName string, amzAsn, bgpAsn, vlan int) string {
	return testAccDxTransitVirtualInterfaceConfig_base(rName, amzAsn) + fmt.Sprintf(`
data "aws_dx_connection" "test" {
  name = %[1]q
}

resource "aws_dx_transit_virtual_interface" "test" {
  address_family = "ipv4"
  bgp_asn        = %[3]d
  dx_gateway_id  = aws_dx_gateway.test.id
  connection_id  = data.aws_dx_connection.test.id
  name           = %[2]q
  vlan           = %[4]d
}
`, connectionName, rName, bgpAsn, vlan)
}

func testAccDxTransitVirtualInterfaceConfig_tags(connectionName, rName string, amzAsn, bgpAsn, vlan int) string {
	return testAccDxTransitVirtualInterfaceConfig_base(rName, amzAsn) + fmt.Sprintf(`
data "aws_dx_connection" "test" {
  name = %[1]q
}

resource "aws_dx_transit_virtual_interface" "test" {
  address_family = "ipv4"
  bgp_asn        = %[3]d
  dx_gateway_id  = aws_dx_gateway.test.id
  connection_id  = data.aws_dx_connection.test.id
  name           = %[2]q
  vlan           = %[4]d

  tags = {
    Name = %[2]q
    Key1 = "Value1"
    Key2 = "Value2a"
  }
}
`, connectionName, rName, bgpAsn, vlan)
}

func testAccDxTransitVirtualInterfaceConfig_tagsUpdated(connectionName, rName string, amzAsn, bgpAsn, vlan int) string {
	return testAccDxTransitVirtualInterfaceConfig_base(rName, amzAsn) + fmt.Sprintf(`
data "aws_dx_connection" "test" {
  name = %[1]q
}

resource "aws_dx_transit_virtual_interface" "test" {
  address_family = "ipv4"
  bgp_asn        = %[3]d
  dx_gateway_id  = aws_dx_gateway.test.id
  connection_id  = data.aws_dx_connection.test.id
  name           = %[2]q
  vlan           = %[4]d

  tags = {
    Name = %[2]q
    Key2 = "Value2b"
    Key3 = "Value3"
  }
}
`, connectionName, rName, bgpAsn, vlan)
}
