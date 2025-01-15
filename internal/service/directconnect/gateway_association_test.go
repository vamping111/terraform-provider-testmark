package directconnect_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/directconnect"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdirectconnect "github.com/hashicorp/terraform-provider-aws/internal/service/directconnect"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDirectConnectGatewayAssociation_basicTransitGatewaySingleAccount(t *testing.T) {
	resourceName := "aws_dx_gateway_association.test"
	resourceNameDxGw := "aws_dx_gateway.test"
	resourceNameTgw := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	var ga directconnect.GatewayAssociation
	var gap directconnect.GatewayAssociationProposal

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, directconnect.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckGatewayAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxGatewayAssociationConfig_basicTransitGatewaySingleAccount(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayAssociationExists(resourceName, &ga, &gap),
					resource.TestCheckResourceAttr(resourceName, "allowed_prefixes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_prefixes.*", "10.255.255.0/30"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_prefixes.*", "10.255.255.8/30"),
					resource.TestCheckResourceAttrPair(resourceName, "associated_gateway_id", resourceNameTgw, "id"),
					acctest.CheckResourceAttrAccountID(resourceName, "associated_gateway_owner_account_id"),
					resource.TestCheckResourceAttr(resourceName, "associated_gateway_type", "transitGateway"),
					resource.TestCheckResourceAttrPair(resourceName, "dx_gateway_id", resourceNameDxGw, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccGatewayAssociationImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGatewayAssociationImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["dx_gateway_id"], rs.Primary.Attributes["associated_gateway_id"]), nil
	}
}

func testAccCheckGatewayAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).DirectConnectConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dx_gateway_association" {
			continue
		}

		_, err := tfdirectconnect.FindGatewayAssociationByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Direct Connect Gateway Association %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCheckGatewayAssociationExists(name string, ga *directconnect.GatewayAssociation, gap *directconnect.GatewayAssociationProposal) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Direct Connect Gateway Association ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DirectConnectConn

		output, err := tfdirectconnect.FindGatewayAssociationByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		if proposalID := rs.Primary.Attributes["proposal_id"]; proposalID != "" {
			output, err := tfdirectconnect.FindGatewayAssociationProposalByID(conn, proposalID)

			if err != nil {
				return err
			}

			*gap = *output
		}

		*ga = *output

		return nil
	}
}

func testAccDxGatewayAssociationConfig_basicTransitGatewaySingleAccount(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_dx_gateway" "test" {
  name            = %[1]q
  amazon_side_asn = "%[2]d"
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_dx_gateway_association" "test" {
  dx_gateway_id         = aws_dx_gateway.test.id
  associated_gateway_id = aws_ec2_transit_gateway.test.id

  allowed_prefixes = [
    "10.255.255.0/30",
    "10.255.255.8/30",
  ]
}
`, rName, rBgpAsn)
}
