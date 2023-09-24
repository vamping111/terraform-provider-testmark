package paas

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindServiceByID(conn *paas.PaaS, id string) (*paas.Service, error) {
	input := &paas.DescribeServiceInput{
		ServiceId: aws.String(id),
	}

	output, err := conn.DescribeService(input)

	// TODO: fix parsing error message (__type -> code) in sdk
	if tfawserr.ErrCodeEquals(err, "Document.NotFound") {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Service == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Service, nil
}
