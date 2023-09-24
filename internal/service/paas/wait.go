package paas

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func waitServiceCreated(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ServiceStatusPending, ServiceStatusClaimed, ServiceStatusProvisioning},
		Target:  []string{ServiceStatusReady},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

func waitServiceUpdated(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ServiceStatusUpdating},
		Target:  []string{ServiceStatusReady},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

func waitServiceDeleted(ctx context.Context, conn *paas.PaaS, id string, timeout time.Duration) (*paas.Service, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ServiceStatusPending, ServiceStatusClaimed},
		Target:  []string{ServiceStatusDeleted},
		Refresh: statusService(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*paas.Service); ok {
		if err != nil {
			setServiceErrorToResourceLastError(output, err)
		}

		return output, err
	}

	return nil, err
}

func setServiceErrorToResourceLastError(service *paas.Service, err error) {
	if status := aws.StringValue(service.Status); status != ServiceStatusError {
		return
	}

	errCode := aws.StringValue(service.ErrorCode)
	errDesc := aws.StringValue(service.ErrorDescription)
	tfresource.SetLastError(err, fmt.Errorf("code: %s, description: %s", errCode, errDesc))
}
