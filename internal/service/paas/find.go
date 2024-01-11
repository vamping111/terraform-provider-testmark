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

func FindBackupUsers(conn *paas.PaaS) ([]*paas.BackupUser, error) {
	input := &paas.ListBackupUsersInput{}

	output, err := conn.ListBackupUsers(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Users == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Users, nil
}

func FindBackupById(conn *paas.PaaS, id string) (*paas.Backup, error) {
	input := &paas.DescribeBackupInput{
		BackupId: aws.String(id),
	}

	output, err := conn.DescribeBackup(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Backup == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Backup, nil
}

func FindBackups(conn *paas.PaaS, serviceClass, serviceId, serviceType string) ([]*paas.Backup, error) {
	input := &paas.ListBackupsInput{}

	if serviceClass != "" {
		input.ServiceClass = aws.String(serviceClass)
	}

	if serviceId != "" {
		input.ServiceId = aws.String(serviceId)
	}

	if serviceType != "" {
		input.ServiceType = aws.String(serviceType)
	}

	output, err := conn.ListBackups(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Backups == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Backups, nil
}
