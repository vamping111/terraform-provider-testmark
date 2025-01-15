package directconnect_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDirectConnectLagDataSource_basic(t *testing.T) {
	key := "DX_LAG_NAME"
	rName := os.Getenv(key)
	if rName == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	datasourceName := "data.aws_dx_lag.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, directconnect.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLagDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "arn"),
					resource.TestCheckResourceAttrSet(datasourceName, "aws_device"),
					resource.TestCheckResourceAttrSet(datasourceName, "bandwidth"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "location"),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "owner_account_id"),
				),
			},
		},
	})
}

func testAccLagDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_dx_lag" "test" {
  name = %[1]q
}
`, rName)
}
